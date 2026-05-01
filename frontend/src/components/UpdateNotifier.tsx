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
import { EventsOn } from "../../wailsjs/runtime/runtime";

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

// UpdateNotifier renderiza o modal de changelog e responde tanto ao evento
// emitido pelo startup do backend quanto a triggerManualCheck (botão
// "Verificar atualizações" no Header). O componente expõe seu callback de
// "show" via window para que o fluxo manual reuse o mesmo dialog/toast.
export function UpdateNotifier({ autoCheckEvent = true }: UpdateNotifierProps = {}) {
  const [info, setInfo] = useState<UpdateInfo | null>(null);
  const [showDialog, setShowDialog] = useState(false);

  const presentUpdate = (data: UpdateInfo) => {
    setInfo(data);
    toast(`Nova versão ${data.version} disponível`, {
      description: `Você está na ${data.current || "dev"}.`,
      duration: 15_000,
      action: {
        label: "Ver detalhes",
        onClick: () => setShowDialog(true),
      },
      cancel: {
        label: "Ignorar",
        onClick: () => handleIgnore(data.version),
      },
    });
  };

  // Registra a função de "apresentar update" no objeto global para que
  // triggerManualCheck (chamado de outros componentes, como o Header) possa
  // reutilizá-la. Single-instance — o SpoolForm monta um único notifier.
  useEffect(() => {
    pendingPresenter = presentUpdate;
    return () => {
      pendingPresenter = null;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Subscreve ao evento "update:available" emitido pelo backend no startup.
  useEffect(() => {
    if (!autoCheckEvent) return;
    const off = EventsOn("update:available", (data: UpdateInfo) => {
      presentUpdate(data);
    });
    return () => {
      off();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [autoCheckEvent]);

  const handleIgnore = async (version: string) => {
    try {
      await IgnoreUpdateVersion(version);
      toast.success(`Versão ${version} silenciada nesta instalação`);
      setShowDialog(false);
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

// pendingPresenter é o handler do UpdateNotifier mais recente montado.
// Mantido em escopo de módulo (não em window) para preservar tipagem e
// permitir testes via vi.mock — o componente registra/desregistra no mount.
let pendingPresenter: ((info: UpdateInfo) => void) | null = null;

/**
 * triggerManualCheck consulta o backend e exibe feedback ao usuário.
 *
 * - Se a versão atual já é a mais recente: toast.success "você está atualizado".
 * - Se há nova versão: delega ao UpdateNotifier (mesmo fluxo do startup).
 * - Em erro: toast.error com a mensagem.
 *
 * Pensado para ser chamado pelo botão "Verificar atualizações" no Header.
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
    if (pendingPresenter) {
      pendingPresenter(info);
    } else {
      // Sem notifier montado: ainda assim damos um feedback útil.
      toast(`Nova versão ${info.version} disponível`, {
        description: info.url,
        duration: 15_000,
      });
    }
  } catch (err) {
    toast.error(`Erro ao verificar atualizações: ${String(err)}`);
  }
}
