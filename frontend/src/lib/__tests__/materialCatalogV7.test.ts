// materialCatalogV7.test.ts — testes de aceitação para a expansão do catálogo v7.
//
// Critério coberto:
//   1. Vendor eSUN + material "eSUN PLA-Basic" → filterByVendor retorna E1009 quando
//      a slice de materiais inclui a entrada expandida do catálogo.
//   2. pickAutoSupplier mapeia corretamente todos os novos codes dos 18 materiais.
//   3. clearMaterialIfVendorChanged preserva materiais novos para o vendor correto
//      e limpa quando o vendor muda.
//
// Não duplica materialSelectLogic.test.ts — este teste usa uma fixture que inclui
// especificamente as 18 novas entradas, exercendo os critérios da história aprovada.

import { describe, expect, it } from "vitest";
import {
  clearMaterialIfVendorChanged,
  filterByVendor,
  pickAutoSupplier,
} from "@/lib/materialSelectLogic";
import type { MaterialOption } from "@/types/spool";

// Subset representativo do catálogo v7 expandido: inclui todas as 18 novas entradas
// mais uma entrada preexistente de cada vendor para referência.
const catalogV7: MaterialOption[] = [
  // Genérico preexistente
  { code: "00001", name: "PLA", vendor: "0000" },
  // Genéricos novos
  { code: "00028", name: "Generic PA-GF", vendor: "0000" },
  { code: "00030", name: "Generic PP-GF", vendor: "0000" },
  // eSUN preexistente
  { code: "E1001", name: "eSUN PLA+", vendor: "ESUN" },
  // eSUN novos
  { code: "E1007", name: "eSUN PLA+HS", vendor: "ESUN" },
  { code: "E1008", name: "eSUN PLA-LW", vendor: "ESUN" },
  { code: "E1009", name: "eSUN PLA-Basic", vendor: "ESUN" },
  { code: "E2004", name: "eSUN PETG-CF", vendor: "ESUN" },
  { code: "E3002", name: "eSUN ABS-CF", vendor: "ESUN" },
  { code: "E3003", name: "eSUN ABS+HS", vendor: "ESUN" },
  { code: "E5001", name: "eSUN TPU-95A", vendor: "ESUN" },
  // Polymaker preexistente
  { code: "P1001", name: "Panchroma PLA Satin", vendor: "POLY" },
  // Polymaker novos
  { code: "P2001", name: "Fiberon PETG-ESD", vendor: "POLY" },
  { code: "P2002", name: "Fiberon PETG-rCF08", vendor: "POLY" },
  { code: "P7001", name: "Fiberon PA6-GF25", vendor: "POLY" },
  { code: "P7002", name: "Fiberon PA612-CF15", vendor: "POLY" },
  { code: "P7003", name: "Fiberon PA6-CF20", vendor: "POLY" },
  { code: "P7004", name: "Fiberon PA12-CF10", vendor: "POLY" },
  { code: "P7005", name: "Fiberon PA612-ESD", vendor: "POLY" },
  { code: "P8001", name: "Fiberon PET-CF17", vendor: "POLY" },
  { code: "P9001", name: "Fiberon PPS-CF10", vendor: "POLY" },
];

// Codes novos por vendor — usados nas asserções da tabela abaixo.
const novosEsun = ["E1007", "E1008", "E1009", "E2004", "E3002", "E3003", "E5001"];
const novosPoly = ["P2001", "P2002", "P7001", "P7002", "P7003", "P7004", "P7005", "P8001", "P9001"];
const novosGenerico = ["00028", "00030"];

describe("catálogo v7 — filterByVendor", () => {
  it("retorna eSUN PLA-Basic (E1009) quando vendor=ESUN", () => {
    const resultado = filterByVendor(catalogV7, "ESUN");
    const codes = resultado.map((m) => m.code);
    expect(codes).toContain("E1009");
    const entrada = resultado.find((m) => m.code === "E1009");
    expect(entrada?.name).toBe("eSUN PLA-Basic");
  });

  it("retorna todos os 7 novos materiais eSUN quando vendor=ESUN", () => {
    const resultado = filterByVendor(catalogV7, "ESUN");
    const codes = resultado.map((m) => m.code);
    for (const code of novosEsun) {
      expect(codes).toContain(code);
    }
  });

  it("retorna todos os 9 novos materiais Polymaker quando vendor=POLY", () => {
    const resultado = filterByVendor(catalogV7, "POLY");
    const codes = resultado.map((m) => m.code);
    for (const code of novosPoly) {
      expect(codes).toContain(code);
    }
  });

  it("retorna os 2 novos genéricos quando vendor=0000", () => {
    const resultado = filterByVendor(catalogV7, "0000");
    const codes = resultado.map((m) => m.code);
    for (const code of novosGenerico) {
      expect(codes).toContain(code);
    }
  });

  it("não inclui materiais de outros vendors no filtro ESUN", () => {
    const resultado = filterByVendor(catalogV7, "ESUN");
    const codes = resultado.map((m) => m.code);
    for (const code of [...novosPoly, ...novosGenerico]) {
      expect(codes).not.toContain(code);
    }
  });
});

describe("catálogo v7 — pickAutoSupplier", () => {
  // Testa que pickAutoSupplier retorna o vendor correto para cada novo code.
  const casos: Array<{ code: string; vendor: string }> = [
    { code: "E1009", vendor: "ESUN" },
    { code: "E1007", vendor: "ESUN" },
    { code: "E1008", vendor: "ESUN" },
    { code: "E2004", vendor: "ESUN" },
    { code: "E3002", vendor: "ESUN" },
    { code: "E3003", vendor: "ESUN" },
    { code: "E5001", vendor: "ESUN" },
    { code: "P2001", vendor: "POLY" },
    { code: "P2002", vendor: "POLY" },
    { code: "P7001", vendor: "POLY" },
    { code: "P7002", vendor: "POLY" },
    { code: "P7003", vendor: "POLY" },
    { code: "P7004", vendor: "POLY" },
    { code: "P7005", vendor: "POLY" },
    { code: "P8001", vendor: "POLY" },
    { code: "P9001", vendor: "POLY" },
    { code: "00028", vendor: "0000" },
    { code: "00030", vendor: "0000" },
  ];

  for (const { code, vendor } of casos) {
    it(`mapeia ${code} → vendor ${vendor}`, () => {
      expect(pickAutoSupplier(catalogV7, code)).toBe(vendor);
    });
  }
});

describe("catálogo v7 — clearMaterialIfVendorChanged", () => {
  it("preserva E1009 quando vendor permanece ESUN", () => {
    expect(clearMaterialIfVendorChanged("E1009", catalogV7, "ESUN")).toBe("E1009");
  });

  it("limpa E1009 quando vendor muda de ESUN para POLY", () => {
    expect(clearMaterialIfVendorChanged("E1009", catalogV7, "POLY")).toBe("");
  });

  it("preserva P9001 quando vendor permanece POLY", () => {
    expect(clearMaterialIfVendorChanged("P9001", catalogV7, "POLY")).toBe("P9001");
  });

  it("limpa P9001 quando vendor muda de POLY para 0000", () => {
    expect(clearMaterialIfVendorChanged("P9001", catalogV7, "0000")).toBe("");
  });
});
