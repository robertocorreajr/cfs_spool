import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Loader2, Sparkles } from "lucide-react";
import appIcon from "@/assets/appicon.png";
import { triggerManualCheck } from "@/components/UpdateNotifier";
import { useEffect, useState } from "react";
import { EventsOn } from "../../wailsjs/runtime/runtime";

interface HeaderProps {
  version: string;
  uid: string;
  /**
   * onCheckForUpdate é injetável para testes — em produção fica indefinido
   * e o botão chama triggerManualCheck diretamente.
   */
  onCheckForUpdate?: () => Promise<void> | void;
  /**
   * hasUpdate força o estado de "atualização disponível" para testes.
   * Em produção é alimentado pelo evento "update:available" do backend.
   */
  hasUpdate?: boolean;
}

export function Header({ version, uid, onCheckForUpdate, hasUpdate }: HeaderProps) {
  const [checking, setChecking] = useState(false);
  const [internalHasUpdate, setInternalHasUpdate] = useState(false);

  // Escuta o evento emitido pelo backend e o mesmo "presente"
  // disparado pelo UpdateNotifier — sinaliza no botão sem precisar
  // de prop drilling. Limpa quando o usuário ignora a versão.
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

  const handleCheck = async () => {
    setChecking(true);
    try {
      const fn = onCheckForUpdate ?? triggerManualCheck;
      await fn();
    } finally {
      setChecking(false);
    }
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
        <Button
          variant={updateAvailable ? "default" : "ghost"}
          size="sm"
          className={
            "relative h-7 gap-1.5 text-xs " +
            (updateAvailable
              ? "bg-emerald-600 hover:bg-emerald-700 text-white animate-pulse"
              : "text-muted-foreground hover:text-foreground")
          }
          aria-label={
            updateAvailable
              ? "Atualização disponível — clique para ver"
              : "Procurar atualizações"
          }
          title={
            updateAvailable
              ? "Atualização disponível"
              : "Procurar atualizações"
          }
          disabled={checking}
          onClick={handleCheck}
        >
          {checking ? (
            <Loader2 className="h-3.5 w-3.5 animate-spin" />
          ) : (
            <Sparkles className="h-3.5 w-3.5" />
          )}
          {updateAvailable ? "Atualização" : "Atualizações"}
          {updateAvailable && (
            <span className="absolute -top-0.5 -right-0.5 flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-red-400 opacity-75" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-red-500" />
            </span>
          )}
        </Button>
      </div>
    </div>
  );
}
