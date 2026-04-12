import { HexColorPicker } from "react-colorful";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { presetColors } from "@/data/presets";

interface ColorPickerProps {
  value: string;
  onChange: (hex: string) => void;
}

export function ColorPicker({ value, onChange }: ColorPickerProps) {
  const handleHexInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const raw = e.target.value.replace(/[^0-9a-fA-F]/g, "").slice(0, 6);
    onChange(raw.toUpperCase());
  };

  const handlePickerChange = (color: string) => {
    // react-colorful retorna #RRGGBB
    onChange(color.replace("#", "").toUpperCase());
  };

  const displayColor = value.length === 6 ? `#${value}` : "#000000";

  return (
    <div className="space-y-3">
      <Label>Cor</Label>
      <div className="flex items-center gap-3">
        <div
          className="w-10 h-10 rounded border border-border shrink-0"
          style={{ backgroundColor: displayColor }}
        />
        <Input
          value={value}
          onChange={handleHexInput}
          placeholder="FF4010"
          maxLength={6}
          className="font-mono uppercase w-28"
        />
      </div>
      <HexColorPicker color={displayColor} onChange={handlePickerChange} />
      <div className="grid grid-cols-12 gap-1">
        {presetColors.map((color) => (
          <button
            key={color}
            type="button"
            className={`w-7 h-7 rounded cursor-pointer border transition-all ${
              value === color
                ? "ring-2 ring-primary ring-offset-1"
                : "border-border hover:scale-110"
            } ${color === "FFFFFF" ? "border-gray-300" : ""}`}
            style={{ backgroundColor: `#${color}` }}
            onClick={() => onChange(color)}
            title={color}
          />
        ))}
      </div>
    </div>
  );
}
