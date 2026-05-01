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

  it("filtra dígitos não numéricos e trunca em 6 chars no campo Serial", async () => {
    render(<SpoolForm />);
    await waitFor(() => expect(mockGetOptions).toHaveBeenCalled());

    const serialInput = screen.getByPlaceholderText("000001") as HTMLInputElement;
    // Reseta para vazio antes (default é "000001")
    fireEvent.change(serialInput, { target: { value: "" } });
    fireEvent.change(serialInput, { target: { value: "abc123def456789" } });
    expect(serialInput.value).toBe("123456");
  });

  it("aplica dados de tag normal lida via callback de tag:read", async () => {
    let tagReadCallback: ((data: unknown) => void) | undefined;
    mockEventsOn.mockImplementation((event: string, cb) => {
      if (event === "tag:read") {
        tagReadCallback = cb as (data: unknown) => void;
      }
      return () => {};
    });

    render(<SpoolForm />);
    await waitFor(() => expect(mockEventsOn).toHaveBeenCalled());
    expect(tagReadCallback).toBeDefined();

    tagReadCallback?.({
      uid: "DEADBEEF",
      date: "2026-04-12",
      supplierCode: "ESUN",
      materialCode: "E1001",
      color: "FF4010",
      lengthCode: "0330",
      serial: "000042",
      isBlank: false,
    });

    await waitFor(() => {
      expect(screen.getByText("DEADBEEF")).toBeInTheDocument();
      expect(screen.getByText("#FF4010")).toBeInTheDocument();
    });
    // Serial vai para o input controlado.
    const serialInput = screen.getByPlaceholderText("000001") as HTMLInputElement;
    expect(serialInput.value).toBe("000042");

    // Toast de sucesso para tag normal.
    expect(toastSuccesses.some((m) => m.includes("DEADBEEF") && m.includes("Tag lida"))).toBe(true);
  });

  it("aplica defaults e dispara toast.info quando isBlank=true", async () => {
    let tagReadCallback: ((data: unknown) => void) | undefined;
    mockEventsOn.mockImplementation((event: string, cb) => {
      if (event === "tag:read") {
        tagReadCallback = cb as (data: unknown) => void;
      }
      return () => {};
    });

    render(<SpoolForm />);
    await waitFor(() => expect(mockEventsOn).toHaveBeenCalled());

    tagReadCallback?.({
      uid: "BEEF1234",
      date: "",
      supplierCode: "",
      materialCode: "",
      color: "",
      lengthCode: "",
      serial: "",
      isBlank: true,
    });

    await waitFor(() => {
      expect(screen.getByText("BEEF1234")).toBeInTheDocument();
    });
    // Defaults: supplier "0276", color "000000", length "0330", serial "000001".
    expect(screen.getByText("#000000")).toBeInTheDocument();
    const serialInput = screen.getByPlaceholderText("000001") as HTMLInputElement;
    expect(serialInput.value).toBe("000001");

    // toast.info é roteado para toastSuccesses no nosso mock; cobre tag virgem.
    expect(toastSuccesses.some((m) => m.includes("BEEF1234") && m.includes("Tag virgem"))).toBe(true);
  });

  it("muda barra de status quando callback tag:status emite 'read'", async () => {
    let statusCallback: ((status: string) => void) | undefined;
    mockEventsOn.mockImplementation((event: string, cb) => {
      if (event === "tag:status") {
        statusCallback = cb as (status: string) => void;
      }
      return () => {};
    });

    render(<SpoolForm />);
    await waitFor(() => expect(mockEventsOn).toHaveBeenCalled());

    expect(screen.getByText(/Aguardando tag/)).toBeInTheDocument();

    statusCallback?.("read");
    await waitFor(() => {
      expect(screen.getByText(/Tag detectada automaticamente/)).toBeInTheDocument();
    });
  });

  it("exibe barra de erro quando callback tag:status emite 'error'", async () => {
    let statusCallback: ((status: string) => void) | undefined;
    mockEventsOn.mockImplementation((event: string, cb) => {
      if (event === "tag:status") {
        statusCallback = cb as (status: string) => void;
      }
      return () => {};
    });

    render(<SpoolForm />);
    await waitFor(() => expect(mockEventsOn).toHaveBeenCalled());

    statusCallback?.("error");
    await waitFor(() => {
      expect(screen.getByText(/Erro ao ler tag/)).toBeInTheDocument();
    });
  });
});
