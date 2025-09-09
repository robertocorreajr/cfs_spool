# Sistema de Tagueamento Automático

Este documento descreve em detalhes o sistema de tagueamento automático implementado para o projeto CFS Spool, que utiliza GitHub Actions para criar tags automáticas seguindo o Versionamento Semântico (SemVer).

## Visão Geral

O sistema automatiza o processo de criação de versões (tags) do projeto usando um workflow do GitHub Actions que:

1. Monitora commits na branch principal (`main`)
2. Detecta palavras-chave específicas na mensagem de commit (`#patch`, `#minor`, `#major`)
3. Incrementa a versão de acordo com a palavra-chave encontrada
4. Cria uma nova tag com a versão incrementada
5. Aciona o workflow de build para gerar automaticamente uma nova release

## Versionamento Semântico

O sistema segue o padrão de [Versionamento Semântico 2.0.0](https://semver.org/lang/pt-BR/), onde:

- **Versão de Patch (z)**: Corrige bugs mantendo compatibilidade retroativa
- **Versão Minor (y)**: Adiciona funcionalidades mantendo compatibilidade retroativa
- **Versão Major (x)**: Realiza mudanças incompatíveis com versões anteriores

## Como Usar

### Incremento Automático via Mensagens de Commit

Para acionar o incremento de versão, adicione uma das seguintes hashtags ao final da mensagem de commit:

- `#patch`: Incrementa a versão de patch
  ```bash
  git commit -m "Corrige falha na leitura de RFID #patch"
  ```

- `#minor`: Incrementa a versão minor
  ```bash
  git commit -m "Adiciona suporte para novos tipos de RFID #minor"
  ```

- `#major`: Incrementa a versão major
  ```bash
  git commit -m "Refatora API de integração com novos leitores #major"
  ```

### Fluxo Recomendado de Trabalho

1. Desenvolva em branches de feature/bugfix
2. Ao concluir, crie um Pull Request para a branch `main`
3. Na mensagem de merge do PR, inclua o hashtag apropriado (#patch, #minor, #major)
4. Após o merge, o workflow será acionado automaticamente

### Acionamento Manual

Você também pode acionar o workflow manualmente através da interface do GitHub:

1. Acesse a aba "Actions" no repositório
2. Selecione o workflow "Auto Tag"
3. Clique em "Run workflow"
4. Selecione a branch `main` e o tipo de incremento desejado
5. Clique em "Run workflow"

## Como Funciona

O sistema consiste em dois workflows principais:

### 1. Workflow de Tagueamento Automático (auto-tag.yml)

Este workflow é responsável por:
- Monitorar commits na branch `main`
- Detectar palavras-chave nas mensagens de commit
- Calcular a nova versão
- Criar uma nova tag com a versão incrementada
- Acionar o workflow de build

### 2. Workflow de Build (build.yml)

Este workflow é responsável por:
- Monitorar a criação de novas tags
- Compilar o projeto
- Criar uma nova release no GitHub
- Anexar os artefatos compilados à release

## Resolução de Problemas

### Workflow Não Está Sendo Acionado

- Verifique se o commit foi feito na branch `main`
- Certifique-se de que a mensagem contém exatamente uma das hashtags: `#patch`, `#minor` ou `#major`
- Verifique os logs do workflow no GitHub Actions

### Tag Não Está Sendo Criada

- Verifique se o workflow tem permissões para criar tags
- Certifique-se de que não há tags com a mesma versão
- Verifique os logs do workflow para entender o motivo da falha

### Build Não Está Sendo Acionado

- Verifique se a tag foi criada corretamente
- Certifique-se de que o workflow de build está configurado para ser acionado pela criação de tags
- Verifique os logs do workflow de build

## Scripts de Utilitários

### Script de Teste do Tagueamento

O script `scripts/teste-auto-tag.sh` está disponível para:
- Verificar a tag atual
- Simular o incremento de versões
- Fornecer comandos para testar o sistema

Para usar o script:
```bash
./scripts/teste-auto-tag.sh
```

## Customização

Os workflows estão definidos nos arquivos:
- `.github/workflows/auto-tag.yml`
- `.github/workflows/build.yml`

Se necessário, você pode modificar esses arquivos para:
- Alterar os eventos que acionam os workflows
- Customizar o formato das tags
- Modificar o processo de build
- Adicionar notificações ou integrações
