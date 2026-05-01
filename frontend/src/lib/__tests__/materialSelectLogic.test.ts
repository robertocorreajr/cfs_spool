import { describe, expect, it } from "vitest";
import {
  clearMaterialIfVendorChanged,
  filterByVendor,
  pickAutoSupplier,
} from "@/lib/materialSelectLogic";
import type { MaterialOption } from "@/types/spool";

const mat: MaterialOption[] = [
  { code: "00001", name: "PLA", vendor: "0000" },
  { code: "04001", name: "CR-PLA", vendor: "0276" },
  { code: "05001", name: "CR-Silk", vendor: "0276" },
  { code: "E1001", name: "eSUN PLA+", vendor: "ESUN" },
  { code: "P1001", name: "Panchroma PLA Satin", vendor: "POLY" },
];

describe("filterByVendor", () => {
  it("retorna apenas materiais do vendor pedido", () => {
    expect(filterByVendor(mat, "0276").map((m) => m.code)).toEqual(["04001", "05001"]);
  });

  it("retorna lista vazia quando vendor não bate em nenhum item", () => {
    expect(filterByVendor(mat, "INEXISTENTE")).toEqual([]);
  });

  it("retorna lista vazia para vendor string vazia (estado inicial)", () => {
    expect(filterByVendor(mat, "")).toEqual([]);
  });

  it("não muta a lista original", () => {
    const cópia = [...mat];
    filterByVendor(mat, "ESUN");
    expect(mat).toEqual(cópia);
  });
});

describe("pickAutoSupplier", () => {
  it("devolve vendor do material existente", () => {
    expect(pickAutoSupplier(mat, "E1001")).toBe("ESUN");
    expect(pickAutoSupplier(mat, "04001")).toBe("0276");
  });

  it("devolve undefined para código inexistente", () => {
    expect(pickAutoSupplier(mat, "ZZZZZ")).toBeUndefined();
  });

  it("devolve undefined para código vazio", () => {
    expect(pickAutoSupplier(mat, "")).toBeUndefined();
  });
});

describe("clearMaterialIfVendorChanged", () => {
  it("preserva material quando ainda pertence ao novo vendor", () => {
    expect(clearMaterialIfVendorChanged("04001", mat, "0276")).toBe("04001");
  });

  it("limpa material quando vendor mudou", () => {
    expect(clearMaterialIfVendorChanged("04001", mat, "ESUN")).toBe("");
  });

  it("retorna vazio quando material atual é vazio (sem trabalho a fazer)", () => {
    expect(clearMaterialIfVendorChanged("", mat, "0276")).toBe("");
  });

  it("limpa material quando código não existe na lista", () => {
    expect(clearMaterialIfVendorChanged("ZZZZZ", mat, "0276")).toBe("");
  });
});
