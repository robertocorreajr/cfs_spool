---
name: feature-factory
description: Orquestra a cadeia dos 7 agentes para construir uma feature completa com 3 checkpoints humanos (história, briefing, PR). Use quando o usuário pedir uma feature nova de tamanho médio/grande.
---

# Skill: feature-factory

Orquestrador da Fábrica de Software. Conduz uma feature pela cadeia dos 7 agentes especializados, **respeitando os 3 checkpoints humanos**.

## Quando usar

- Pedidos de feature nova de tamanho **médio ou grande** (mexe em modelo + API/binding + UI, ou em mais de 2 arquivos).
- Pedidos que cruzam fronteiras (camada backend ↔ camada frontend, conforme `## Layers` do `CLAUDE.md`).
- Sempre que houver risco de ambiguidade no problema (o que vai surgir como pergunta no `story-writer`).

## Quando NÃO usar

- Typos, ajustes locais, mudanças de configuração curtas.
- Mudanças contidas em **≤ 20 linhas em 1 arquivo**.
- Operações de manutenção (rename, mover arquivo, atualizar dependência).

Nesses casos, faça direto na conversa. A fábrica tem overhead — vale quando o trabalho passa do trivial.

## Argumentos

`/feature-factory <descrição rascunho da feature>`

A descrição pode ser curta — o `story-writer` vai estruturar.

## Passos obrigatórios (na ordem)

### 1. Pesquisa (`codebase-researcher`)

Invocar `codebase-researcher` com a descrição rascunho.

Aguardar relatório com: arquivos relevantes, padrões existentes, pontos de reuso, riscos, perguntas em aberto.

Não prosseguir sem o relatório.

### 2. História (`story-writer`) — **Checkpoint #1**

Invocar `story-writer` com:
- Descrição rascunho do usuário.
- Relatório do `codebase-researcher`.

Apresentar a história gerada ao usuário e **PARAR**:

> ⏸ **Checkpoint #1** — Revise a história acima.
> - Se estiver tudo certo, diga "aprovado".
> - Se faltar algo ou estiver errado, peça ajustes (a história será regerada — não vamos fazer patch).

Não chame o próximo agente antes da aprovação explícita.

### 3. Briefing (`spec-writer`) — **Checkpoint #2**

Após aprovação da história, invocar `spec-writer` com:
- História aprovada.
- Relatório do `codebase-researcher`.

Apresentar o briefing técnico ao usuário e **PARAR**:

> ⏸ **Checkpoint #2** — Revise o briefing técnico.
> - Procure sinais de alerta destacados em negrito (ex.: "guardar IDs em memória", falta de tenant check).
> - Aprovado? Confirme. Se não, peça ajustes — o briefing será regerado.

### 4. Build do backend (`backend-builder`)

Após aprovação do briefing, invocar `backend-builder` com o briefing.

Antes de invocar, confirme que o `CLAUDE.md` da raiz tem um bloco `## Layers` declarando os globs de cada camada. Se não tiver, pare e peça ao usuário para declarar (sem isso, os builders não sabem o escopo de escrita).

Esperar como saída:
- Código backend dentro dos globs declarados como `backend:` em `## Layers` (rotas/handlers/bindings, serviços, DB, migrations, jobs, testes unitários).
- **Resumo do contrato de API** (ou de bindings nativos) em Markdown.

Se o backend reportar que o briefing está ambíguo ou impossível, pare e volte ao `spec-writer`. Não improvise.

### 5. Build do frontend (`frontend-builder`)

Invocar `frontend-builder` passando:
- Resumo do contrato de API do backend (essencial — entrada principal).
- Briefing do `spec-writer`.

Tratar feedback de divergência:
- Se o frontend reportar que o formato da API não cabe na UI, **NÃO** aplique patch local.
- Volte ao `backend-builder` com a divergência. Repita 4→5 até alinhar.

### 6. Testes de aceitação (`test-verifier`)

Invocar `test-verifier`. Esperar saída de execução.

Para cada teste que falhar:
- Identificar qual builder o `test-verifier` indicou.
- Devolver ao builder correspondente para correção.
- Re-rodar 6.

Não prosseguir enquanto algum teste de aceitação estiver falhando.

### 7. Validação (`implementation-validator`) — **Checkpoint #3**

Invocar `implementation-validator`. Ler o relatório classificado por severidade.

- **Crítico** ou **Importante** → devolver ao builder certo (`backend-builder` ou `frontend-builder`). Voltar a 4/5.
- **Apenas Menor** ou limpo → apresentar resumo final ao usuário e **PARAR**:

> ⏸ **Checkpoint #3** — `implementation-validator` limpo.
> - Revise o `git diff` completo.
> - Quando estiver pronto, abra o PR: `gh pr create --base main --head <branch>`.

## Regras

- **Nunca pule um checkpoint.** Eles são o ponto onde o julgamento humano importa.
- **Cada agente recebe contexto limpo.** Não acumule histórico entre invocações; passe explicitamente o que é entrada (relatório, história, briefing, contrato de API).
- Quando reportar progresso, **cite `arquivo:linha`** sempre que mencionar código.
- Se em qualquer ponto o usuário pedir para abortar a fábrica, pare imediatamente. Não tente "salvar" o trabalho parcial — registre o que foi feito e devolva o controle.
- Se uma suposição arquitetural se mostrar errada no meio do caminho, **jogue a conversa fora** e recomece do passo apropriado (provavelmente história ou briefing). Não corrija com patch — drift é o assassino silencioso.
