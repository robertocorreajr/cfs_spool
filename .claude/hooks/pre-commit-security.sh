#!/usr/bin/env bash
#
# pre-commit-security.sh
# Hook do Claude Code (PreToolUse / Bash) que intercepta `git commit`
# e bloqueia se houver indícios de secrets ou arquivos sensíveis staged.
#
# Wired em .claude/settings.json. Recebe via stdin um JSON do tipo:
#   {"tool_name":"Bash","tool_input":{"command":"git commit -m ..."}}
#
# Exit codes:
#   0 → permite o commit
#   1 → bloqueia (Claude Code mostra stderr ao usuário)

set -euo pipefail

# ── Dependência: jq ─────────────────────────────────────────────────
if ! command -v jq >/dev/null 2>&1; then
  echo "[pre-commit-security] jq não encontrado; hook desativado." >&2
  echo "  Instale com: brew install jq" >&2
  exit 0
fi

# ── 1. Ler e parsear stdin ──────────────────────────────────────────
INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // empty')
CMD=$(echo "$INPUT" | jq -r '.tool_input.command // empty')

# ── 2. Só atuar em git commit ───────────────────────────────────────
if [[ "$TOOL_NAME" != "Bash" ]]; then
  exit 0
fi
if ! [[ "$CMD" =~ (^|[[:space:]\;\&\|])git[[:space:]]+commit([[:space:]]|$) ]]; then
  exit 0
fi

# ── 3. Coletar arquivos staged ──────────────────────────────────────
STAGED=$(git diff --cached --name-only 2>/dev/null || true)
if [[ -z "$STAGED" ]]; then
  exit 0
fi

# ── 4. Checar nomes de arquivo proibidos ────────────────────────────
BLOCK_ENV=$(echo "$STAGED" | grep -E '(^|/)\.env($|\.[^/]+$)' | grep -vE '\.(example|sample|template)$' || true)
BLOCK_KEYS=$(echo "$STAGED" | grep -iE '(credentials\.json$|service-account.*\.json$|id_(rsa|ed25519|ecdsa|dsa)$|\.(pem|key|p12|pfx|asc)$)' || true)

BLOCK_FILES=""
[[ -n "$BLOCK_ENV"  ]] && BLOCK_FILES+="$BLOCK_ENV"$'\n'
[[ -n "$BLOCK_KEYS" ]] && BLOCK_FILES+="$BLOCK_KEYS"$'\n'

if [[ -n "${BLOCK_FILES// /}" && "${BLOCK_FILES//[$'\n\t ']/}" != "" ]]; then
  echo "[pre-commit-security] BLOQUEADO: arquivos sensíveis no staging:" >&2
  echo "$BLOCK_FILES" | sed '/^$/d' | sed 's/^/  - /' >&2
  echo "" >&2
  echo "  Remova com: git restore --staged <arquivo>" >&2
  echo "  E adicione ao .gitignore antes de tentar de novo." >&2
  exit 1
fi

# ── 5. Checar conteúdo suspeito no diff staged ──────────────────────
DIFF=$(git diff --cached -U0 2>/dev/null || true)

# Cada padrão é checado separadamente para mensagens claras.
declare -a PATTERNS=(
  'AKIA[0-9A-Z]{16}'                                                   # AWS access key id
  'aws_secret_access_key[[:space:]]*=[[:space:]]*["'\''A-Za-z0-9/+=]'  # AWS secret literal
  '-----BEGIN (RSA|OPENSSH|EC|DSA|PGP) PRIVATE KEY-----'                # PEM private key
  'gh[pousr]_[A-Za-z0-9]{36,}'                                          # GitHub token
  'xox[abpr]-[A-Za-z0-9-]{10,}'                                         # Slack token
  '(api[_-]?key|password|secret|token)[[:space:]]*[:=][[:space:]]*['\''"][^'\''" ]{16,}['\''"]'  # generic
)

HITS=""
for pat in "${PATTERNS[@]}"; do
  match=$(echo "$DIFF" | grep -nE -- "$pat" || true)
  if [[ -n "$match" ]]; then
    HITS+="• padrão: $pat"$'\n'"$match"$'\n\n'
  fi
done

if [[ -n "$HITS" ]]; then
  echo "[pre-commit-security] BLOQUEADO: provável secret no diff staged." >&2
  echo "" >&2
  echo "$HITS" >&2
  echo "Se for falso positivo confirmado, refaça o commit fora do Claude Code" >&2
  echo "ou ajuste o padrão em .claude/hooks/pre-commit-security.sh." >&2
  exit 1
fi

exit 0
