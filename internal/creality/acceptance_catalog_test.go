package creality

// acceptance_catalog_test.go — testes de aceitação para a expansão do catálogo v7.
//
// Critério coberto: "Para CADA das 18 entradas, o triple (code, name, vendor) existe e
// gravação/leitura preservam" (round-trip lossless via ASCIIConcat + ParseFieldsCompat +
// GetMaterialName).  Também cobre o critério de round-trip de codes alfanuméricos
// (E1009, P9001).
//
// Não duplica TestRoundTripMaterialE1009 — este teste abrange todos os 18 novos codes +
// o 00035 corrigido numa tabela orientada a dados.

import "testing"

// TestRoundTripAllNewMaterials valida o ciclo completo encode→decode para cada uma das
// 18 novas entradas do catálogo v7 e para a entrada corrigida 00035.
// Garante que ASCIIConcat + ParseFieldsCompat preservam o código do material e que
// GetMaterialName retorna o nome exato esperado para cada entrada.
func TestRoundTripAllNewMaterials(t *testing.T) {
	testes := []struct {
		code string
		name string
	}{
		// Entrada corrigida
		{"00035", "Generic PLA-LW"},
		// 18 novas entradas
		{"00028", "Generic PA-GF"},
		{"00030", "Generic PP-GF"},
		{"E1007", "eSUN PLA+HS"},
		{"E1008", "eSUN PLA-LW"},
		{"E1009", "eSUN PLA-Basic"},
		{"E2004", "eSUN PETG-CF"},
		{"E3002", "eSUN ABS-CF"},
		{"E3003", "eSUN ABS+HS"},
		{"E5001", "eSUN TPU-95A"},
		{"P2001", "Fiberon PETG-ESD"},
		{"P2002", "Fiberon PETG-rCF08"},
		{"P7001", "Fiberon PA6-GF25"},
		{"P7002", "Fiberon PA612-CF15"},
		{"P7003", "Fiberon PA6-CF20"},
		{"P7004", "Fiberon PA12-CF10"},
		{"P7005", "Fiberon PA612-ESD"},
		{"P8001", "Fiberon PET-CF17"},
		{"P9001", "Fiberon PPS-CF10"},
	}

	for _, tt := range testes {
		tt := tt // captura para evitar data race em subtestes paralelos
		t.Run(tt.code, func(t *testing.T) {
			f := Fields{
				Batch:    "A2",
				Date:     "10524",
				Supplier: "0276",
				Material: tt.code,
				Color:    "0FF4010",
				Length:   "0330",
				Serial:   "000001",
				Reserve:  "0000",
			}

			ascii, err := f.ASCIIConcat()
			if err != nil {
				t.Fatalf("ASCIIConcat falhou para code=%q: %v", tt.code, err)
			}

			parsed, err := ParseFieldsCompat(ascii)
			if err != nil {
				t.Fatalf("ParseFieldsCompat falhou para code=%q: %v", tt.code, err)
			}

			if parsed.Material != tt.code {
				t.Errorf("Material após round-trip = %q, esperado %q", parsed.Material, tt.code)
			}

			nome := parsed.GetMaterialName()
			if nome != tt.name {
				t.Errorf("GetMaterialName() = %q, esperado %q (code=%q)", nome, tt.name, tt.code)
			}
		})
	}
}
