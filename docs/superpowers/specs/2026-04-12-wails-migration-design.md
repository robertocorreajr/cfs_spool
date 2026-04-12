# CFS Spool — Migração para Wails + React + shadcn/ui

## Contexto

O CFS Spool é uma aplicação Go para ler e gravar tags RFID (MIFARE Classic) de filamentos Creality. Atualmente funciona como um servidor HTTP local (porta 8080) com frontend web embutido (HTML/CSS/JS vanilla). A aplicação apresenta um bug no tratamento de cores (aceita caracteres hex inválidos) e a experiência de usuário é limitada pelo frontend vanilla.

**Objetivo**: migrar para uma aplicação desktop nativa usando **Wails v2** (Go backend + frontend web), com **React + TypeScript**, **shadcn/ui** e **Tailwind CSS** no frontend. O resultado será uma aplicação mais bonita, com melhor UX (especialmente no color picker), e empacotada nativamente para macOS, Linux e Windows.

---

## Arquitetura

### Stack

| Camada | Tecnologia |
|--------|-----------|
| Framework desktop | Wails v2 |
| Backend | Go 1.24+ (CGO para PC/SC) |
| Frontend | React 18 + TypeScript + Vite |
| UI Components | shadcn/ui |
| Styling | Tailwind CSS |
| Color picker | react-colorful |
| RFID | github.com/ebfe/scard (PC/SC) |
| Comunicação | Wails bindings (Go ↔ JS direto, sem REST) |

### Fluxo de dados

```
Wails Window (nativa)
  └─ React + shadcn/ui + Tailwind CSS
       │
       └─ Wails Bindings (chamadas Go diretas)
            ├─ app.go → ReadTag() / WriteTag() / GetOptions() / ValidateColor()
            ├─ internal/rfid/reader.go → PC/SC APDU → ACR122U → MIFARE Tag
            ├─ internal/creality/crypto.go → AES-ECB (2 chaves)
            └─ internal/creality/fields.go → layout 38 bytes
```

### Estrutura do projeto

```
cfs_spool/
├── main.go                          # Wails entry point (go:embed frontend/dist)
├── app.go                           # Métodos Go expostos via bindings
├── app_options.go                   # Dados estáticos (materiais, fornecedores, comprimentos)
├── wails.json                       # Configuração Wails
├── go.mod / go.sum
├── internal/                        # INALTERADO
│   ├── creality/
│   │   ├── crypto.go
│   │   └── fields.go
│   └── rfid/
│       └── reader.go
├── frontend/                        # React app (substitui web/)
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   ├── components.json              # shadcn/ui config
│   ├── index.html
│   └── src/
│       ├── main.tsx
│       ├── App.tsx
│       ├── globals.css
│       ├── lib/utils.ts             # cn() do shadcn
│       ├── components/
│       │   ├── ui/                  # shadcn (auto-gerado)
│       │   ├── SpoolForm.tsx        # Formulário principal
│       │   ├── ColorPicker.tsx      # Color picker completo
│       │   ├── MaterialSelect.tsx   # Selects vinculados fornecedor/material
│       │   ├── LengthSelect.tsx     # Comprimento + input customizado
│       │   ├── TagStatus.tsx        # Status da leitura (UID)
│       │   └── Header.tsx
│       ├── data/presets.ts          # 35 cores pré-definidas
│       ├── hooks/useSpoolForm.ts    # State do formulário
│       └── types/spool.ts          # Interfaces TS para bindings Go
├── build/                           # Assets Wails (ícone, Info.plist, etc.)
├── assets/icons/cfs-spool.svg
├── tests/
├── Makefile
├── CLAUDE.md
├── README.md / README.pt-BR.md
└── .github/workflows/
```

---

## Go Backend — Bindings (app.go)

### Tipos

```go
type TagData struct {
    UID          string `json:"uid"`
    Date         string `json:"date"`          // YYYY-MM-DD (para input date)
    DateDisplay  string `json:"dateDisplay"`    // "12 de Abril de 2026"
    SupplierCode string `json:"supplierCode"`   // "0276"
    SupplierName string `json:"supplierName"`   // "Creality"
    MaterialCode string `json:"materialCode"`   // "04001"
    MaterialName string `json:"materialName"`   // "CR-PLA"
    Color        string `json:"color"`          // "77BB41" (6 chars, uppercase, sem prefixo)
    LengthCode   string `json:"lengthCode"`     // "0330"
    LengthDisplay string `json:"lengthDisplay"` // "330cm (1kg)"
    Serial       string `json:"serial"`         // "000001"
}

type WriteRequest struct {
    Date     string `json:"date"`     // YYYY-MM-DD
    Supplier string `json:"supplier"` // código 4 chars
    Material string `json:"material"` // código 5 chars
    Color    string `json:"color"`    // 6 chars hex (sem # ou prefixo 0)
    Length   string `json:"length"`   // código 4 chars ou gramas
    Serial   string `json:"serial"`   // até 6 dígitos
}

type OptionsResponse struct {
    Materials []MaterialOption `json:"materials"`
    Vendors   []VendorOption   `json:"vendors"`
    Lengths   []LengthOption   `json:"lengths"`
}
```

### Métodos expostos

| Método | Descrição |
|--------|-----------|
| `ReadTag() (*TagData, error)` | Abre leitor → lê UID → autentica → lê blocos 4-6 → decripta → parseia → retorna dados com códigos E nomes, data em YYYY-MM-DD |
| `WriteTag(data WriteRequest) error` | Valida → converte data/material/comprimento → SetColor com validação hex → encripta → grava tag |
| `GetOptions() OptionsResponse` | Retorna listas estáticas de materiais, fornecedores e comprimentos |
| `ValidateColor(hex string) (string, error)` | Valida `^[0-9A-Fa-f]{6}$`, retorna uppercase ou erro |

### Helpers privados (migrados de cmd/app/main.go)

- `convertDate(isoDate string) (string, error)` — "2026-04-12" → "26412"
- `parseDateToISO(date5 string) string` — "26412" → "2026-04-12" (NOVO, reverso). Formato interno: char[0-1]=ano (YY), char[2]=mês (1-9 para Jan-Set, então conforme FormatDate em fields.go), char[3-4]=dia (DD). Ex: "26412" = ano 26, mês 4 (Abril), dia 12
- `convertMaterial(name string) string` — nome → código 5 chars
- `convertLength(code string) string` — mapeia comprimentos padrão ou converte gramas
- `padSerial(serial string) string` — preenche com zeros à esquerda até 6 dígitos

---

## Frontend React

### Componentes shadcn/ui utilizados

Button, Card, Input, Label, Select, SelectContent, SelectItem, SelectTrigger, SelectValue, Separator, Sonner (toast), Badge

### Dependência extra

`react-colorful` — color picker leve (3KB gzipped), fornece `HexColorPicker`

### Layout da tela única

```
┌─────────────────────────────────────────────────┐
│  CFS Spool                            v3.0.0    │
├─────────────────────────────────────────────────┤
│  [🔍 Ler Tag]              UID: A1B2C3D4       │
│─────────────────────────────────────────────────│
│                                                  │
│  Data           [____2026-04-12____________]     │
│  Fornecedor     [▼ Creality________________]     │
│  Material       [▼ CR-PLA_________________]      │
│                                                  │
│  Cor   [■ preview] [__77BB41__]                  │
│        ┌────────────────────────┐                │
│        │  HexColorPicker        │                │
│        │  (gradiente + hue bar) │                │
│        └────────────────────────┘                │
│        [■][■][■][■][■][■][■][■][■][■][■][■]    │
│        [■][■][■][■][■][■][■][■][■][■][■][■]    │
│        [■][■][■][■][■][■][■][■][■][■][■]       │
│                                                  │
│  Comprimento    [▼ 330cm (1kg)____________]      │
│  Serial         [____000001_______________]      │
│                                                  │
│  [            ✍️ Gravar Tag              ]       │
└─────────────────────────────────────────────────┘
```

### Hierarquia de componentes

```
App.tsx
├── <Toaster /> (sonner — notificações globais)
├── <Header /> (título + versão)
└── <SpoolForm />
    ├── <TagStatus /> (botão Ler Tag + exibição UID)
    ├── <Separator />
    ├── <Input type="date" /> (data)
    ├── <MaterialSelect /> (fornecedor + material vinculados)
    │   ├── <Select> Fornecedor
    │   └── <Select> Material (filtrado por fornecedor)
    ├── <ColorPicker />
    │   ├── preview (div com background-color)
    │   ├── <Input> hex (6 chars, validado)
    │   ├── <HexColorPicker /> (react-colorful)
    │   └── grade de swatches (35 botões coloridos)
    ├── <LengthSelect />
    │   ├── <Select> comprimentos padrão
    │   └── <Input> gramas customizado (condicional)
    ├── <Input> serial
    └── <Button> Gravar Tag
```

### Color Picker — Detalhes

- `react-colorful` fornece `HexColorPicker` que aceita/retorna formato `#RRGGBB`
- Conversão interna: componente trabalha com 6 chars sem `#`, converte na interface com react-colorful
- Input hex: controlled input com validação `/^[0-9a-fA-F]{0,6}$/` — **corrige o bug de cor no frontend**
- Swatches: 35 `<button>` em grid Tailwind (`grid grid-cols-12 gap-1`), cada um `w-7 h-7 rounded cursor-pointer border`, com `ring-2 ring-primary` quando selecionado
- Cores branca/preta têm `border-border` para visibilidade

### Filtragem Fornecedor/Material

Mesma lógica do frontend atual:
- Fornecedor "Genérico" (0000) → mostra materiais com código < "04000"
- Fornecedor "Creality" (0276) → mostra materiais com código >= "04000"
- Ao selecionar material: auto-seleciona fornecedor correspondente
- Ao selecionar fornecedor: filtra lista de materiais

### Fluxo "Ler Tag → Preencher Formulário"

1. Usuário clica "Ler Tag"
2. Botão fica desabilitado, mostra loading
3. React chama `window.go.main.App.ReadTag()` (binding Wails)
4. Go: abre leitor → lê UID → autentica → lê blocos → decripta → parseia → retorna `TagData`
5. React recebe `TagData` e atualiza todos os campos do formulário:
   - `date` = `tagData.date` (YYYY-MM-DD)
   - `supplier` = `tagData.supplierCode` (seleciona no dropdown)
   - `material` = `tagData.materialCode` (seleciona no dropdown)
   - `color` = `tagData.color` (6 chars hex, atualiza picker e preview)
   - `length` = `tagData.lengthCode` (seleciona no dropdown)
   - `serial` = `tagData.serial`
6. Exibe UID e toast de sucesso
7. Usuário edita qualquer campo e clica "Gravar Tag"

### Estado do formulário (useSpoolForm hook)

```typescript
interface SpoolFormState {
  date: string;          // YYYY-MM-DD
  supplier: string;      // código 4 chars
  material: string;      // código 5 chars
  color: string;         // 6 chars hex uppercase
  length: string;        // código 4 chars
  customGrams: string;   // gramas (quando comprimento customizado)
  serial: string;        // até 6 dígitos
  uid: string;           // somente leitura (exibição)
  isReading: boolean;
  isWriting: boolean;
}
```

Usa `useState` por campo — simples e adequado para ~9 campos.

---

## Correção do Bug de Cor

### Problema

1. Frontend aceita caracteres não-hex (ex: "GG0000") — valida apenas comprimento
2. Backend `SetColor()` valida apenas `len(color) != 6`, não os caracteres
3. Inconsistência de case (às vezes maiúsculo, às vezes minúsculo)

### Solução (dupla validação)

**Frontend** (`ColorPicker.tsx`):
- Input controlado aceita apenas `[0-9a-fA-F]`
- Normaliza para uppercase ao enviar

**Backend** (`app.go` → `ValidateColor()`):
- Regex `^[0-9A-Fa-f]{6}$`
- Normaliza para uppercase antes de chamar `fields.SetColor()`
- `fields.SetColor()` em `fields.go` permanece inalterado

---

## Build e Empacotamento

### Desenvolvimento

```bash
wails dev    # Hot-reload: Go backend + Vite frontend
```

### Produção

```bash
wails build                          # Binário para plataforma atual
wails build -platform darwin/arm64   # macOS Apple Silicon
wails build -platform darwin/amd64   # macOS Intel
wails build -platform linux/amd64    # Linux
wails build -platform windows/amd64  # Windows
```

### Makefile

```makefile
dev:
	wails dev

build:
	wails build

build-all:
	wails build -platform darwin/arm64
	wails build -platform darwin/amd64
	wails build -platform linux/amd64
	wails build -platform windows/amd64

test:
	go test -v ./internal/...

clean:
	rm -rf build/bin/ frontend/dist/ frontend/node_modules/

install-frontend:
	cd frontend && npm install

install-wails:
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

---

## CI/CD

### build.yml — Reescrito

```yaml
jobs:
  build:
    strategy:
      matrix:
        include:
          - os: macos-14
            platform: darwin/arm64
          - os: ubuntu-latest
            platform: linux/amd64
          - os: windows-latest
            platform: windows/amd64
    steps:
      - Setup Go 1.24
      - Setup Node.js 20
      - Install wails CLI
      - Install PC/SC deps (plataforma-específico)
      - cd frontend && npm ci
      - wails build -platform ${{ matrix.platform }}
      - Upload artifact
```

### auto-tag.yml — Mudanças mínimas
- Remover referências Docker
- Manter lógica de versionamento semântico (#major, #minor, patch)

---

## Remoções

| Item | Motivo |
|------|--------|
| `cmd/app/main.go` | Substituído por `main.go` + `app.go` na raiz |
| `cmd/cfs-spool/` | CLI removida |
| `web/` | Substituído por `frontend/` |
| `Dockerfile` | App nativa não precisa de container |
| Docker targets no Makefile | Idem |

---

## Documentação

### CLAUDE.md — Reescrever com:
- Nova arquitetura (Wails + React + shadcn)
- Comandos `wails dev` / `wails build`
- Estrutura de diretórios atualizada
- Fluxo de dados atualizado (sem REST)
- Componentes React no lugar de API endpoints

### README.md / README.pt-BR.md — Reescrever com:
- Instruções de instalação como app nativo (sem abrir browser)
- Pré-requisitos atualizados (Node.js + Wails CLI)
- Capturas de tela da nova interface
- Remover referências Docker e API REST

---

## Fases de Implementação

### Fase 1: Setup Wails + Backend Go
- Scaffold Wails com template react-ts
- Criar `main.go`, `app.go`, `app_options.go`
- Migrar lógica de conversão
- Implementar 4 métodos de binding
- Implementar `parseDateToISO()` e `ValidateColor()`
- Verificar `wails dev` inicia

### Fase 2: Frontend React + shadcn
- Configurar Tailwind + shadcn/ui
- Instalar componentes shadcn + react-colorful
- Criar todos os componentes React
- Implementar filtragem fornecedor/material
- Labels pt-BR

### Fase 3: Integração e Testes
- Testar ciclo ler → editar → gravar
- Testar color picker (35 presets + customizado + validação)
- Testar cenários de erro
- Confirmar fix do bug de cor

### Fase 4: Limpeza
- Deletar `cmd/`, `web/`, `Dockerfile`
- Atualizar `.gitignore`
- Reescrever `Makefile`

### Fase 5: Documentação e CI/CD
- Reescrever `CLAUDE.md`, `README.md`, `README.pt-BR.md`
- Reescrever `build.yml`
- Atualizar `auto-tag.yml`

---

## Verificação

### Como testar end-to-end
1. `wails dev` — app abre em janela nativa
2. Clicar "Ler Tag" com ACR122U conectado e tag presente → campos preenchidos
3. Alterar cor via color picker (gradiente + swatches + hex manual)
4. Tentar digitar "GG0000" no campo hex → deve ser rejeitado
5. Clicar "Gravar Tag" → toast de sucesso
6. Ler novamente → dados devem corresponder ao que foi gravado
7. Testar sem leitor → mensagem de erro amigável
8. Testar com tag virgem → mensagem informativa
9. `wails build` → binário nativo funcional

### Testes unitários
- `go test ./internal/...` — testes existentes continuam passando
- Novos testes para `ValidateColor()`, `parseDateToISO()`, `convertDate()` em `app_test.go`
