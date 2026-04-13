import { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { ColorPicker } from "@/components/ColorPicker";
import { MaterialSelect } from "@/components/MaterialSelect";
import { LengthSelect } from "@/components/LengthSelect";
import { toast } from "sonner";
import { ReadTag, WriteTag, GetOptions } from "../../wailsjs/go/main/App";
import type { OptionsResponse } from "@/types/spool";

interface SpoolFormProps {
  onUidChange?: (uid: string) => void;
}

export function SpoolForm({ onUidChange }: SpoolFormProps) {
  // Opcoes dos dropdowns
  const [options, setOptions] = useState<OptionsResponse>({ materials: [], vendors: [], lengths: [] });

  // Campos do formulario
  const [date, setDate] = useState(new Date().toISOString().split("T")[0]);
  const [supplier, setSupplier] = useState("0276");
  const [material, setMaterial] = useState("");
  const [color, setColor] = useState("000000");
  const [length, setLength] = useState("0330");
  const [customGrams, setCustomGrams] = useState("");
  const [serial, setSerial] = useState("");

  // Estado
  const [uid, setUid] = useState("");
  const [isReading, setIsReading] = useState(false);
  const [isWriting, setIsWriting] = useState(false);

  useEffect(() => {
    GetOptions().then(setOptions).catch(() =>
      toast.error("Erro ao carregar opcoes")
    );
  }, []);

  const handleRead = async () => {
    setIsReading(true);
    try {
      const data = await ReadTag();
      setUid(data.uid);
      onUidChange?.(data.uid);
      setDate(data.date);
      setSupplier(data.supplierCode);
      setMaterial(data.materialCode);
      setColor(data.color);
      setLength(data.lengthCode);
      setSerial(data.serial);
      toast.success(`Tag lida com sucesso! UID: ${data.uid}`);
    } catch (err: any) {
      toast.error(err?.message || String(err));
    } finally {
      setIsReading(false);
    }
  };

  const handleWrite = async () => {
    if (!material) {
      toast.error("Selecione um material");
      return;
    }
    if (color.length !== 6) {
      toast.error("Cor deve ter 6 caracteres hex");
      return;
    }

    setIsWriting(true);
    try {
      const lengthValue = length === "CUSTOM" ? customGrams : length;
      await WriteTag({
        date,
        supplier,
        material,
        color,
        length: lengthValue,
        serial: serial || "000001",
      });
      toast.success("Tag gravada com sucesso!");
    } catch (err: any) {
      toast.error(err?.message || String(err));
    } finally {
      setIsWriting(false);
    }
  };

  return (
    <Card className="max-w-2xl mx-auto">
      <CardContent className="pt-6 space-y-6">
        {/* Secao Ler Tag */}
        <div className="flex items-center gap-4">
          <Button onClick={handleRead} disabled={isReading} variant="default" className="shrink-0">
            {isReading ? "Lendo..." : "Ler Tag"}
          </Button>
          {uid && (
            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">UID:</span>
              <Badge variant="outline" className="font-mono">{uid}</Badge>
            </div>
          )}
        </div>

        <Separator />

        {/* Data */}
        <div className="space-y-2">
          <Label>Data</Label>
          <Input type="date" value={date} onChange={(e) => setDate(e.target.value)} />
        </div>

        {/* Fornecedor + Material */}
        <MaterialSelect
          supplier={supplier}
          material={material}
          onSupplierChange={setSupplier}
          onMaterialChange={setMaterial}
          materials={options.materials}
          vendors={options.vendors}
        />

        {/* Cor */}
        <ColorPicker value={color} onChange={setColor} />

        {/* Comprimento */}
        <LengthSelect
          length={length}
          customGrams={customGrams}
          onLengthChange={setLength}
          onCustomGramsChange={setCustomGrams}
          lengths={options.lengths}
        />

        {/* Serial */}
        <div className="space-y-2">
          <Label>Serial</Label>
          <Input
            value={serial}
            onChange={(e) => setSerial(e.target.value.replace(/\D/g, "").slice(0, 6))}
            placeholder="000001"
            maxLength={6}
            className="font-mono"
          />
        </div>

        <Separator />

        {/* Gravar */}
        <Button
          onClick={handleWrite}
          disabled={isWriting}
          className="w-full"
          size="lg"
        >
          {isWriting ? "Gravando..." : "Gravar Tag"}
        </Button>
      </CardContent>
    </Card>
  );
}
