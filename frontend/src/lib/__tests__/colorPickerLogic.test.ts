import { describe, expect, it } from "vitest";
import {
  cleanHexInput,
  isCompleteHex,
  shouldSyncFromProp,
} from "@/lib/colorPickerLogic";

describe("cleanHexInput", () => {
  it("preserva hex válido", () => {
    expect(cleanHexInput("FF4010")).toBe("FF4010");
    expect(cleanHexInput("ff4010")).toBe("ff4010");
  });

  it("remove caracteres não-hex", () => {
    expect(cleanHexInput("FF#4010")).toBe("FF4010");
    expect(cleanHexInput("F G H 1 2 3")).toBe("F123");
  });

  it("trunca em 6 chars (preserva os 6 primeiros válidos)", () => {
    expect(cleanHexInput("FFAABB1234")).toBe("FFAABB");
  });

  it("remove # e espaços", () => {
    expect(cleanHexInput(" #FF4010 ")).toBe("FF4010");
  });

  it("retorna vazio para entrada totalmente inválida", () => {
    expect(cleanHexInput("ZZZZZZ")).toBe("");
    expect(cleanHexInput("")).toBe("");
  });
});

describe("isCompleteHex", () => {
  it("verdadeiro para 6 chars hex válidos", () => {
    expect(isCompleteHex("FF4010")).toBe(true);
    expect(isCompleteHex("000000")).toBe(true);
    expect(isCompleteHex("aabbcc")).toBe(true);
  });

  it("falso para tamanho diferente de 6", () => {
    expect(isCompleteHex("FF40")).toBe(false);
    expect(isCompleteHex("FF40100")).toBe(false);
    expect(isCompleteHex("")).toBe(false);
  });

  it("falso quando contém char não-hex", () => {
    expect(isCompleteHex("ZZ4010")).toBe(false);
    expect(isCompleteHex("FF#010")).toBe(false);
  });
});

describe("shouldSyncFromProp", () => {
  it("não sincroniza quando valor é igual", () => {
    expect(shouldSyncFromProp("FF4010", "FF4010")).toBe(false);
  });

  it("não sincroniza quando algum lado é incompleto", () => {
    expect(shouldSyncFromProp("FF4010", "FF40")).toBe(false);
    expect(shouldSyncFromProp("", "FF4010")).toBe(false);
  });

  it("não sincroniza quando diferem só no case (mesma cor)", () => {
    expect(shouldSyncFromProp("FF4010", "ff4010")).toBe(false);
  });

  it("sincroniza quando ambos completos e a cor difere", () => {
    expect(shouldSyncFromProp("FF4010", "00FF00")).toBe(true);
  });
});
