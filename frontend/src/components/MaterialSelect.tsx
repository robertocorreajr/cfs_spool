import { useState } from "react";
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
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
  const [supplierOpen, setSupplierOpen] = useState(false);
  const [materialOpen, setMaterialOpen] = useState(false);

  const filteredMaterials = materials.filter((m) => m.vendor === supplier);
  const selectedVendor = vendors.find((v) => v.code === supplier);
  const selectedMaterial = materials.find((m) => m.code === material);

  const handleMaterialChange = (code: string) => {
    onMaterialChange(code);
    setMaterialOpen(false);
    // Auto-selecionar fornecedor baseado no vendor do material
    const mat = materials.find((m) => m.code === code);
    if (mat) {
      onSupplierChange(mat.vendor);
    }
  };

  const handleSupplierChange = (code: string) => {
    onSupplierChange(code);
    setSupplierOpen(false);
    // Limpar material se não pertencer ao novo fornecedor
    if (material) {
      const mat = materials.find((m) => m.code === material);
      if (mat && mat.vendor !== code) {
        onMaterialChange("");
      }
    }
  };

  return (
    <div className="grid grid-cols-2 gap-3">
      <div className="space-y-1.5">
        <Label className="text-xs font-medium text-muted-foreground">Fornecedor</Label>
        <Popover open={supplierOpen} onOpenChange={setSupplierOpen}>
          <PopoverTrigger asChild>
            <Button variant="outline" role="combobox" aria-expanded={supplierOpen} className="w-full justify-between font-normal">
              {selectedVendor?.name ?? "Buscar fornecedor..."}
              <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-[--radix-popover-trigger-width] p-0">
            <Command>
              <CommandInput placeholder="Buscar fornecedor..." />
              <CommandList>
                <CommandEmpty>Nenhum fornecedor encontrado.</CommandEmpty>
                <CommandGroup>
                  {vendors.map((v) => (
                    <CommandItem key={v.code} value={v.name} onSelect={() => handleSupplierChange(v.code)}>
                      <Check className={cn("mr-2 h-4 w-4", supplier === v.code ? "opacity-100" : "opacity-0")} />
                      {v.name}
                    </CommandItem>
                  ))}
                </CommandGroup>
              </CommandList>
            </Command>
          </PopoverContent>
        </Popover>
      </div>
      <div className="space-y-1.5">
        <Label className="text-xs font-medium text-muted-foreground">Material</Label>
        <Popover open={materialOpen} onOpenChange={setMaterialOpen}>
          <PopoverTrigger asChild>
            <Button variant="outline" role="combobox" aria-expanded={materialOpen} className="w-full justify-between font-normal">
              {selectedMaterial?.name ?? "Buscar material..."}
              <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-[--radix-popover-trigger-width] p-0">
            <Command>
              <CommandInput placeholder="Buscar material..." />
              <CommandList>
                <CommandEmpty>Nenhum material encontrado.</CommandEmpty>
                <CommandGroup>
                  {filteredMaterials.map((m) => (
                    <CommandItem key={m.code} value={m.name} onSelect={() => handleMaterialChange(m.code)}>
                      <Check className={cn("mr-2 h-4 w-4", material === m.code ? "opacity-100" : "opacity-0")} />
                      {m.name}
                    </CommandItem>
                  ))}
                </CommandGroup>
              </CommandList>
            </Command>
          </PopoverContent>
        </Popover>
      </div>
    </div>
  );
}
