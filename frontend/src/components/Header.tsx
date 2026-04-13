import { Badge } from "@/components/ui/badge";
import appIcon from "@/assets/appicon.png";

interface HeaderProps {
  version: string;
  uid: string;
}

export function Header({ version, uid }: HeaderProps) {
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
      </div>
    </div>
  );
}
