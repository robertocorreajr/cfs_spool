import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";

// Mock dos bindings Wails — antes do import do componente.
const mockWriteTag = vi.fn();
const mockGetOptions = vi.fn();
const mockGetVersion = vi.fn();
// Retorna unsubscribe no-op; tipado para aceitar (event, callback).
const mockEventsOn = vi.fn(
  (_event: string, _cb: (...args: unknown[]) => void) => () => {},
);

vi.mock("../../../wailsjs/go/main/App", () => ({
  WriteTag: (req: unknown) => mockWriteTag(req),
  GetOptions: () => mockGetOptions(),
  GetVersion: () => mockGetVersion(),
}));

vi.mock("../../../wailsjs/runtime/runtime", () => ({
  EventsOn: (event: string, cb: (...args: unknown[]) => void) =>
    mockEventsOn(event, cb),
}));

// sonner injeta um portal global; mockamos para asserções determinísticas.
const toastErrors: string[] = [];
const toastSuccesses: string[] = [];
vi.mock("sonner", () => ({
  toast: {
    error: (msg: string) => toastErrors.push(msg),
    success: (msg: string) => toastSuccesses.push(msg),
    info: (msg: string) => toastSuccesses.push(msg),
  },
}));

import { SpoolForm } from "@/components/SpoolForm";

const optionsFixture = {
  materials: [
    { code: "04001", name: "CR-PLA", vendor: "0276" },
    { code: "E1001", name: "eSUN PLA+", vendor: "ESUN" },
  ],
  vendors: [
    { code: "0276", name: "Creality" },
    { code: "ESUN", name: "eSUN" },
  ],
  lengths: [
    { code: "0330", name: "330cm (1kg)", grams: "1000" },
    { code: "CUSTOM", name: "Personalizado", grams: "0" },
  ],
};

describe("SpoolForm", () => {
  beforeEach(() => {
    mockWriteTag.mockReset().mockResolvedValue(undefined);
    mockGetOptions.mockReset().mockResolvedValue(optionsFixture);
    mockGetVersion.mockReset().mockResolvedValue("v3.0.0");
    mockEventsOn.mockReset().mockReturnValue(() => {});
    toastErrors.length = 0;
    toastSuccesses.length = 0;
  });

  it("registra listeners para os eventos tag:status e tag:read", async () => {
    render(<SpoolForm />);
    await waitFor(() => {
      expect(mockEventsOn).toHaveBeenCalledWith("tag:status", expect.any(Function));
      expect(mockEventsOn).toHaveBeenCalledWith("tag:read", expect.any(Function));
    });
  });

  it("carrega options e versão na montagem", async () => {
    render(<SpoolForm />);
    await waitFor(() => {
      expect(mockGetOptions).toHaveBeenCalled();
      expect(mockGetVersion).toHaveBeenCalled();
    });
    // Versão acaba renderizada no Header.
    await waitFor(() => {
      expect(screen.getByText("v3.0.0")).toBeInTheDocument();
    });
  });

  it("bloqueia gravação quando material não está selecionado", async () => {
    render(<SpoolForm />);
    await waitFor(() => expect(mockGetOptions).toHaveBeenCalled());

    const writeButton = screen.getByRole("button", { name: /Gravar/ });
    fireEvent.click(writeButton);

    await waitFor(() => {
      expect(toastErrors).toContain("Selecione um material");
    });
    expect(mockWriteTag).not.toHaveBeenCalled();
  });

  it("renderiza barra de status 'Aguardando tag' por padrão", async () => {
    render(<SpoolForm />);
    await waitFor(() =>
      expect(screen.getByText(/Aguardando tag/)).toBeInTheDocument(),
    );
  });

  it("propaga toast de erro quando GetOptions falha", async () => {
    mockGetOptions.mockRejectedValueOnce(new Error("falha pcsc"));
    render(<SpoolForm />);
    await waitFor(() => {
      expect(toastErrors.length).toBeGreaterThan(0);
    });
    expect(toastErrors[0]).toMatch(/Erro ao carregar/);
  });
});
