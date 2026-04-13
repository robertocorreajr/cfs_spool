import { useState } from "react";
import { HexColorPicker } from "react-colorful";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { presetColors } from "@/data/presets";

interface ColorPickerProps {
  value: string;
  onChange: (color: string) => void;
}

export function ColorPicker({ value, onChange }: ColorPickerProps) {
  const [open, setOpen] = useState(false);
  const [hexInput, setHexInput] = useState(value);

  const handleHexChange = (hex: string) => {
    const clean = hex.replace(/[^0-9A-Fa-f]/g, "").slice(0, 6);
    setHexInput(clean);
    if (clean.length === 6) {
      onChange(clean.toUpperCase());
    }
  };

  const handlePickerChange = (hex: string) => {
    // react-colorful retorna #RRGGBB
    const clean = hex.replace("#", "").toUpperCase();
    onChange(clean);
    setHexInput(clean);
  };

  const handlePresetClick = (color: string) => {
    onChange(color);
    setHexInput(color);
  };

  // Sincronizar input quando valor muda externamente
  if (value !== hexInput && value.length === 6 && hexInput.length === 6 && value.toUpperCase() !== hexInput.toUpperCase()) {
    setHexInput(value);
  }

  return (
    <div className="space-y-1.5">
      <Label className="text-xs font-medium text-muted-foreground">Cor</Label>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <button
            type="button"
            className="flex items-center gap-2 w-full border rounded-md px-3 py-2 hover:bg-accent/50 transition-colors cursor-pointer"
          >
            <div
              className="w-7 h-7 rounded-md border shrink-0"
              style={{ backgroundColor: `#${value}` }}
            />
            <span className="font-mono text-sm font-medium">#{value}</span>
          </button>
        </PopoverTrigger>
        <PopoverContent className="w-64 p-3" align="start">
          <HexColorPicker
            color={`#${value}`}
            onChange={handlePickerChange}
            style={{ width: "100%", height: "140px" }}
          />
          <div className="flex flex-wrap gap-1 mt-3">
            {presetColors.map((color) => (
              <button
                key={color}
                type="button"
                className={`w-5 h-5 rounded-sm border transition-all ${
                  value === color ? "ring-2 ring-primary ring-offset-1" : "hover:scale-110"
                }`}
                style={{ backgroundColor: `#${color}` }}
                onClick={() => handlePresetClick(color)}
                title={`#${color}`}
              />
            ))}
          </div>
          <div className="flex items-center gap-1.5 mt-3">
            <span className="text-sm text-muted-foreground">#</span>
            <Input
              value={hexInput}
              onChange={(e) => handleHexChange(e.target.value)}
              className="font-mono text-sm h-8"
              maxLength={6}
              placeholder="000000"
            />
          </div>
        </PopoverContent>
      </Popover>
    </div>
  );
}
