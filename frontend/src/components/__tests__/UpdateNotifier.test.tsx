import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";

const mockCheckForUpdate = vi.fn();
const mockIgnoreUpdateVersion = vi.fn();
const mockOpenURL = vi.fn();
const mockEventsEmit = vi.fn();

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
  EventsEmit: (...args: unknown[]) => mockEventsEmit(...args),
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
  mockEventsEmit.mockReset();
});

describe("UpdateNotifier", () => {
  it("não renderiza nada quando não há update detectado", () => {
    mockCheckForUpdate.mockResolvedValue({
      ...releaseFixture,
      available: false,
    });
    const { container } = render(<UpdateNotifier autoCheckEvent={false} />);
    expect(container.firstChild).toBeNull();
    expect(toastMessages).toHaveLength(0);
  });

  it("não dispara toast automaticamente quando recebe update:available", async () => {
    render(<UpdateNotifier />);

    const cb = subscribers["update:available"]?.[0];
    expect(cb).toBeDefined();
    cb?.(releaseFixture);

    // Comportamento novo: nenhum toast nem modal abrem sozinhos —
    // o usuário precisa clicar no ícone do Header (que dispara update:show).
    expect(toastMessages).toHaveLength(0);
    expect(screen.queryByText(/Auto-update/)).not.toBeInTheDocument();
  });

  it("abre o modal quando recebe update:show após update:available", async () => {
    render(<UpdateNotifier />);

    subscribers["update:available"]?.[0]?.(releaseFixture);
    subscribers["update:show"]?.[0]?.();

    await waitFor(() => {
      expect(screen.getByText(/Auto-update/)).toBeInTheDocument();
    });
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

  it("triggerManualCheck emite update:available e update:show quando há update", async () => {
    mockCheckForUpdate.mockResolvedValue(releaseFixture);

    await triggerManualCheck();

    const events = mockEventsEmit.mock.calls.map((c) => c[0]);
    expect(events).toContain("update:available");
    expect(events).toContain("update:show");
  });

  it("triggerManualCheck mostra error quando o backend falha", async () => {
    mockCheckForUpdate.mockRejectedValue(new Error("rede caiu"));

    await triggerManualCheck();

    const err = toastMessages.find((t) => t.kind === "error");
    expect(err).toBeDefined();
    expect(err?.msg).toContain("rede caiu");
  });

  it("clicar em 'Abrir release no GitHub' chama OpenURL com a URL da release", async () => {
    render(<UpdateNotifier />);
    subscribers["update:available"]?.[0]?.(releaseFixture);
    subscribers["update:show"]?.[0]?.();

    await waitFor(() => {
      expect(screen.getByText(/Auto-update/)).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole("button", { name: /Abrir release/i }));
    expect(mockOpenURL).toHaveBeenCalledWith(releaseFixture.url);
  });

  it("clicar em 'Ignorar esta versão' chama IgnoreUpdateVersion + emite update:cleared", async () => {
    mockIgnoreUpdateVersion.mockResolvedValue(undefined);
    render(<UpdateNotifier />);
    subscribers["update:available"]?.[0]?.(releaseFixture);
    subscribers["update:show"]?.[0]?.();

    await waitFor(() => {
      expect(screen.getByText(/Auto-update/)).toBeInTheDocument();
    });

    fireEvent.click(
      screen.getByRole("button", { name: /Ignorar esta versão/i }),
    );

    await waitFor(() => {
      expect(mockIgnoreUpdateVersion).toHaveBeenCalledWith("v3.1.0");
    });
    expect(mockEventsEmit).toHaveBeenCalledWith("update:cleared");
  });

  it("pull no mount emite update:available quando backend reporta versão nova", async () => {
    mockCheckForUpdate.mockResolvedValue(releaseFixture);

    render(<UpdateNotifier />);

    await waitFor(() => {
      expect(mockEventsEmit).toHaveBeenCalledWith(
        "update:available",
        expect.objectContaining({ version: "v3.1.0" }),
      );
    });
  });
});
