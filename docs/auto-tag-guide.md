# Guia de Tagueamento Automático de Versões

Este documento explica como funciona o sistema automático de tagueamento semântico do projeto CFS Spool.

## O que é Versionamento Semântico?

O projeto segue o padrão [SemVer (Semantic Versioning)](https://semver.org/lang/pt-BR/), onde as versões têm o formato `vMAJOR.MINOR.PATCH`:

- **MAJOR**: alterações incompatíveis com versões anteriores
- **MINOR**: adições de funcionalidades compatíveis com versões anteriores
- **PATCH**: correções de bugs compatíveis com versões anteriores

## Sistema Automático de Tags

O projeto possui um sistema de GitHub Actions que automaticamente cria tags de versão quando código é enviado para a branch principal (`main`).

### Como Funciona

1. Quando código é enviado para a branch `main`, o workflow `auto-tag.yml` é acionado
2. O workflow analisa a mensagem do commit mais recente para decidir o tipo de incremento de versão
3. Uma nova tag é criada e enviada ao repositório com a versão incrementada

### Controlando o Tipo de Incremento

Por padrão, a versão **patch** é incrementada. Para controlar explicitamente o tipo de incremento, adicione uma das flags abaixo na mensagem de commit:

| Flag | Incremento | Exemplo |
|------|------------|---------|
| `#patch` | Incrementa o número de patch | `v1.0.0` → `v1.0.1` |
| `#minor` | Incrementa o número de minor e reseta o patch | `v1.0.1` → `v1.1.0` |
| `#major` | Incrementa o número de major e reseta minor e patch | `v1.1.0` → `v2.0.0` |

#### Exemplos:

```bash
# Incrementar patch (padrão)
git commit -m "Corrige problema no seletor de cores"

# Incrementar patch (explícito)
git commit -m "Corrige problema no seletor de cores #patch"

# Incrementar minor
git commit -m "Adiciona novo seletor de cores #minor"

# Incrementar major
git commit -m "Refatora API completamente #major"
```

### Acionamento Manual

Também é possível acionar o workflow de tagueamento manualmente através da interface do GitHub:

1. Acesse a aba "Actions" no GitHub
2. Selecione o workflow "Auto Tag"
3. Clique em "Run workflow"
4. Escolha o tipo de incremento: `patch`, `minor`, ou `major`
5. Confirme clicando em "Run workflow"

## Pipeline de Build

Uma vez criada a tag, o workflow `build.yml` é automaticamente acionado para:

1. Executar testes automatizados
2. Compilar binários para todas as plataformas
3. Criar instaladores nativos (DMG, AppImage, EXE)
4. Publicar uma nova release no GitHub

## Dicas para Mensagens de Commit

- **Correções de bugs**: use `#patch` ou deixe o padrão
- **Novas funcionalidades**: use `#minor`
- **Mudanças grandes/incompatíveis**: use `#major`

## Fluxo Recomendado

1. Desenvolva em branches de feature
2. Crie Pull Requests para a branch `main`
3. Após revisar e aprovar o PR, faça merge para `main`
4. O sistema de tagueamento automático criará a nova versão
5. A pipeline de CI/CD criará automaticamente os instaladores para a nova versão

## Resolução de Problemas

Se o tagueamento automático não for acionado ou falhar:

1. Verifique se o commit foi feito diretamente na branch `main`
2. Certifique-se que as GitHub Actions estão habilitadas no repositório
3. Verifique os logs na aba "Actions" no GitHub
4. Use o tagueamento manual através da interface do GitHub como alternativa
