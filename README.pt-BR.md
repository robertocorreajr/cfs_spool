# CFS Spool - Leitor/Gravador de Tags RFID Creality

🏷️ **Sistema completo para leitura e gravação de tags RFID do Creality File System (CFS)**

[![English](https://img.shields.io/badge/lang-en-blue.svg)](README.md)
[![Portuguese](https://img.shields.io/badge/lang-pt--BR-green.svg)](README.pt-BR.md)

## 📋 Descrição

O CFS Spool é uma aplicação completa desenvolvida em Go que oferece tanto interface de linha de comando quanto interface web para interagir com tags RFID MIFARE Classic utilizadas no sistema de filamentos da Creality. A ferramenta permite ler e gravar informações de bobinas de filamento como material, cor, lote, data de fabricação e outros metadados armazenados de forma criptografada nas tags.

## ✨ Funcionalidades

### 🖥️ Interface Web
- 🎨 **Seletor avançado de cores**: Escolha entre 35 cores predefinidas ou use o seletor de cores para qualquer cor personalizada
- 🧠 **Lógica inteligente**: Auto-seleção de fornecedor baseado no material escolhido
- 📝 **Preenchimento automático**: Campos opcionais com padding automático
- 📖 **Leitura visual**: Preview das cores lidas das tags existentes
- 🔄 **Interface responsiva**: Funciona em desktop e mobile
- � **Criptografia/descriptografia AES-ECB**: Suporte completo ao sistema de criptografia Creality
- 🔄 **Autenticação robusta**: Múltiplos métodos de fallback para leitura

### 🛠️ Recursos Avançados
- 🎯 **Derivação de chaves**: Algoritmo completo baseado no UID da tag
- 🔒 **Compatibilidade**: Funciona com tags novas (FFFFFFFFFFFF) e usadas (chave derivada)
- 🧪 **Ferramentas de diagnóstico**: Suite completa para troubleshooting
- 📦 **Instaladores nativos**: DMG para macOS, AppImage para Linux, executável para Windows

## 🚀 Instalação

### 📥 Downloads Prontos (Recomendado)

Baixe a versão mais recente dos instaladores nativos:

**[⬇️ Releases - GitHub](https://github.com/robertocorreajr/cfs_spool/releases/latest)**

- 🍎 **macOS**: `CFS-Spool-macOS.dmg` (instalador drag-and-drop)
- 🐧 **Linux**: `CFS-Spool-Linux.AppImage` (portável)
- 🪟 **Windows**: `CFS-Spool-Windows.exe` (instalador)

### 🛠️ Compilação Manual

#### Pré-requisitos

- **Go 1.21+**
- **Leitor RFID compatível** (testado com ACR122U)
- **PC/SC Smart Card Daemon** 
  - macOS: já incluso
  - Linux: `sudo apt install pcscd libpcsclite-dev`
  - Windows: driver do leitor RFID

#### Compilação

```bash
git clone https://github.com/robertocorreajr/cfs_spool.git
cd cfs_spool

# Compilar aplicativo web
go build -o cfs-spool-app ./cmd/app
```

## 📱 Uso

### 🖥️ Interface Web (Recomendado)

1. **Executar aplicação**:
   ```bash
   ./cfs-spool-app
   # ou no Windows: CFS-Spool.exe
   ```

2. **Acessar interface**: Navegador abre automaticamente em `http://localhost:8080`

3. **Usar interface**:
   - **Aba "Ler Tag"**: Coloque tag no leitor e clique "Ler Tag"
   - **Aba "Gravar Tag"**: Preencha dados e clique "Gravar Tag"

#### 🎨 Recursos da Interface Web

- **Seleção de cores**:
  - 35 cores predefinidas com preview visual
  - Seletor de cores para escolher qualquer cor personalizada
  - Preview em tempo real da cor selecionada
- **Preenchimento inteligente**: 
  - Lote vazio → `000`
  - Serial vazio → `000001`
  - Auto-padding com zeros à esquerda
- **Lógica inteligente**:
  - Material Generic → Fornecedor Generic (automático)
  - Material Creality → Fornecedor 0276 (Creality)
  - Filtragem de materiais por fornecedor

### Exemplo de Saída

```
╔══════════════════════════════════════════╗
║           INFORMAÇÕES DA TAG             ║
╚══════════════════════════════════════════╝
📦 Lote:        1A5
📅 Data:        20 de Janeiro de 2024
🏭 Fornecedor:  0276 (Creality)
🧪 Material:    CR-PLA (padrão)
🎨 Cor:         #77BB41 (hex)
📏 Comprimento: 330cm (1kg de filamento)
🔢 Serial:      000001
```

## 🛠️ Hardware Suportado

### 🛒 Hardware Recomendado (Links de Afiliados)

#### AliExpress (Internacional)
- **🏷️ [Leitor RFID ACR122U](https://s.click.aliexpress.com/e/_ok8qAl9)** – Leitor usado no desenvolvimento (compatibilidade garantida)
- **📇 [Etiquetas MIFARE Classic 1K](https://s.click.aliexpress.com/e/_oBPVnEb)** – Tags compatíveis testadas no projeto

#### Mercado Livre (Brasil)
- **🏷️ [Leitor e Gravador RFID](https://mercadolivre.com/sec/2QgqvkG)** – Opção nacional para compra do leitor RFID

### Leitores RFID Testados
- **ACR122U** ✅ (recomendado)
- **Outros leitores PC/SC** (compatibilidade não garantida)

### Tags Suportadas
- **MIFARE Classic 1K** ✅
- **MIFARE Classic 4K** ✅
- **Tags Creality CFS** ✅

## 🔧 Desenvolvimento

### Estrutura do Projeto

```
cfs-spool/
├── cmd/
│   └── app/                # 🖥️ Aplicativo Web
│       └── main.go         # Servidor web com API REST
├── internal/
│   ├── creality/           # Lógica específica da Creality
│   │   ├── crypto.go       # Criptografia AES-ECB
│   │   └── fields.go       # Parsing e formatação de campos
│   └── rfid/               # Comunicação RFID
│       └── reader.go       # Interface PC/SC
├── web/                    # 🎨 Frontend da interface web
│   ├── index.html          # Interface HTML/CSS/JS
│   └── favicon.svg         # Ícone da aplicação
├── tests/                  # 🧪 Ferramentas de teste
│   ├── test_auth_read.go   # Teste de autenticação
│   ├── test_basic_read.go  # Teste de leitura básica
│   ├── test_decode_cfs.go  # Teste de decodificação
│   └── test_read_diagnosis.go # Diagnóstico completo
├── assets/                 # 🎨 Recursos visuais
│   ├── icons/              # Ícones para instaladores
│   └── dmg-background.svg  # Fundo do instalador macOS
├── .github/workflows/      # 🚀 CI/CD
│   └── build.yml           # Pipeline de build automático
├── scripts/                # 📦 Scripts de release
│   └── release.sh          # Script de empacotamento
└── Dockerfile              # 🐳 Container Docker
```

### API REST (Interface Web)

A interface web expõe uma API REST simples:

- `GET /api/status` - Status da aplicação
- `GET /api/options` - Opções para dropdowns (materiais, fornecedores, etc.)
- `POST /api/read-tag` - Leitura de tag RFID
- `POST /api/write` - Gravação de tag RFID

### Dependências

- `github.com/ebfe/scard` - Interface PC/SC para comunicação RFID
- `crypto/aes` - Criptografia AES (biblioteca padrão)
- Interface web nativa (sem dependências externas)

### 🧪 Ferramentas de Diagnóstico

```bash
# Diagnóstico completo de leitura RFID
go run tests/test_read_diagnosis.go

# Teste de autenticação
go run tests/test_auth_read.go

# Teste de decodificação CFS
go run tests/test_decode_cfs.go
```

## 📊 Referência Técnica

### Vendors conhecidos

| **Vendor Code** | **Marca / Observação**                             |
|:---------------:|:--------------------------------------------------:|
|  0x0276         | Creality • Hyper • Ender • HP (linhas oficiais)    |
|  0xFFFF         | Genérico (qualquer fabricante não-oficial)         |

### Materials conhecidos

| **Material Code** | **Descrição**         |
|:-----------------:|:---------------------:|
|  00001            | Generic PLA           |
|  00002            | Generic PLA-Silk      |
|  00003            | Generic PETG          |
|  00004            | Generic ABS           |
|  00005            | Generic TPU           |
|  00006            | Generic PLA-CF        |
|  00007            | Generic ASA           |
|  00008            | Generic PA            |
|  00009            | Generic PA-CF         |
|  00010            | Generic BVOH          |
|  00011            | Generic PVA           |
|  00012            | Generic HIPS          |
|  00013            | Generic PET-CF        |
|  00014            | Generic PETG-CF       |
|  00015            | Generic PA6-CF        |
|  00016            | Generic PAHT-CF       |
|  00017            | Generic PPS           |
|  00018            | Generic PPS-CF        |
|  00019            | Generic PP            |
|  00020            | Generic PET           |
|  00021            | Generic PC            |
|  01001            | Hyper PLA             |
|  02001            | Hyper PLA-CF          |
|  03001            | Hyper ABS             |
|  04001            | CR-PLA                |
|  05001            | CR-Silk               |
|  06001            | CR-PETG               |
|  06002            | Hyper PETG            |
|  07001            | CR-ABS                |
|  08001            | Ender-PLA             |
|  09001            | EN-PLA+               |
|  09002            | Ender Fast PLA        |
|  10001            | HP-TPU                |
|  11001            | CR-Nylon              |
|  13001            | CR-PLA Carbon         |
|  14001            | CR-PLA Matte          |
|  15001            | CR-PLA Fluo           |
|  16001            | CR-TPU                |
|  17001            | CR-Wood               |
|  18001            | HP Ultra PLA          |
|  19001            | HP-ASA                |

### Formato da Tag CFS

O sistema Creality CFS armazena dados nos setores 1-2 das tags MIFARE Classic:

- **Setor 1 (Blocos 4-6)**: Dados criptografados do filamento
- **Criptografia**: AES-ECB com chaves derivadas do UID
- **Chave S1**: Derivada do UID usando chave "q3bu^t1nqfZ(pf$1"
- **Payload**: Descriptografado com chave "H@CFkRnz@KAtBJp2"

#### Algoritmo de Autenticação

1. **Tags novas**: Key A = `FFFFFFFFFFFF` (padrão MIFARE)
2. **Tags usadas**: Key A = derivada do UID usando algoritmo AES
3. **Fallback**: Múltiplas tentativas com diferentes métodos

## 🎨 Paleta de Cores Predefinidas

A interface web inclui 35 cores predefinidas baseadas no sistema Creality:

| Categoria | Cores |
|-----------|-------|
| **Azuis** | #25C4DA, #0099A7, #0B359A, #0A4AB6, #11B6EE, #90C6F5 |
| **Laranjas/Amarelos** | #FA7C0C, #F7B30F, #E5C20F, #B18F2E, #F8E911, #F6D311 |
| **Marrons** | #8D766D, #6C4E43 |
| **Vermelhos/Rosas** | #E62E2E, #EE2862, #EA2A2B, #E83D89, #AE2E65 |
| **Roxos** | #611C8B, #8D60C7, #B287C9 |
| **Verdes** | #006764, #018D80, #42B5AE, #1D822D, #54B351, #72E115 |
| **Cinzas** | #474747, #668798, #B1BEC6, #58636E |
| **Especiais** | #F2EFCE, #FFFFFF, #000000 |

## 🚀 Releases e Versionamento

- **v2.2.0+**: Versão apenas web com seletor de cores avançado
- **v1.2.0+**: Interface web completa com paleta de cores
- **v1.1.1**: Correção crítica na derivação de chaves
- **v1.1.0**: Primeira versão com instaladores nativos
- **v1.0.x**: Versões CLI básicas (descontinuadas)

### 📦 Sistema de Tags e Build Automático

#### 🏷️ Tagueamento Automático de Versões

O projeto inclui tagueamento automático de versões quando código é enviado para a branch principal:

- 🔄 **Auto-incremento**: Versão patch aumenta automaticamente
- 🚀 **Versionamento Semântico**: Controle o tipo de versão usando flags na mensagem de commit:
  - `git commit -m "Mensagem #patch"` - incrementa patch (v1.0.0 → v1.0.1)
  - `git commit -m "Mensagem #minor"` - incrementa minor (v1.0.0 → v1.1.0)
  - `git commit -m "Mensagem #major"` - incrementa major (v1.0.0 → v2.0.0)
- ⚙️ **Acionamento Manual**: Disponível pela interface do GitHub Actions

#### 🏗️ Pipeline de Build Automático

Cada tag `v*` gera automaticamente:
- 🍎 Instalador DMG para macOS (com ícone customizado)
- 🐧 AppImage portável para Linux
- 🪟 Executável para Windows com installer
- 🐳 Imagem Docker multi-arquitetura

## ❓ FAQ

### Como usar a interface web?

- Inicie o aplicativo usando `./cfs-spool-app`
- Acesse http://localhost:8080 em seu navegador
- Use a interface intuitiva para leitura e gravação de tags

### Como escolher cores personalizadas?

Você tem total flexibilidade:
- Escolha uma das 35 cores predefinidas clicando na paleta abaixo do campo de texto
- Digite qualquer código hexadecimal manualmente no campo de texto (6 dígitos)
- Clique no quadrado colorido à esquerda do campo de texto para abrir o seletor de cores e escolher qualquer cor do espectro

### Campos opcionais não funcionam?

Os campos **Lote** e **Serial** são opcionais:
- Lote vazio → automaticamente `000`
- Serial vazio → automaticamente `000001`
- Preenchimento com zeros à esquerda automático

### Como diagnosticar problemas de leitura?

```bash
go run tests/test_read_diagnosis.go
```

Este comando testa sistematicamente todos os métodos de autenticação.

## 🤝 Contribuição

Contribuições são bem-vindas! Por favor:

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

### 🔧 Desenvolvimento Local

```bash
# Executar Aplicativo Web
go run cmd/app/main.go

# Executar Testes
go run tests/test_read_diagnosis.go
```

## 📄 Licença

Este projeto está sob licença MIT. Veja os detalhes em cada arquivo fonte.

## ⚠️ Disclaimer

Este projeto é desenvolvido para fins educacionais e de interoperabilidade. Não é afiliado à Creality 3D Technology Co., Ltd.

---

**🏷️ CFS Spool v2.1.26** - Sistema completo para tags RFID Creality  
*Desenvolvido com ❤️ em Go*
