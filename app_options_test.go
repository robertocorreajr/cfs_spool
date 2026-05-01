package main

import "testing"

// TestMaterialCodesUnique verifica que não há códigos duplicados em `materials`.
// Códigos duplicados quebram a filtragem por vendor e o lookup por Code, gerando
// comportamento inconsistente nos dropdowns e na gravação da tag.
func TestMaterialCodesUnique(t *testing.T) {
	vistos := make(map[string]string, len(materials))
	for _, m := range materials {
		if anterior, dup := vistos[m.Code]; dup {
			t.Errorf("código duplicado em materials: %q aparece em %q e %q", m.Code, anterior, m.Name)
			continue
		}
		vistos[m.Code] = m.Name
	}
}

// TestMaterialsVendorsConsistency verifica que todo Vendor referenciado em
// `materials` existe em `vendors`. Sem isso, o dropdown de filtragem por vendor
// não exibe o material correspondente.
func TestMaterialsVendorsConsistency(t *testing.T) {
	codigosValidos := make(map[string]struct{}, len(vendors))
	for _, v := range vendors {
		codigosValidos[v.Code] = struct{}{}
	}

	for _, m := range materials {
		if _, ok := codigosValidos[m.Vendor]; !ok {
			t.Errorf("material %q (%s) referencia vendor %q inexistente em vendors", m.Code, m.Name, m.Vendor)
		}
	}
}

// TestLengthsHasCustom verifica que a entrada "CUSTOM" está presente em `lengths`.
// O frontend depende dessa entrada para exibir o campo de gramas customizado;
// sua ausência quebra a UX de comprimentos personalizados.
func TestLengthsHasCustom(t *testing.T) {
	for _, l := range lengths {
		if l.Code == "CUSTOM" {
			return
		}
	}
	t.Error("entrada com Code=\"CUSTOM\" não encontrada em lengths")
}
