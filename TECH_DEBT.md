# Débitos Técnicos — Testes

Inventário de testes ausentes identificados em 2026-04-18. Referência para
planejamento da próxima versão. Nenhum destes itens altera comportamento de
produção — são apenas gaps de cobertura.

## Legenda de prioridade

- **P0** — lógica crítica que roda em produção sem cobertura; risco de
  regressão silenciosa no ciclo write → read de tags RFID.
- **P1** — cobertura parcial existente; faltam casos de borda relevantes.
- **P2** — código de dependência externa (PC/SC) ou UI; regressão detectável
  visualmente, mas o mock dá custo maior.

## Backend Go

### P0 — Criptografia e payload (internal/creality)

- [ ] `internal/creality/crypto.go:18` `DeriveS1KeyFromUID` — UID inválido,
      UIDs de 4/7/10 bytes.
- [ ] `internal/creality/crypto.go:31` `EncryptPayloadToBlocks` — padding,
      payloads com tamanho diferente de 48 bytes.
- [ ] `internal/creality/crypto.go:67` `DecryptBlocks` — hex inválido,
      tamanhos errados.
- [ ] `internal/creality/fields.go:102` `ParseFields` — campos truncados,
      tamanhos inválidos.
- [ ] `internal/creality/fields.go:45` `SetColor` — cores inválidas (não-hex,
      menos de 6 caracteres).

### P0 — Conversões e round-trip (app.go)

- [ ] `app.go:535` `convertLength` — valores custom (3000g, 5000g, 10000g),
      overflow em 65535 cm, fórmula `cm = gramas / 3`. Motivador original
      deste inventário.
- [ ] **Nova função** `DecodeLengthToGrams` em `internal/creality/fields.go` —
      inverso de `convertLength` para fechar o round-trip. Não existe hoje;
      `FormatLength` só exibe a string bruta ("03E8cm") para valores custom.
- [ ] Round-trip write→read completo: `convertLength` → `EncryptPayloadToBlocks`
      → `DecryptBlocks` → `ParseFields` → `DecodeLengthToGrams` preservando
      o valor original em gramas.
- [ ] `app.go:488` `convertMaterial` — mapeamentos e edge cases.
- [ ] `app.go:384` `vendorToSupplier`, `app.go:392` `materialToVendor`,
      `app.go:407` `vendorName` — todos os mapeamentos.

### P1 — Formatters e parsers (fields.go)

- [ ] `internal/creality/fields.go:169` `FormatDate` — datas com mês A/B/C
      (base 36).
- [ ] `internal/creality/fields.go:213` `FormatColor` — cores inválidas.
- [ ] `internal/creality/fields.go:221` `FormatLength` — comprimentos
      customizados (hoje exibe "03E8cm" cru, sem converter para gramas).
- [ ] `internal/creality/fields.go:234` `GetMaterialName` — códigos
      desconhecidos.
- [ ] `internal/creality/fields.go:334` `GetSupplierName` — códigos inválidos.
- [ ] `internal/creality/fields.go:347` `IsBlankTag` — heurísticas de batch,
      reserve e color.

### P1 — Bindings Wails (app.go)

- [ ] `app.go:197` `ReadTag` — precisa de mock PC/SC.
- [ ] `app.go:282` `WriteTag` — validação campo a campo + round-trip
      criptográfico.
- [ ] `app.go:360` `GetOptions` — vendors/materials/lengths não vazios e
      entrada "CUSTOM" presente em lengths.
- [ ] `app.go:162` `GetVersion` — formato semântico.
- [ ] `app.go:143` `handleTagPresent`, `app.go:156` `handleTagRemoved` —
      emissão de eventos Wails.
- [ ] `StartTagWatcher` / `StopTagWatcher` — goroutines e sinais PC/SC.

### P1 — Validações de options (app_options.go)

- [ ] `materials` — todos os códigos únicos.
- [ ] `vendors` — correspondência com `materials`.
- [ ] `lengths` — entrada "CUSTOM" presente.

### P2 — Reader PC/SC (internal/rfid/reader.go)

Pré-requisito: abstrair `scard` numa interface mockável antes de testar.

- [ ] `internal/rfid/reader.go:32` `Open` — mock de contexto PC/SC.
- [ ] `internal/rfid/reader.go:58` `UID` — mock de cartão conectado.
- [ ] `internal/rfid/reader.go:362` `ReadBlockDirect`, `:93` `WriteBlock`,
      `:404` `TryReadBlock` — APDUs felizes e de erro.
- [ ] `internal/rfid/reader.go:357` `transmit` — cobertura de erros APDU
      (6300, 6A82 etc.).

## Frontend TypeScript/React

### P0 — Infra (pré-requisito)

- [ ] Adicionar Vitest + `@testing-library/react` ao `frontend/package.json`.
- [ ] Criar `frontend/vitest.config.ts` com ambiente jsdom.
- [ ] Script `npm run test` no `frontend/package.json`.

### P1 — Componentes

- [ ] `frontend/src/components/LengthSelect.tsx:28` — input de custom grams,
      validação numérica, envio de `customGrams` para o handler pai, fluxo
      explícito para 3000g.
- [ ] `frontend/src/components/ColorPicker.tsx` — validação hex (rejeitar
      mais de 6 caracteres e não-hex), sincronização input ↔ gradient ↔
      swatches, 35 presets clicáveis.
- [ ] `frontend/src/components/MaterialSelect.tsx` — filtro por vendor,
      auto-seleção de supplier ao trocar material.
- [ ] `frontend/src/components/SpoolForm.tsx` — handlers de read/write,
      estados de loading, mensagens de erro, listeners de evento Wails.
- [ ] `frontend/src/components/Header.tsx` — exibição de versão vinda de
      `GetVersion()`.

## Infra de teste e CI

- [ ] Adicionar alvo `make test-frontend` no `Makefile`.
- [ ] Estender hook `pre-push` (`.git/hooks/pre-push`) para rodar testes de
      frontend.
- [ ] Rodar `make test` + `make test-frontend` no workflow antes da auto-tag
      em `.github/workflows/auto-tag.yml`.

## Priorização sugerida para a próxima versão

1. Round-trip de custom length (motivador original) — bloqueador de
   confiança em escrita de tags.
2. Crypto + `ParseFields` (P0) — fundação do payload.
3. Infra Vitest + `LengthSelect` e `ColorPicker` (P1).
4. Mock PC/SC e `reader.go` (P2) — escopo maior, deixar por último.
