---
name: build-with-tests
description: Executa apenas a fase de build da fábrica (backend-builder + frontend-builder + test-verifier) quando o briefing já foi aprovado fora desta skill. Use para tocar a implementação de um spec pronto.
---

# Skill: build-with-tests

Subset da Fábrica de Software. Apenas a **fase de construção**: builders + verificação de aceitação. Pressupõe que **história e briefing já foram aprovados** fora desta skill.

## Quando usar

- O briefing técnico do `spec-writer` já existe e foi aprovado pelo humano.
- Você quer apenas tocar a implementação sem repassar pelos passos de planejamento.
- A feature é cross-stack (mexe na camada backend **e** na camada frontend, conforme `## Layers` do `CLAUDE.md`); para algo só de uma camada, invoque o builder correspondente direto.

## Quando NÃO usar

- Não há briefing aprovado → use `/feature-factory`.
- A história nunca foi formalizada → use `/feature-factory`.
- Mudança trivial (≤ 20 linhas em 1 arquivo) → faça direto.

## Argumentos

`/build-with-tests <referência ao briefing aprovado>`

A referência pode ser:
- Caminho de um arquivo de briefing salvo no repo.
- Bloco Markdown colado na conversa imediatamente antes.
- Link para um documento externo (Linear, Notion, Drive).

## Passos

### 0. Confirmação de pré-requisito

Antes de chamar qualquer agente, confirmar com o usuário:

> Qual é o briefing aprovado para esta feature?
> - Se for um arquivo: confirme o caminho.
> - Se for inline: confirme o trecho.
> - Se não houver briefing aprovado, **encerre esta skill** e sugira `/feature-factory`.

Não invocar agentes sem briefing referenciado.

### 1. Build do backend (`backend-builder`)

Antes de invocar, confirme que o `CLAUDE.md` da raiz tem um bloco `## Layers` declarando os globs de cada camada. Se não tiver, pare e peça ao usuário para declarar.

Invocar `backend-builder` com o briefing.

Esperar:
- Código backend dentro dos globs declarados como `backend:` em `## Layers`.
- Resumo do contrato de API (ou de bindings nativos) em Markdown.

Se o backend reportar problemas no briefing, **pare** — o spec precisa ser revisado pelo humano. Volte ao usuário.

### 2. Build do frontend (`frontend-builder`)

Invocar `frontend-builder` com:
- Resumo do contrato de API (do backend).
- Briefing.

Tratar divergência como na `feature-factory`: volte ao `backend-builder` se a UI não couber na API. Não aplique patch local.

### 3. Testes de aceitação (`test-verifier`)

Invocar `test-verifier`. Para cada teste que falhar, devolver ao builder indicado. Loop até verde.

### 4. Validação (opcional)

Perguntar ao usuário:

> Quer rodar o `implementation-validator` antes de fechar?
> - **Sim** (recomendado para features que tocam regras de negócio sensíveis).
> - **Não** — encerro com testes verdes e relatório.

Se sim, invocar `implementation-validator`. Tratar Crítico/Importante como na `feature-factory`. Se limpo, apresentar resumo final.

## Regras

- Esta skill **não substitui** os checkpoints da fábrica completa. Os 2 primeiros (história, briefing) tiveram que acontecer fora.
- Não passa por `codebase-researcher`, `story-writer` nem `spec-writer` — se algum desses for necessário, encerre e sugira `/feature-factory`.
- Reporte ao usuário **a cada transição** entre builders (com `arquivo:linha` quando referenciar código).
- Se em qualquer ponto descobrir que o briefing está ambíguo ou desatualizado, **pare**. Não tente reconstruir o spec aqui.
