package main

import (
	"strings"
	"testing"
)

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

// TestMaterialPrefixInvariant valida a regra completa de prefixos para 100% das entradas:
//   - Códigos numéricos de 5 dígitos (00001–29999) → Vendor deve ser "0000" ou "0276"
//   - Prefixo "E" seguido de 4 dígitos → Vendor deve ser "ESUN"
//   - Prefixo "P" seguido de 4 dígitos → Vendor deve ser "POLY"
func TestMaterialPrefixInvariant(t *testing.T) {
	for _, m := range materials {
		code := m.Code
		switch {
		case len(code) == 5 && isAllDigits(code):
			// Genérico ou Creality — vendor deve ser "0000" ou "0276"
			if m.Vendor != "0000" && m.Vendor != "0276" {
				t.Errorf("code=%q (numérico 5 dígitos): vendor esperado \"0000\" ou \"0276\", recebido %q (name=%q)",
					code, m.Vendor, m.Name)
			}
		case len(code) == 5 && code[0] == 'E' && isAllDigits(code[1:]):
			if m.Vendor != "ESUN" {
				t.Errorf("code=%q (prefixo E): vendor esperado \"ESUN\", recebido %q (name=%q)",
					code, m.Vendor, m.Name)
			}
		case len(code) == 5 && code[0] == 'P' && isAllDigits(code[1:]):
			if m.Vendor != "POLY" {
				t.Errorf("code=%q (prefixo P): vendor esperado \"POLY\", recebido %q (name=%q)",
					code, m.Vendor, m.Name)
			}
		default:
			t.Errorf("code=%q não obedece nenhuma regra de prefixo conhecida (name=%q, vendor=%q)",
				code, m.Name, m.Vendor)
		}
	}
}

// isAllDigits reporta se s contém apenas caracteres numéricos ASCII.
func isAllDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

// TestMaterialsNameVendorUnique valida que o par (Name, Vendor) é único em `materials`.
// Pares duplicados causariam ambiguidade no lookup nome→código em convertMaterial.
func TestMaterialsNameVendorUnique(t *testing.T) {
	type par struct{ name, vendor string }
	vistos := make(map[par]string, len(materials))
	for _, m := range materials {
		chave := par{m.Name, m.Vendor}
		if anterior, dup := vistos[chave]; dup {
			t.Errorf("par (name=%q, vendor=%q) duplicado: primeiro em code=%q, repetido em code=%q",
				m.Name, m.Vendor, anterior, m.Code)
			continue
		}
		vistos[chave] = m.Code
	}
}

// TestMaterialsAlphabeticalOrder valida que a slice `materials` está ordenada
// alfabeticamente por Name (case-insensitive, crescente).
func TestMaterialsAlphabeticalOrder(t *testing.T) {
	for i := 0; i < len(materials)-1; i++ {
		a := strings.ToLower(materials[i].Name)
		b := strings.ToLower(materials[i+1].Name)
		if a > b {
			t.Errorf("materials fora de ordem: materials[%d].Name=%q > materials[%d].Name=%q",
				i, materials[i].Name, i+1, materials[i+1].Name)
		}
	}
}

// TestMaterialsContainsExpectedEntries verifica que as 18 novas entradas e a
// entrada corrigida (00035) estão presentes com Name e Vendor exatos.
func TestMaterialsContainsExpectedEntries(t *testing.T) {
	esperadas := []struct {
		code   string
		name   string
		vendor string
	}{
		// Entrada corrigida
		{"00035", "Generic PLA-LW", "0000"},
		// 18 novas entradas
		{"00028", "Generic PA-GF", "0000"},
		{"00030", "Generic PP-GF", "0000"},
		{"E1007", "eSUN PLA+HS", "ESUN"},
		{"E1008", "eSUN PLA-LW", "ESUN"},
		{"E1009", "eSUN PLA-Basic", "ESUN"},
		{"E2004", "eSUN PETG-CF", "ESUN"},
		{"E3002", "eSUN ABS-CF", "ESUN"},
		{"E3003", "eSUN ABS+HS", "ESUN"},
		{"E5001", "eSUN TPU-95A", "ESUN"},
		{"P2001", "Fiberon PETG-ESD", "POLY"},
		{"P2002", "Fiberon PETG-rCF08", "POLY"},
		{"P7001", "Fiberon PA6-GF25", "POLY"},
		{"P7002", "Fiberon PA612-CF15", "POLY"},
		{"P7003", "Fiberon PA6-CF20", "POLY"},
		{"P7004", "Fiberon PA12-CF10", "POLY"},
		{"P7005", "Fiberon PA612-ESD", "POLY"},
		{"P8001", "Fiberon PET-CF17", "POLY"},
		{"P9001", "Fiberon PPS-CF10", "POLY"},
	}

	// Índice por código para acesso O(1)
	byCode := make(map[string]MaterialOption, len(materials))
	for _, m := range materials {
		byCode[m.Code] = m
	}

	for _, e := range esperadas {
		m, ok := byCode[e.code]
		if !ok {
			t.Errorf("code=%q não encontrado em materials (esperado name=%q, vendor=%q)",
				e.code, e.name, e.vendor)
			continue
		}
		if m.Name != e.name {
			t.Errorf("code=%q: name esperado %q, recebido %q", e.code, e.name, m.Name)
		}
		if m.Vendor != e.vendor {
			t.Errorf("code=%q: vendor esperado %q, recebido %q", e.code, e.vendor, m.Vendor)
		}
	}
}
