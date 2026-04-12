# CFS Spool - Gerenciador RFID de Filamentos Creality

Aplicativo desktop nativo (Wails v2 + React + shadcn/ui) para leitura e gravacao de tags RFID do Creality File System (CFS) utilizadas em bobinas de filamento para impressoras 3D Creality.

[![English](https://img.shields.io/badge/lang-en-blue.svg)](README.md)
[![Portuguese](https://img.shields.io/badge/lang-pt--BR-green.svg)](README.pt-BR.md)

## Funcionalidades

- **Aplicativo desktop nativo** -- sem necessidade de navegador ou servidor, executa diretamente no seu sistema operacional
- **Seletor avancado de cores**: 35 cores predefinidas baseadas no sistema Creality + seletor de cores para qualquer cor hexadecimal personalizada
- **Logica inteligente**: Auto-selecao de fornecedor baseado no material escolhido
- **Preenchimento automatico**: Campos opcionais (Lote, Serial) com padding automatico
- **Leitura visual**: Preview das cores lidas das tags existentes
- **Criptografia/descriptografia AES-ECB**: Suporte completo ao sistema de criptografia Creality
- **Autenticacao robusta**: Multiplos metodos de fallback para leitura RFID
- **Derivacao de chaves**: Algoritmo completo baseado no UID da tag
- **Compatibilidade**: Funciona com tags novas (FFFFFFFFFFFF) e usadas (chave derivada)

## Instalacao

### Downloads Prontos (Recomendado)

Baixe a versao mais recente para sua plataforma:

**[Releases - GitHub](https://github.com/robertocorreajr/cfs_spool/releases/latest)**

- **macOS (Apple Silicon)**: `cfs-spool-darwin-arm64` -- aplicativo nativo, basta executar
- **Linux (x86_64)**: `cfs-spool-linux-amd64` -- aplicativo nativo, basta executar
- **Windows (x86_64)**: `cfs-spool-windows-amd64.exe` -- aplicativo nativo, basta executar

### Compilacao a partir do Codigo Fonte

#### Pre-requisitos

- **Go 1.24+**
- **Node.js 18+**
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Leitor RFID compativel** (testado com ACR122U)
- **Headers PC/SC**:
  - macOS: incluso no sistema (nenhuma acao necessaria)
  - Linux: `sudo apt install pcscd libpcsclite-dev libgtk-3-dev libwebkit2gtk-4.1-dev`
  - Windows: incluso no sistema (winscard)

#### Comandos de Build

```bash
git clone https://github.com/robertocorreajr/cfs_spool.git
cd cfs_spool

# Modo de desenvolvimento (hot-reload)
wails dev

# Build de producao
wails build

# Compilar para todas as plataformas
make build-all
```

O binario de producao e gerado em `build/bin/`.

## Uso

1. Conecte seu leitor RFID ACR122U
2. Inicie o aplicativo (duplo clique ou execute pelo terminal)
3. **Ler Tag**: Coloque a tag no leitor e clique em "Ler Tag"
4. **Gravar Tag**: Preencha os campos e clique em "Gravar Tag"

### Selecao de Cores

Voce tem total flexibilidade para escolher cores:
- Clique em uma das **35 cores predefinidas** da paleta
- Digite qualquer **codigo hexadecimal de 6 digitos** manualmente no campo de texto
- Clique no **quadrado do seletor de cores** para abrir o espectro completo de cores

### Preenchimento Automatico Inteligente

- Lote vazio automaticamente vira `000`
- Serial vazio automaticamente vira `000001`
- Padding automatico com zeros a esquerda

### Logica Inteligente

- Material Generic seleciona fornecedor Generic automaticamente
- Material Creality seleciona fornecedor 0276 (Creality) automaticamente
- Lista de materiais e filtrada pelo fornecedor selecionado

## Hardware Suportado

### Hardware Recomendado (Links de Afiliados)

#### AliExpress (Internacional)
- **[Leitor RFID ACR122U](https://s.click.aliexpress.com/e/_ok8qAl9)** -- Leitor usado no desenvolvimento (compatibilidade garantida)
- **[Etiquetas MIFARE Classic 1K](https://s.click.aliexpress.com/e/_oBPVnEb)** -- Tags compativeis testadas no projeto

#### Mercado Livre (Brasil)
- **[Leitor e Gravador RFID](https://mercadolivre.com/sec/2QgqvkG)** -- Opcao nacional para compra do leitor RFID

### Leitores RFID Testados
- **ACR122U** (recomendado)
- **Outros leitores PC/SC** (compatibilidade nao garantida)

### Tags Suportadas
- **MIFARE Classic 1K**
- **MIFARE Classic 4K**
- **Tags Creality CFS**

## Desenvolvimento

### Estrutura do Projeto

```
cfs_spool/
├── main.go                 # Ponto de entrada do app Wails
├── app.go                  # Struct App com metodos vinculados ao Wails
├── app_options.go          # Handler de opcoes (materiais, fornecedores)
├── wails.json              # Configuracao do Wails
├── frontend/               # Frontend React + shadcn/ui
│   ├── src/                # Componentes e paginas React
│   ├── package.json        # Dependencias Node.js
│   ├── tailwind.config.js  # Configuracao Tailwind CSS
│   └── vite.config.ts      # Configuracao do bundler Vite
├── internal/
│   ├── creality/           # Logica especifica da Creality
│   │   ├── crypto.go       # Criptografia AES-ECB
│   │   └── fields.go       # Parsing e formatacao de campos
│   └── rfid/               # Comunicacao RFID
│       └── reader.go       # Interface PC/SC
├── build/                  # Saida de build e assets do Wails
├── tests/                  # Ferramentas de diagnostico
├── assets/                 # Recursos visuais
├── .github/workflows/      # Pipelines CI/CD
├── Makefile                # Atalhos de build
└── go.mod                  # Definicao do modulo Go
```

### Desenvolvimento Local

```bash
# Executar em modo de desenvolvimento com hot-reload
wails dev

# Executar testes Go
go test -v ./internal/... ./...

# Ferramentas de diagnostico RFID (requerem hardware)
go run tests/test_read_diagnosis.go
go run tests/test_auth_read.go
go run tests/test_decode_cfs.go
```

### Dependencias

- `github.com/wailsapp/wails/v2` -- Framework de aplicativo desktop
- `github.com/ebfe/scard` -- Interface PC/SC para comunicacao RFID
- `crypto/aes` -- Criptografia AES (biblioteca padrao do Go)
- React + shadcn/ui + Tailwind CSS (frontend)

## Referencia Tecnica

### Vendors Conhecidos

| **Vendor Code** | **Marca / Observacao** |
|:---:|:---:|
| 0x0276 | Creality, Hyper, Ender, HP (linhas oficiais) |
| 0xFFFF | Generico (qualquer fabricante nao-oficial) |

### Materials Conhecidos

| **Material Code** | **Descricao** |
|:---:|:---:|
| 00001 | Generic PLA |
| 00002 | Generic PLA-Silk |
| 00003 | Generic PETG |
| 00004 | Generic ABS |
| 00005 | Generic TPU |
| 00006 | Generic PLA-CF |
| 00007 | Generic ASA |
| 00008 | Generic PA |
| 00009 | Generic PA-CF |
| 00010 | Generic BVOH |
| 00011 | Generic PVA |
| 00012 | Generic HIPS |
| 00013 | Generic PET-CF |
| 00014 | Generic PETG-CF |
| 00015 | Generic PA6-CF |
| 00016 | Generic PAHT-CF |
| 00017 | Generic PPS |
| 00018 | Generic PPS-CF |
| 00019 | Generic PP |
| 00020 | Generic PET |
| 00021 | Generic PC |
| 01001 | Hyper PLA |
| 02001 | Hyper PLA-CF |
| 03001 | Hyper ABS |
| 04001 | CR-PLA |
| 05001 | CR-Silk |
| 06001 | CR-PETG |
| 06002 | Hyper PETG |
| 07001 | CR-ABS |
| 08001 | Ender-PLA |
| 09001 | EN-PLA+ |
| 09002 | Ender Fast PLA |
| 10001 | HP-TPU |
| 11001 | CR-Nylon |
| 13001 | CR-PLA Carbon |
| 14001 | CR-PLA Matte |
| 15001 | CR-PLA Fluo |
| 16001 | CR-TPU |
| 17001 | CR-Wood |
| 18001 | HP Ultra PLA |
| 19001 | HP-ASA |

### Formato da Tag CFS

O sistema Creality CFS armazena dados nos setores 1-2 das tags MIFARE Classic:

- **Setor 1 (Blocos 4-6)**: Dados criptografados do filamento
- **Criptografia**: AES-ECB com chaves derivadas do UID
- **Chave S1**: Derivada do UID usando chave `q3bu^t1nqfZ(pf$1`
- **Payload**: Descriptografado com chave `H@CFkRnz@KAtBJp2`

#### Layout dos Campos (38 bytes)

```
Date(5) + Supplier(4) + Batch(2) + Material(5) + Color(7) + Length(4) + Serial(6) + Reserve(4)
```

- Formato da cor: `"0" + 6 caracteres hex` (ex: `"077BB41"`)
- Batch padrao: `"A2"`, Reserve padrao: `"0000"`

#### Algoritmo de Autenticacao

1. **Tags novas**: Key A = `FFFFFFFFFFFF` (padrao MIFARE)
2. **Tags usadas**: Key A = derivada do UID usando algoritmo AES
3. **Fallback**: Multiplas tentativas com diferentes metodos

### Paleta de Cores Predefinidas

A interface inclui 35 cores predefinidas baseadas no sistema Creality:

| Categoria | Cores |
|---|---|
| **Azuis** | #25C4DA, #0099A7, #0B359A, #0A4AB6, #11B6EE, #90C6F5 |
| **Laranjas/Amarelos** | #FA7C0C, #F7B30F, #E5C20F, #B18F2E, #F8E911, #F6D311 |
| **Marrons** | #8D766D, #6C4E43 |
| **Vermelhos/Rosas** | #E62E2E, #EE2862, #EA2A2B, #E83D89, #AE2E65 |
| **Roxos** | #611C8B, #8D60C7, #B287C9 |
| **Verdes** | #006764, #018D80, #42B5AE, #1D822D, #54B351, #72E115 |
| **Cinzas** | #474747, #668798, #B1BEC6, #58636E |
| **Especiais** | #F2EFCE, #FFFFFF, #000000 |

## Releases e Versionamento

O tagueamento de versoes e automatico via `.github/workflows/auto-tag.yml`:

- **Auto-incremento**: Versao patch aumenta automaticamente no push para main
- **Versionamento Semantico**: Controle o tipo de versao usando flags na mensagem de commit:
  - `git commit -m "Mensagem #patch"` -- incrementa patch (v1.0.0 -> v1.0.1)
  - `git commit -m "Mensagem #minor"` -- incrementa minor (v1.0.0 -> v1.1.0)
  - `git commit -m "Mensagem #major"` -- incrementa major (v1.0.0 -> v2.0.0)
- **Acionamento manual**: Disponivel pela interface do GitHub Actions

Cada tag `v*` automaticamente compila binarios nativos para macOS, Linux e Windows via Wails.

## FAQ

### Como escolher cores personalizadas?

Voce tem total flexibilidade:
- Escolha uma das 35 cores predefinidas clicando na paleta
- Digite qualquer codigo hexadecimal manualmente no campo de texto (6 digitos)
- Clique no quadrado colorido para abrir o seletor de cores e escolher qualquer cor do espectro

### Campos opcionais nao funcionam?

Os campos **Lote** e **Serial** sao opcionais:
- Lote vazio automaticamente vira `000`
- Serial vazio automaticamente vira `000001`
- Preenchimento com zeros a esquerda automatico

### Como diagnosticar problemas de leitura?

```bash
go run tests/test_read_diagnosis.go
```

Este comando testa sistematicamente todos os metodos de autenticacao.

## Contribuicao

Contribuicoes sao bem-vindas! Por favor:

1. Faca um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudancas (`git commit -am 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## Licenca

Este projeto esta sob licenca MIT. Veja os detalhes em cada arquivo fonte.

## Disclaimer

Este projeto e desenvolvido para fins educacionais e de interoperabilidade. Nao e afiliado a Creality 3D Technology Co., Ltd.
