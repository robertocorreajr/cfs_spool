# CFS Spool - Leitor/Gravador de Tags RFID Creality

🏷️ **Ferramenta para leitura e escrita de tags RFID do sistema Creality File System (CFS)**

## 📋 Descrição

O CFS Spool é uma ferramenta em linha de comando desenvolvida em Go para interagir com tags RFID MIFARE Classic utilizadas no sistema de filamentos da Creality. A ferramenta permite ler informações de bobinas de filamento como material, cor, lote, data de fabricação e outros metadados armazenados de forma criptografada nas tags.

## ✨ Funcionalidades

- 📖 **Leitura de tags CFS**: Decodifica informações completas do filamento
- 🔐 **Descriptografia AES-ECB**: Suporte ao sistema de criptografia Creality
- 🎨 **Interface amigável**: Apresentação clara das informações com emojis
- 🔧 **Modo debug**: Exibição de dados técnicos para desenvolvimento
- 🔄 **Autenticação robusta**: Múltiplos métodos de fallback para leitura

## 🚀 Instalação

### Pré-requisitos

- **Go 1.21+**
- **Leitor RFID compatível** (testado com ACR122U)
- **PC/SC Smart Card Daemon** (no macOS já vem instalado)

### Compilação

```bash
git clone https://github.com/robertocorreajr/cfs_spool.git
cd cfs_spool
go build -o cfs-spool ./cmd/cfs-spool
```

### Execução direta do Git

```bash
go run github.com/robertocorreajr/cfs_spool/cmd/cfs-spool@latest read-tag
```

## 📱 Uso

### Leitura de Tags

```bash
# Leitura básica
./cfs-spool read-tag

# Modo debug (dados técnicos)
./cfs-spool read-tag -debug
```

### Exemplo de Saída

```
╔══════════════════════════════════════════╗
║           INFORMAÇÕES DA TAG             ║
╚══════════════════════════════════════════╝
📦 Lote:        1A5
📅 Data:        20 de Janeiro de 2024
🏭 Fornecedor:  1B3D
🧪 Material:    CR-PLA (padrão)
🎨 Cor:         #77BB41 (hex)
📏 Comprimento: 330cm (1kg de filamento)
🔢 Serial:      000001
```

## 🛠️ Hardware Suportado

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
├── cmd/cfs-spool/          # Ponto de entrada da aplicação
│   ├── main.go             # CLI principal
│   └── write_tag.go        # Comandos de leitura/escrita
├── internal/
│   ├── creality/           # Lógica específica da Creality
│   │   ├── crypto.go       # Criptografia AES-ECB
│   │   └── fields.go       # Parsing e formatação de campos
│   └── rfid/               # Comunicação RFID
│       └── reader.go       # Interface PC/SC
├── go.mod                  # Dependências Go
└── README.md              # Esta documentação
```

### Dependências

- `github.com/ebfe/scard` - Interface PC/SC para comunicação RFID
- `crypto/aes` - Criptografia AES (biblioteca padrão)

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

## 🌐 Interface Web (Em Desenvolvimento)

Uma interface web está sendo desenvolvida para facilitar o uso da ferramenta:

- **Acesso local**: `http://localhost:8080`
- **Leitura via browser**: Interface amigável para leitura de tags
- **Compatibilidade**: Funciona com leitores RFID conectados via USB

## 🤝 Contribuição

Contribuições são bem-vindas! Por favor:

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## ⚠️ Disclaimer

Este projeto é desenvolvido para fins educacionais e de interoperabilidade. Não é afiliado à Creality 3D Technology Co., Ltd.
