import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MaterialSelect } from "@/components/MaterialSelect";
import type { MaterialOption, VendorOption } from "@/types/spool";

const vendors: VendorOption[] = [
  { code: "0276", name: "Creality" },
  { code: "0000", name: "Genérico" },
  { code: "ESUN", name: "eSUN" },
];

const materials: MaterialOption[] = [
  { code: "00001", name: "PLA", vendor: "0000" },
  { code: "04001", name: "CR-PLA", vendor: "0276" },
  { code: "E1001", name: "eSUN PLA+", vendor: "ESUN" },
];

describe("MaterialSelect", () => {
  it("renderiza os labels Fornecedor e Material", () => {
    render(
      <MaterialSelect
        supplier="0276"
        material="04001"
        onSupplierChange={vi.fn()}
        onMaterialChange={vi.fn()}
        materials={materials}
        vendors={vendors}
      />,
    );
    expect(screen.getByText("Fornecedor")).toBeInTheDocument();
    expect(screen.getByText("Material")).toBeInTheDocument();
  });

  it("exibe o nome do vendor selecionado no botão", () => {
    render(
      <MaterialSelect
        supplier="ESUN"
        material=""
        onSupplierChange={vi.fn()}
        onMaterialChange={vi.fn()}
        materials={materials}
        vendors={vendors}
      />,
    );
    expect(screen.getByText("eSUN")).toBeInTheDocument();
  });

  it("exibe o nome do material selecionado no botão", () => {
    render(
      <MaterialSelect
        supplier="0276"
        material="04001"
        onSupplierChange={vi.fn()}
        onMaterialChange={vi.fn()}
        materials={materials}
        vendors={vendors}
      />,
    );
    expect(screen.getByText("CR-PLA")).toBeInTheDocument();
  });

  it("usa placeholder quando vendor não está selecionado", () => {
    render(
      <MaterialSelect
        supplier=""
        material=""
        onSupplierChange={vi.fn()}
        onMaterialChange={vi.fn()}
        materials={materials}
        vendors={vendors}
      />,
    );
    expect(screen.getByText("Buscar fornecedor...")).toBeInTheDocument();
    expect(screen.getByText("Buscar material...")).toBeInTheDocument();
  });

  // Nota: a filtragem do dropdown e a auto-seleção de supplier ao escolher
  // um material rodam dentro do Popover/Command (cmdk) do Radix, que abre
  // em portal e depende de medidas do DOM real para o virtualizer. Os
  // testes acima cobrem o estado renderizado; o handler interno
  // handleMaterialChange/handleSupplierChange é exercido manualmente em
  // `wails dev` e indiretamente pelo SpoolForm.
});
