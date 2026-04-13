import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { LengthOption } from "@/types/spool";

interface LengthSelectProps {
  length: string;
  customGrams: string;
  onLengthChange: (value: string) => void;
  onCustomGramsChange: (value: string) => void;
  lengths: LengthOption[];
}

export function LengthSelect({
  length, customGrams, onLengthChange, onCustomGramsChange, lengths,
}: LengthSelectProps) {
  return (
    <div className="space-y-1.5">
      <Label className="text-xs font-medium text-muted-foreground">Comprimento</Label>
      <Select value={length} onValueChange={onLengthChange}>
        <SelectTrigger><SelectValue placeholder="Selecione..." /></SelectTrigger>
        <SelectContent>
          {lengths.map((l) => (
            <SelectItem key={l.code} value={l.code}>{l.name}</SelectItem>
          ))}
        </SelectContent>
      </Select>
      {length === "CUSTOM" && (
        <Input
          type="number"
          value={customGrams}
          onChange={(e) => onCustomGramsChange(e.target.value)}
          placeholder="Gramas (ex: 750)"
        />
      )}
    </div>
  );
}
