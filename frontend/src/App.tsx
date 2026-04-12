import { Toaster } from "sonner";
import { Header } from "@/components/Header";
import { SpoolForm } from "@/components/SpoolForm";

function App() {
  return (
    <div className="min-h-screen bg-background">
      <Toaster position="top-right" richColors />
      <Header />
      <main className="p-6">
        <SpoolForm />
      </main>
    </div>
  );
}

export default App;
