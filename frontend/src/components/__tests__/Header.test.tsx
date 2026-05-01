import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
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
    // Apenas dois badges existem no design (UID + versão); ambos sumem com props vazias.
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
});
