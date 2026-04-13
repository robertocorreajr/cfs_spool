import { SpoolForm } from "@/components/SpoolForm";
import { Toaster } from "sonner";

function App() {
  return (
    <>
      <SpoolForm />
      <Toaster position="top-right" richColors />
    </>
  );
}

export default App;
