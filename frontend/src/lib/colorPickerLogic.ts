// Lógica pura extraída de ColorPicker.tsx para permitir cobertura
// unitária sem renderizar Radix Popover (que abre em portal e depende
// de medidas DOM reais — instáveis em jsdom).

/** Filtra apenas dígitos hex e trunca em 6 caracteres. Não muda o case. */
export function cleanHexInput(raw: string): string {
  return raw.replace(/[^0-9A-Fa-f]/g, "").slice(0, 6);
}

/** Verdadeiro quando o hex está no formato exato de 6 chars válidos. */
export function isCompleteHex(hex: string): boolean {
  return /^[0-9A-Fa-f]{6}$/.test(hex);
}

/**
 * Decide se a prop externa `value` (controlada pelo pai) deve sobrescrever
 * o estado local `hexInput`. Espelha a heurística defensiva do componente:
 * só reescreve quando ambos têm 6 chars e diferem case-insensitive — evita
 * disputar foco com edição em curso.
 */
export function shouldSyncFromProp(value: string, hexInput: string): boolean {
  if (value === hexInput) return false;
  if (value.length !== 6 || hexInput.length !== 6) return false;
  return value.toUpperCase() !== hexInput.toUpperCase();
}
