import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";

// Mocks dos bindings — devem vir antes do import do componente.
const mockCheckForUpdate = vi.fn();
const mockIgnoreUpdateVersion = vi.fn();
const mockOpenURL = vi.fn();

// Subscribers do EventsOn — guardamos por nome de evento para podermos
// disparar manualmente nos testes (simulando o startup do backend).
type Subscriber = (...args: unknown[]) => void;
const subscribers: Record<string, Subscriber[]> = {};

vi.mock("../../../wailsjs/go/main/App", () => ({
  CheckForUpdate: () => mockCheckForUpdate(),
  IgnoreUpdateVersion: (v: string) => mockIgnoreUpdateVersion(v),
  OpenURL: (url: string) => mockOpenURL(url),
}));

vi.mock("../../../wailsjs/runtime/runtime", () => ({
  EventsOn: (event: string, cb: Subscriber) => {
    subscribers[event] = subscribers[event] || [];
    subscribers[event].push(cb);
    return () => {
      const list = subscribers[event];
      if (list) {
        const idx = list.indexOf(cb);
        if (idx >= 0) list.splice(idx, 1);
      }
    };
  },
}));

const toastMessages: { kind: string; msg: string; opts?: any }[] = [];
vi.mock("sonner", () => ({
  toast: Object.assign(
    (msg: string, opts?: any) => {
      toastMessages.push({ kind: "default", msg, opts });
    },
    {
      success: (msg: string) => toastMessages.push({ kind: "success", msg }),
      error: (msg: string) => toastMessages.push({ kind: "error", msg }),
      info: (msg: string) => toastMessages.push({ kind: "info", msg }),
    },
  ),
}));

import { UpdateNotifier, triggerManualCheck } from "@/components/UpdateNotifier";

const releaseFixture = {
  available: true,
  ignored: false,
  current: "v3.0.0",
  version: "v3.1.0",
  name: "v3.1.0",
  url: "https://example.com/r/v3.1.0",
  publishedAt: "2026-05-01T12:00:00Z",
  body: "## Novidades\n- Auto-update",
};

beforeEach(() => {
  toastMessages.length = 0;
  Object.keys(subscribers).forEach((k) => delete subscribers[k]);
  mockCheckForUpdate.mockReset();
  mockIgnoreUpdateVersion.mockReset();
  mockOpenURL.mockReset();
});

describe("UpdateNotifier", () => {
  it("não renderiza nada antes de receber um evento", () => {
    const { container } = render(<UpdateNotifier />);
    // Nenhum dialog aberto, nenhum toast disparado ainda.
    expect(container.firstChild).toBeNull();
    expect(toastMessages).toHaveLength(0);
  });

  it("dispara toast quando recebe update:available do backend", () => {
    render(<UpdateNotifier />);

    // Simula o evento que o startup do backend emite.
    const cb = subscribers["update:available"]?.[0];
    expect(cb).toBeDefined();
    cb?.(releaseFixture);

    expect(toastMessages).toHaveLength(1);
    expect(toastMessages[0].msg).toContain("v3.1.0");
    expect(toastMessages[0].opts.description).toContain("v3.0.0");
  });

  it("triggerManualCheck mostra success quando não há update", async () => {
    mockCheckForUpdate.mockResolvedValue({
      ...releaseFixture,
      available: false,
    });

    await triggerManualCheck();

    expect(toastMessages.some((t) => t.kind === "success")).toBe(true);
    expect(
      toastMessages.find((t) => t.kind === "success")?.msg,
    ).toContain("mais recente");
  });

  it("triggerManualCheck reusa o presenter do UpdateNotifier quando há update", async () => {
    render(<UpdateNotifier />);

    mockCheckForUpdate.mockResolvedValue(releaseFixture);
    await triggerManualCheck();

    // Deve ter disparado o toast vindo do presenter (default kind).
    expect(toastMessages.some((t) => t.kind === "default")).toBe(true);
  });

  it("triggerManualCheck mostra error quando o backend falha", async () => {
    mockCheckForUpdate.mockRejectedValue(new Error("rede caiu"));

    await triggerManualCheck();

    const err = toastMessages.find((t) => t.kind === "error");
    expect(err).toBeDefined();
    expect(err?.msg).toContain("rede caiu");
  });

  it("clicar em 'Ver detalhes' abre o dialog com o changelog", async () => {
    render(<UpdateNotifier />);

    const cb = subscribers["update:available"]?.[0];
    cb?.(releaseFixture);

    // O toast tem uma `action` com onClick — exercitamos manualmente.
    const action = toastMessages[0].opts.action;
    expect(action.label).toBe("Ver detalhes");
    action.onClick();

    await waitFor(() => {
      expect(screen.getByText(/Auto-update/)).toBeInTheDocument();
    });
    expect(
      screen.getByRole("button", { name: /Abrir release/i }),
    ).toBeInTheDocument();
  });

  it("clicar em 'Abrir release no GitHub' chama OpenURL com a URL da release", async () => {
    render(<UpdateNotifier />);

    const cb = subscribers["update:available"]?.[0];
    cb?.(releaseFixture);
    toastMessages[0].opts.action.onClick();

    await waitFor(() => {
      expect(screen.getByText(/Auto-update/)).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole("button", { name: /Abrir release/i }));
    expect(mockOpenURL).toHaveBeenCalledWith(releaseFixture.url);
  });

  it("clicar em 'Ignorar' chama IgnoreUpdateVersion com a tag certa", async () => {
    mockIgnoreUpdateVersion.mockResolvedValue(undefined);
    render(<UpdateNotifier />);

    const cb = subscribers["update:available"]?.[0];
    cb?.(releaseFixture);

    const cancelAction = toastMessages[0].opts.cancel;
    expect(cancelAction.label).toBe("Ignorar");
    cancelAction.onClick();

    await waitFor(() => {
      expect(mockIgnoreUpdateVersion).toHaveBeenCalledWith("v3.1.0");
    });
  });
});
