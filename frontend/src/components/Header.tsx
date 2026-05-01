import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { RefreshCw } from "lucide-react";
import appIcon from "@/assets/appicon.png";
import { triggerManualCheck } from "@/components/UpdateNotifier";
import { useState } from "react";

interface HeaderProps {
  version: string;
  uid: string;
  /**
   * onCheckForUpdate é injetável para testes — em produção fica indefinido
   * e o botão chama triggerManualCheck diretamente.
   */
  onCheckForUpdate?: () => Promise<void> | void;
}

export function Header({ version, uid, onCheckForUpdate }: HeaderProps) {
  const [checking, setChecking] = useState(false);

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
          variant="ghost"
          size="icon"
          className="h-7 w-7"
          aria-label="Verificar atualizações"
          title="Verificar atualizações"
          disabled={checking}
          onClick={handleCheck}
        >
          <RefreshCw className={"h-4 w-4 " + (checking ? "animate-spin" : "")} />
        </Button>
      </div>
    </div>
  );
}
