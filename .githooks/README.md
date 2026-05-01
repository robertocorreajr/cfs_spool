# Git hooks versionados

Este diretório contém os hooks Git mantidos no repositório (em vez do
`.git/hooks/` local de cada clone).

## Instalação

Depois de clonar o repo, rode uma vez:

```bash
make install-hooks
```

Equivale a:

```bash
git config core.hooksPath .githooks
```

## Hooks disponíveis

### `pre-push`

Antes de qualquer push roda, em ordem:

1. `go build ./...`
2. `cd frontend && npx tsc --noEmit`
3. `go test ./...`
4. `cd frontend && npm test`

Se qualquer etapa falhar, o push é cancelado.

A mesma cadeia roda em CI no workflow `Tests`
(`.github/workflows/tests.yml`); o hook é a primeira linha de defesa
local.
