import { Toaster } from "sonner";
import { Header } from "@/components/Header";
import { SpoolForm } from "@/components/SpoolForm";
import { useState, useEffect } from "react";
import { GetVersion } from "../wailsjs/go/main/App";

function App() {
  const [version, setVersion] = useState("");
  const [uid, setUid] = useState("");

  useEffect(() => {
    GetVersion().then(setVersion);
  }, []);

  return (
    <div className="min-h-screen bg-background">
      <Toaster position="top-right" richColors />
      <Header version={version} uid={uid} />
      <main className="p-6">
        <SpoolForm onUidChange={setUid} />
      </main>
    </div>
  );
}

export default App;
