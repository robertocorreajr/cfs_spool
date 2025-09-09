# Guia Rápido: Sistema de Tagueamento Automático

Este documento é um guia de referência rápida para o sistema de tagueamento automático implementado no projeto CFS Spool.

## Como Funciona

O sistema utiliza workflows do GitHub Actions para detectar mensagens de commit com os sufixos especiais e criar tags automaticamente seguindo o versionamento semântico (SemVer).

## Tipos de Incremento

Ao fazer um commit na branch `main`, você pode adicionar um dos seguintes sufixos à mensagem:

- **`#patch`**: Incrementa a versão de patch (z) em x.y.z
  - *Exemplo*: `"Corrige bug na interface #patch"`
  - *Resultado*: v1.0.0 → v1.0.1

- **`#minor`**: Incrementa a versão minor (y) em x.y.z
  - *Exemplo*: `"Adiciona novo recurso #minor"`
  - *Resultado*: v1.0.0 → v1.1.0

- **`#major`**: Incrementa a versão major (x) em x.y.z
  - *Exemplo*: `"Muda API incompatível com versão anterior #major"`
  - *Resultado*: v1.0.0 → v2.0.0

## Fluxo de Trabalho

1. Faça suas alterações em uma branch de feature/fix
2. Ao finalizar, mescle a branch com `main`
3. Na mensagem do commit de merge, adicione o sufixo apropriado (#patch, #minor, #major)
4. O workflow irá:
   - Detectar o sufixo na mensagem
   - Incrementar a versão automaticamente
   - Criar uma nova tag
   - Acionar o workflow de build e release

## Verificação Manual

Para verificar se o sistema está funcionando corretamente:

1. Execute o script `./scripts/teste-auto-tag.sh`
2. Verifique a última tag criada
3. Siga os passos indicados para testar o workflow

## Acionamento Manual

Você também pode acionar o workflow manualmente:

1. Acesse a aba Actions no GitHub
2. Clique em "Auto Tag"
3. Clique em "Run workflow"
4. Selecione o tipo de incremento desejado
5. Clique em "Run workflow"

## Solução de Problemas

Se o sistema não estiver criando tags automaticamente:

1. Verifique se o commit foi feito na branch `main`
2. Confira se a mensagem contém o sufixo correto (#patch, #minor, #major)
3. Verifique os logs do workflow na aba Actions do GitHub
4. Certifique-se de que o workflow tem permissões para criar tags
