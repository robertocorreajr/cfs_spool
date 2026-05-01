import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Download, ExternalLink, X } from "lucide-react";
import { toast } from "sonner";
import {
  CheckForUpdate,
  IgnoreUpdateVersion,
  OpenURL,
} from "../../wailsjs/go/main/App";
import { EventsEmit, EventsOn } from "../../wailsjs/runtime/runtime";

// UpdateInfo espelha o struct exposto pelo backend (Go App.UpdateInfo).
// Mantemos a tipagem local em vez de importar do `wailsjs/go/models`
// porque models.ts é regenerado pelo wails — o type aqui fica explícito
// no contrato do componente.
export interface UpdateInfo {
  available: boolean;
  ignored: boolean;
  current: string;
  version: string;
  name: string;
  url: string;
  publishedAt: string;
  body: string;
}

interface UpdateNotifierProps {
  /**
   * autoCheckEvent indica se o componente deve subscrever ao evento
   * "update:available" emitido pelo backend no startup. Em testes podemos
   * deixar `false` para evitar dependência do runtime Wails real.
   * Default: true.
   */
  autoCheckEvent?: boolean;
}

// UpdateNotifier mantém o estado da release detectada e renderiza o
// modal de changelog. Não dispara toasts nem abre o modal automaticamente —
// o fluxo é: backend detecta → emite "update:available" + Header pisca →
// usuário clica no ícone → "update:show" abre o dialog.
export function UpdateNotifier({ autoCheckEvent = true }: UpdateNotifierProps = {}) {
  const [info, setInfo] = useState<UpdateInfo | null>(null);
  const [showDialog, setShowDialog] = useState(false);

  // Carrega a release mais nova (uma vez no mount). Reaproveita o evento
  // "update:available" caso ele chegue depois (ex.: o usuário ficou offline
  // no startup e a rede voltou — backend pode reemitir no futuro).
  useEffect(() => {
    if (!autoCheckEvent) return;

    const offAvail = EventsOn("update:available", (data: UpdateInfo) => {
      setInfo(data);
    });
    const offShow = EventsOn("update:show", () => {
      setShowDialog(true);
    });

    let cancelled = false;
    (async () => {
      try {
        const data = (await CheckForUpdate()) as UpdateInfo;
        if (cancelled) return;
        if (data.available && !data.ignored) {
          setInfo(data);
          // Avisa o Header para acender o ícone — útil quando o pull do
          // frontend é mais rápido que o evento do backend.
          EventsEmit("update:available", data);
        }
      } catch {
        // Silencioso: erro de rede / 404 / rate limit não trava o app.
      }
    })();

    return () => {
      cancelled = true;
      offAvail();
      offShow();
    };
  }, [autoCheckEvent]);

  const handleIgnore = async (version: string) => {
    try {
      await IgnoreUpdateVersion(version);
      toast.success(`Versão ${version} silenciada nesta instalação`);
      setShowDialog(false);
      setInfo(null);
      // Avisa o Header para apagar o ícone piscando.
      EventsEmit("update:cleared");
    } catch (err) {
      toast.error(`Erro ao silenciar versão: ${String(err)}`);
    }
  };

  const handleOpenRelease = () => {
    if (!info) return;
    OpenURL(info.url);
    setShowDialog(false);
  };

  if (!info) {
    return null;
  }

  return (
    <Dialog open={showDialog} onOpenChange={setShowDialog}>
      <DialogContent className="max-w-2xl max-h-[80vh] flex flex-col">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Download className="h-5 w-5 text-primary" />
            Nova versão {info.version} disponível
          </DialogTitle>
          <DialogDescription>
            Você está na versão {info.current || "dev"}. Confira as novidades
            abaixo antes de baixar.
          </DialogDescription>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto rounded-md border bg-muted/30 p-4 text-sm font-mono whitespace-pre-wrap">
          {info.body || "(sem changelog publicado)"}
        </div>

        <DialogFooter className="gap-2 sm:gap-2">
          <Button variant="ghost" onClick={() => handleIgnore(info.version)}>
            <X className="mr-2 h-4 w-4" />
            Ignorar esta versão
          </Button>
          <Button onClick={handleOpenRelease}>
            <ExternalLink className="mr-2 h-4 w-4" />
            Abrir release no GitHub
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

/**
 * triggerManualCheck é o fluxo do botão "Verificar atualizações" — mas como
 * o ícone do Header agora só aparece quando há novidade, esta função
 * abre o dialog se o estado já tem release carregada e roda um check
 * fresco caso contrário. Mantida exportada para retrocompatibilidade
 * com testes existentes.
 */
export async function triggerManualCheck(): Promise<void> {
  try {
    const info = (await CheckForUpdate()) as UpdateInfo;
    if (!info.available) {
      toast.success(
        `Você está na versão mais recente (${info.current || "dev"}).`,
      );
      return;
    }
    EventsEmit("update:available", info);
    EventsEmit("update:show");
  } catch (err) {
    toast.error(`Erro ao verificar atualizações: ${String(err)}`);
  }
}
