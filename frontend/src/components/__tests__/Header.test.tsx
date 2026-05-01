import { describe, expect, it, vi } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";

// Mocks dos bindings — devem vir antes do import do componente. O Header
// agora chama triggerManualCheck quando o usuário clica em "Verificar
// atualizações"; isso depende dos bindings Wails que não existem em jsdom.
vi.mock("../../../wailsjs/go/main/App", () => ({
  CheckForUpdate: vi.fn(),
  IgnoreUpdateVersion: vi.fn(),
  OpenURL: vi.fn(),
}));
vi.mock("../../../wailsjs/runtime/runtime", () => ({
  EventsOn: () => () => {},
}));
vi.mock("sonner", () => ({
  toast: Object.assign(() => {}, {
    success: () => {},
    error: () => {},
    info: () => {},
  }),
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

  it("dispara onCheckForUpdate ao clicar no botão de atualizações", async () => {
    const onCheck = vi.fn().mockResolvedValue(undefined);
    render(<Header version="v3.0.0" uid="" onCheckForUpdate={onCheck} />);

    fireEvent.click(
      screen.getByRole("button", { name: /Verificar atualizações/i }),
    );

    await waitFor(() => {
      expect(onCheck).toHaveBeenCalledTimes(1);
    });
  });
});
