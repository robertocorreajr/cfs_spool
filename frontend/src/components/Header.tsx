import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Sparkles } from "lucide-react";
import appIcon from "@/assets/appicon.png";
import { useEffect, useState } from "react";
import { EventsEmit, EventsOn } from "../../wailsjs/runtime/runtime";

interface HeaderProps {
  version: string;
  uid: string;
  /**
   * hasUpdate força o estado de "atualização disponível" em testes.
   * Em produção é alimentado pelos eventos "update:available" /
   * "update:cleared" emitidos pelo UpdateNotifier e pelo backend.
   */
  hasUpdate?: boolean;
}

// Header exibe título, UID, versão e — quando há nova release — um
// ícone que pisca no canto superior direito. Clique abre o modal de
// changelog (via evento "update:show" consumido pelo UpdateNotifier).
// Sem update disponível, o ícone fica oculto.
export function Header({ version, uid, hasUpdate }: HeaderProps) {
  const [internalHasUpdate, setInternalHasUpdate] = useState(false);

  useEffect(() => {
    const offAvailable = EventsOn("update:available", () => {
      setInternalHasUpdate(true);
    });
    const offCleared = EventsOn("update:cleared", () => {
      setInternalHasUpdate(false);
    });
    return () => {
      offAvailable();
      offCleared();
    };
  }, []);

  const updateAvailable = hasUpdate ?? internalHasUpdate;

  const handleClick = () => {
    EventsEmit("update:show");
  };

  return (
    <div className="flex items-center justify-between px-5 py-3 border-b bg-background">
      <div className="flex items-center gap-2.5">
        <img src={appIcon} alt="CFS Spool" className="w-9 h-9 rounded-lg" />
        <span className="text-lg font-bold tracking-tight">CFS Spool</span>
      </div>
      <div className="flex items-center gap-2">
        {uid && (
          <Badge variant="secondary" className="font-normal">
            UID: <span className="font-mono font-medium ml-1">{uid}</span>
          </Badge>
        )}
        {version && (
          <Badge variant="outline" className="font-normal text-muted-foreground">
            {version}
          </Badge>
        )}
        {updateAvailable && (
          <Button
            variant="ghost"
            size="icon"
            className="relative h-8 w-8 text-emerald-600 hover:text-emerald-700 hover:bg-emerald-50 animate-pulse"
            aria-label="Atualização disponível — clique para ver"
            title="Atualização disponível"
            onClick={handleClick}
          >
            <Sparkles className="h-4 w-4" />
            <span className="absolute -top-0.5 -right-0.5 flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-red-400 opacity-75" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-red-500" />
            </span>
          </Button>
        )}
      </div>
    </div>
  );
}
