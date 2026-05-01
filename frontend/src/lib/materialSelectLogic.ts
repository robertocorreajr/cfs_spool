// Lógica pura extraída de MaterialSelect.tsx para permitir cobertura
// unitária sem renderizar Radix Popover/Command (que dependem de medidas
// DOM reais e portais — instáveis em jsdom).

import type { MaterialOption } from "@/types/spool";

/** Devolve apenas os materiais cujo vendor bate com o fornecedor selecionado. */
export function filterByVendor(
  materials: MaterialOption[],
  vendor: string,
): MaterialOption[] {
  return materials.filter((m) => m.vendor === vendor);
}

/** Devolve o vendor associado a um material, ou undefined se não existir. */
export function pickAutoSupplier(
  materials: MaterialOption[],
  code: string,
): string | undefined {
  return materials.find((m) => m.code === code)?.vendor;
}

/**
 * Retorna o material atual quando ainda pertence ao novo vendor; caso
 * contrário devolve string vazia (UI deve limpar o campo). Se o material
 * atual não está na lista, considera "não pertence".
 */
export function clearMaterialIfVendorChanged(
  material: string,
  materials: MaterialOption[],
  newVendor: string,
): string {
  if (!material) return "";
  const found = materials.find((m) => m.code === material);
  if (!found) return "";
  return found.vendor === newVendor ? material : "";
}
