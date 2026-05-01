import { describe, expect, it, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { LengthSelect } from "@/components/LengthSelect";
import type { LengthOption } from "@/types/spool";

const mockLengths: LengthOption[] = [
  { code: "0083", name: "83cm (250g)", grams: "250" },
  { code: "0330", name: "330cm (1kg)", grams: "1000" },
  { code: "CUSTOM", name: "Personalizado", grams: "0" },
];

describe("LengthSelect", () => {
  it("renderiza o label de Comprimento", () => {
    render(
      <LengthSelect
        length="0330"
        customGrams=""
        onLengthChange={vi.fn()}
        onCustomGramsChange={vi.fn()}
        lengths={mockLengths}
      />,
    );
    expect(screen.getByText("Comprimento")).toBeInTheDocument();
  });

  it("não exibe input de gramas quando length não é CUSTOM", () => {
    render(
      <LengthSelect
        length="0330"
        customGrams=""
        onLengthChange={vi.fn()}
        onCustomGramsChange={vi.fn()}
        lengths={mockLengths}
      />,
    );
    expect(screen.queryByPlaceholderText(/Gramas/)).not.toBeInTheDocument();
  });

  it("exibe input de gramas quando length é CUSTOM", () => {
    render(
      <LengthSelect
        length="CUSTOM"
        customGrams=""
        onLengthChange={vi.fn()}
        onCustomGramsChange={vi.fn()}
        lengths={mockLengths}
      />,
    );
    const input = screen.getByPlaceholderText(/Gramas/);
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute("type", "number");
  });

  it("dispara onCustomGramsChange ao digitar gramatura customizada (fluxo 3000g)", () => {
    const onCustomGramsChange = vi.fn();
    render(
      <LengthSelect
        length="CUSTOM"
        customGrams=""
        onLengthChange={vi.fn()}
        onCustomGramsChange={onCustomGramsChange}
        lengths={mockLengths}
      />,
    );
    const input = screen.getByPlaceholderText(/Gramas/);
    fireEvent.change(input, { target: { value: "3000" } });
    expect(onCustomGramsChange).toHaveBeenCalledWith("3000");
  });

  it("reflete o valor de customGrams no input controlado", () => {
    render(
      <LengthSelect
        length="CUSTOM"
        customGrams="750"
        onLengthChange={vi.fn()}
        onCustomGramsChange={vi.fn()}
        lengths={mockLengths}
      />,
    );
    expect(screen.getByPlaceholderText(/Gramas/)).toHaveValue(750);
  });
});
