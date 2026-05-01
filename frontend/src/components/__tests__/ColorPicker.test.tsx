import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { ColorPicker } from "@/components/ColorPicker";

describe("ColorPicker", () => {
  it("renderiza o label de Cor", () => {
    render(<ColorPicker value="000000" onChange={vi.fn()} />);
    expect(screen.getByText("Cor")).toBeInTheDocument();
  });

  it("exibe o valor hex no botão de gatilho do popover", () => {
    render(<ColorPicker value="FF4010" onChange={vi.fn()} />);
    expect(screen.getByText("#FF4010")).toBeInTheDocument();
  });

  it("aplica a cor como background do swatch (#RRGGBB)", () => {
    const { container } = render(
      <ColorPicker value="123456" onChange={vi.fn()} />,
    );
    // jsdom normaliza hex → rgb(); aceitamos as duas representações.
    const swatch = container.querySelector("div[style*='background']");
    expect(swatch).not.toBeNull();
    const style = swatch?.getAttribute("style") ?? "";
    // 0x12=18, 0x34=52, 0x56=86
    expect(style).toMatch(/#123456|rgb\(\s*18,\s*52,\s*86\s*\)/);
  });

  it("re-renderiza com novo valor quando a prop muda externamente", () => {
    const { rerender } = render(
      <ColorPicker value="000000" onChange={vi.fn()} />,
    );
    expect(screen.getByText("#000000")).toBeInTheDocument();

    rerender(<ColorPicker value="FFFFFF" onChange={vi.fn()} />);
    expect(screen.getByText("#FFFFFF")).toBeInTheDocument();
  });

  // Nota: a interatividade dentro do Popover do Radix (preset swatches,
  // input hex e gradient) usa portal/jsdom-incompatibilidades específicas
  // que tornam os testes frágeis. Cobrimos a sincronização externa →
  // estado interno aqui; a confirmação ponta-a-ponta fica em
  // verificação manual durante `wails dev`.
});
