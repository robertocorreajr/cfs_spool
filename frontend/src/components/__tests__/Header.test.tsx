import { describe, expect, it, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

const mockEventsEmit = vi.fn();

// Mocks dos bindings — devem vir antes do import do componente.
vi.mock("../../../wailsjs/runtime/runtime", () => ({
  EventsOn: () => () => {},
  EventsEmit: (...args: unknown[]) => mockEventsEmit(...args),
}));

import { Header } from "@/components/Header";

describe("Header", () => {
  it("renderiza o título do app", () => {
    render(<Header version="" uid="" />);
    expect(screen.getByText("CFS Spool")).toBeInTheDocument();
  });

  it("exibe a versão recebida via prop", () => {
    render(<Header version="v3.0.0" uid="" />);
    expect(screen.getByText("v3.0.0")).toBeInTheDocument();
  });

  it("não renderiza o badge de versão quando version é string vazia", () => {
    const { container } = render(<Header version="" uid="" />);
    expect(container.querySelectorAll("[class*=Badge]").length).toBe(0);
  });

  it("exibe o UID quando preenchido", () => {
    render(<Header version="v3.0.0" uid="DEADBEEF" />);
    expect(screen.getByText(/UID/)).toBeInTheDocument();
    expect(screen.getByText("DEADBEEF")).toBeInTheDocument();
  });

  it("não exibe UID quando vazio", () => {
    render(<Header version="v3.0.0" uid="" />);
    expect(screen.queryByText(/UID/)).not.toBeInTheDocument();
  });

  it("não renderiza o ícone de atualização quando hasUpdate é false", () => {
    render(<Header version="v3.0.0" uid="" hasUpdate={false} />);
    expect(
      screen.queryByRole("button", { name: /Atualização disponível/i }),
    ).not.toBeInTheDocument();
  });

  it("renderiza ícone piscando quando hasUpdate é true", () => {
    render(<Header version="v3.0.0" uid="" hasUpdate />);
    expect(
      screen.getByRole("button", { name: /Atualização disponível/i }),
    ).toBeInTheDocument();
  });

  it("clicar no ícone emite update:show para abrir o modal", () => {
    mockEventsEmit.mockClear();
    render(<Header version="v3.0.0" uid="" hasUpdate />);

    fireEvent.click(
      screen.getByRole("button", { name: /Atualização disponível/i }),
    );

    expect(mockEventsEmit).toHaveBeenCalledWith("update:show");
  });
});
