import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import type { MaterialOption, VendorOption } from "@/types/spool";

interface MaterialSelectProps {
  supplier: string;
  material: string;
  onSupplierChange: (value: string) => void;
  onMaterialChange: (value: string) => void;
  materials: MaterialOption[];
  vendors: VendorOption[];
}

export function MaterialSelect({
  supplier, material, onSupplierChange, onMaterialChange, materials, vendors,
}: MaterialSelectProps) {
  const filteredMaterials = materials.filter((m) => {
    if (supplier === "0276") return m.code >= "04000";
    if (supplier === "0000") return m.code < "04000";
    return true;
  });

  const handleMaterialChange = (code: string) => {
    onMaterialChange(code);
    // Auto-selecionar fornecedor
    if (code >= "04000") {
      onSupplierChange("0276");
    } else {
      onSupplierChange("0000");
    }
  };

  return (
    <div className="grid grid-cols-2 gap-4">
      <div className="space-y-2">
        <Label>Fornecedor</Label>
        <Select value={supplier} onValueChange={onSupplierChange}>
          <SelectTrigger><SelectValue placeholder="Selecione..." /></SelectTrigger>
          <SelectContent>
            {vendors.map((v) => (
              <SelectItem key={v.code} value={v.code}>{v.name}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2">
        <Label>Material</Label>
        <Select value={material} onValueChange={handleMaterialChange}>
          <SelectTrigger><SelectValue placeholder="Selecione..." /></SelectTrigger>
          <SelectContent>
            {filteredMaterials.map((m) => (
              <SelectItem key={m.code} value={m.code}>{m.name}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}
