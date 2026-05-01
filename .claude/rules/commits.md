# Padrão de Commits

Use Conventional Commits em português (pt-BR):

- feat: adiciona nova funcionalidade
- fix: corrige um bug
- docs: atualiza documentação
- refactor: refatora sem mudar comportamento
- test: adiciona ou corrige testes
- chore: tarefas de manutenção

Exemplo: `feat: adiciona autenticação por OAuth`

## Marcas de skip de CI

Dois workflows estão configurados:

- **Tests** (`.github/workflows/tests.yml`) — roda em todos os pushes/PRs e valida `make test` + `make test-frontend`. Não cria nada.
- **Auto Tag** (`.github/workflows/auto-tag.yml`) — só roda em `main`; cria tag patch e dispara release.

Convenção:

- `[skip release]` no título do commit → pula apenas Auto Tag. Tests continua rodando. **Use sempre que o commit não justificar release** (test, refactor sem mudança comportamental, docs com testes que valem).
- `[skip ci]` → pula tudo (interpretado pelo GitHub Actions). Use só em docs puros (README typo, screenshot) onde nem teste precisa rodar.
- Sem marca → Tests + Auto Tag rodam (release nova será criada).

Exemplos:

- `test: adiciona cobertura de fields.go [skip release]`
- `refactor: extrai lógica pura [skip release]`
- `feat: adiciona check de versão` (sem marca — merece release)
- `docs: corrige typo no README [skip ci]` (raro)
