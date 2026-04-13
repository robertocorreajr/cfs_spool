import { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ColorPicker } from "@/components/ColorPicker";
import { MaterialSelect } from "@/components/MaterialSelect";
import { LengthSelect } from "@/components/LengthSelect";
import { toast } from "sonner";
import { WriteTag, GetOptions, GetVersion } from "../../wailsjs/go/main/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { Header } from "@/components/Header";
import { Save } from "lucide-react";
import type { OptionsResponse } from "@/types/spool";

type TagStatus = "waiting" | "read" | "error";

export function SpoolForm() {
  const [options, setOptions] = useState<OptionsResponse>({ materials: [], vendors: [], lengths: [] });
  const [version, setVersion] = useState("");

  // Campos do formulario
  const [date, setDate] = useState(new Date().toISOString().split("T")[0]);
  const [supplier, setSupplier] = useState("0276");
  const [material, setMaterial] = useState("");
  const [color, setColor] = useState("000000");
  const [length, setLength] = useState("0330");
  const [customGrams, setCustomGrams] = useState("");
  const [serial, setSerial] = useState("000001");

  // Estado
  const [uid, setUid] = useState("");
  const [tagStatus, setTagStatus] = useState<TagStatus>("waiting");
  const [isWriting, setIsWriting] = useState(false);
  const [writeCount, setWriteCount] = useState(0);

  useEffect(() => {
    GetOptions().then(setOptions).catch(() => toast.error("Erro ao carregar opcoes"));
    GetVersion().then(setVersion);
  }, []);

  // Escuta eventos do watcher de tags
  useEffect(() => {
    const offStatus = EventsOn("tag:status", (status: TagStatus) => {
      setTagStatus(status);
    });
    const offRead = EventsOn("tag:read", (data: any) => {
      setTagStatus("read");
      applyTagData(data);
    });
    return () => { offStatus(); offRead(); };
  }, []);

  const applyTagData = (data: any) => {
    setUid(data.uid);
    setDate(data.date || new Date().toISOString().split("T")[0]);
    setSupplier(data.supplierCode || "0276");
    setMaterial(data.materialCode || "");
    setColor(data.color || "000000");
    setLength(data.lengthCode || "0330");
    setSerial(data.serial || "000001");
    setWriteCount(0);
    if (data.isBlank) {
      toast.info(`Tag virgem — UID: ${data.uid}`);
    } else {
      toast.success(`Tag lida — UID: ${data.uid}`);
    }
  };

  const handleSerialChange = (value: string) => {
    const clean = value.replace(/\D/g, "").slice(0, 6);
    setSerial(clean);
    setWriteCount(0);
  };

  const incrementSerial = () => {
    const num = parseInt(serial || "0", 10);
    const next = Math.min(num + 1, 999999);
    setSerial(String(next).padStart(6, "0"));
  };

  const handleWrite = async () => {
    if (!material) { toast.error("Selecione um material"); return; }
    if (color.length !== 6) { toast.error("Cor deve ter 6 caracteres hex"); return; }
    setIsWriting(true);
    try {
      const lengthValue = length === "CUSTOM" ? customGrams : length;
      await WriteTag({ date, supplier, material, color, length: lengthValue, serial: serial || "000001" });
      const newCount = writeCount + 1;
      setWriteCount(newCount);
      if (newCount >= 2) {
        incrementSerial();
        setWriteCount(0);
      }
      toast.success(`Tag gravada! (${newCount}/2)`);
    } catch (err: any) {
      toast.error(err?.message || String(err));
    } finally {
      setIsWriting(false);
    }
  };

  // Barra de status da tag
  const statusBar = () => {
    if (tagStatus === "waiting") return (
      <div className="flex items-center gap-2 px-5 py-2 bg-amber-50 border-b border-amber-200">
        <div className="w-2 h-2 rounded-full bg-amber-400 animate-pulse" />
        <span className="text-xs font-medium text-amber-700">Aguardando tag no leitor...</span>
      </div>
    );
    if (tagStatus === "read") return (
      <div className="flex items-center gap-2 px-5 py-2 bg-green-50 border-b border-green-200">
        <div className="w-2 h-2 rounded-full bg-green-500" />
        <span className="text-xs font-medium text-green-700">Tag detectada automaticamente</span>
      </div>
    );
    return (
      <div className="flex items-center gap-2 px-5 py-2 bg-red-50 border-b border-red-200">
        <div className="w-2 h-2 rounded-full bg-red-500" />
        <span className="text-xs font-medium text-red-700">Erro ao ler tag</span>
      </div>
    );
  };

  // Layout de pagina completa
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Header version={version} uid={uid} />
      {statusBar()}
      <div className="flex-1 p-4 pb-24">
        <Card className="max-w-2xl mx-auto">
          <CardContent className="pt-5 space-y-4">
            <MaterialSelect
              supplier={supplier}
              material={material}
              onSupplierChange={setSupplier}
              onMaterialChange={setMaterial}
              materials={options.materials}
              vendors={options.vendors}
            />
            <div className="grid grid-cols-2 gap-3">
              <ColorPicker value={color} onChange={setColor} />
              <LengthSelect
                length={length}
                customGrams={customGrams}
                onLengthChange={setLength}
                onCustomGramsChange={setCustomGrams}
                lengths={options.lengths}
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1.5">
                <Label className="text-xs font-medium text-muted-foreground">Data</Label>
                <Input type="date" value={date} onChange={(e) => setDate(e.target.value)} />
              </div>
              <div className="space-y-1.5">
                <Label className="text-xs font-medium text-muted-foreground">
                  Serial {writeCount > 0 && <span className="ml-1.5 text-muted-foreground/60">({writeCount}/2)</span>}
                </Label>
                <Input
                  value={serial}
                  onChange={(e) => handleSerialChange(e.target.value)}
                  placeholder="000001"
                  maxLength={6}
                  className="font-mono"
                />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Rodape fixo com botao Gravar */}
      <div className="fixed bottom-0 left-0 right-0 p-4 bg-background/95 backdrop-blur border-t">
        <div className="max-w-2xl mx-auto">
          <Button onClick={handleWrite} disabled={isWriting} className="w-full" size="lg">
            <Save className="mr-2 h-4 w-4" />
            {isWriting ? "Gravando..." : "Gravar Tag"}
          </Button>
        </div>
      </div>

    </div>
  );
}
