import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Loader2, Search } from "lucide-react";
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
          size="sm"
          className="h-7 gap-1.5 text-xs text-muted-foreground hover:text-foreground"
          aria-label="Procurar atualizações"
          title="Procurar atualizações"
          disabled={checking}
          onClick={handleCheck}
        >
          {checking ? (
            <Loader2 className="h-3.5 w-3.5 animate-spin" />
          ) : (
            <Search className="h-3.5 w-3.5" />
          )}
          Atualizações
        </Button>
      </div>
    </div>
  );
}
