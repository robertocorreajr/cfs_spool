import { useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { GetVersion } from "../../wailsjs/go/main/App";

export function Header() {
  const [version, setVersion] = useState("");

  useEffect(() => {
    GetVersion().then(setVersion);
  }, []);

  return (
    <div className="flex items-center justify-between px-6 py-4 border-b">
      <h1 className="text-2xl font-bold">CFS Spool</h1>
      {version && <Badge variant="secondary">{version}</Badge>}
    </div>
  );
}
