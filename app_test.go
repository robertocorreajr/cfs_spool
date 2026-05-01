package main

import (
	"testing"
)

func TestValidateColor(t *testing.T) {
	app := NewApp()

	testes := []struct {
		entrada  string
		esperado string
		erro     bool
	}{
		{"FF4010", "FF4010", false},
		{"ff4010", "FF4010", false},
		{"000000", "000000", false},
		{"FFFFFF", "FFFFFF", false},
		{"aAbBcC", "AABBCC", false},
		{"#FF4010", "FF4010", false},  // com #
		{" FF4010 ", "FF4010", false}, // com espaços
		{"GG0000", "", true},          // caracteres inválidos
		{"12345", "", true},           // muito curto
		{"1234567", "", true},         // muito longo
		{"", "", true},                // vazio
		{"ZZZZZZ", "", true},          // não-hex
	}

	for _, tt := range testes {
		resultado, err := app.ValidateColor(tt.entrada)
		if tt.erro && err == nil {
			t.Errorf("ValidateColor(%q) deveria retornar erro", tt.entrada)
		}
		if !tt.erro && err != nil {
			t.Errorf("ValidateColor(%q) retornou erro inesperado: %v", tt.entrada, err)
		}
		if resultado != tt.esperado {
			t.Errorf("ValidateColor(%q) = %q, esperado %q", tt.entrada, resultado, tt.esperado)
		}
	}
}

func TestConvertDate(t *testing.T) {
	testes := []struct {
		entrada  string
		esperado string
		erro     bool
	}{
		{"2026-04-12", "26412", false},
		{"2024-01-05", "24105", false},
		{"2024-11-15", "24B15", false}, // mês 11 > 9, deve ser representado como 1-dígito
		{"2024-12-25", "24C25", false}, // mês 12 > 9
		{"invalido", "", true},
	}

	for _, tt := range testes {
		resultado, err := convertDate(tt.entrada)
		if tt.erro && err == nil {
			t.Errorf("convertDate(%q) deveria retornar erro", tt.entrada)
		}
		if !tt.erro && err != nil {
			t.Errorf("convertDate(%q) retornou erro inesperado: %v", tt.entrada, err)
		}
		if !tt.erro && resultado != tt.esperado {
			t.Errorf("convertDate(%q) = %q, esperado %q", tt.entrada, resultado, tt.esperado)
		}
	}
}

func TestParseDateToISO(t *testing.T) {
	testes := []struct {
		entrada  string
		esperado string
	}{
		{"26412", "2026-04-12"},
		{"24105", "2024-01-05"},
		{"25915", "2025-09-15"},
		{"24A20", "2024-10-20"},
		{"24B15", "2024-11-15"},
		{"24C25", "2024-12-25"},
		{"", ""},          // vazio
		{"123", ""},       // muito curto
		{"12345X", ""},    // muito longo
	}

	for _, tt := range testes {
		resultado := parseDateToISO(tt.entrada)
		if resultado != tt.esperado {
			t.Errorf("parseDateToISO(%q) = %q, esperado %q", tt.entrada, resultado, tt.esperado)
		}
	}
}

func TestPadSerial(t *testing.T) {
	testes := []struct {
		entrada  string
		esperado string
	}{
		{"", "000001"},
		{"1", "000001"},
		{"123", "000123"},
		{"000001", "000001"},
		{"1234567", "123456"}, // trunca
		{" 42 ", "000042"},   // espaços
	}

	for _, tt := range testes {
		resultado := padSerial(tt.entrada)
		if resultado != tt.esperado {
			t.Errorf("padSerial(%q) = %q, esperado %q", tt.entrada, resultado, tt.esperado)
		}
	}
}

func TestConvertDateRoundTrip(t *testing.T) {
	// Testar que convertDate + parseDateToISO é ida e volta
	datas := []string{
		"2026-04-12",
		"2024-01-05",
		"2025-09-15",
	}

	for _, data := range datas {
		interno, err := convertDate(data)
		if err != nil {
			t.Fatalf("convertDate(%q) erro: %v", data, err)
		}
		volta := parseDateToISO(interno)
		if volta != data {
			t.Errorf("round-trip falhou: %q -> %q -> %q", data, interno, volta)
		}
	}
}

// TestConvertLength documenta o comportamento ATUAL de convertLength.
//
// Atenção: a função tem um bug de short-circuit — entradas com 4 chars
// passam direto SEM consultar o lookup. Isso significa que gramaturas
// como "1000" e "2000" (que estão no gramMap) nunca alcançam o lookup
// e são gravadas literalmente como código hex de 4 chars.
// Bug rastreado na issue #39 — quando corrigida, atualizar os esperados
// marcados com "(BUG:...)" abaixo.
func TestConvertLength(t *testing.T) {
	testes := []struct {
		nome     string
		entrada  string
		esperado string
	}{
		// 4 chars: passa direto (intenção parece ser "já está em formato RFID")
		{"4 chars passa direto", "0053", "0053"},
		{"4 chars uppercase", "ABCD", "ABCD"},
		{"1000g passa direto (BUG: deveria ir pro lookup gramMap)", "1000", "1000"},
		{"2000g passa direto (BUG: deveria ir pro lookup gramMap)", "2000", "2000"},
		{"3000g passa direto (BUG: deveria virar 03E8)", "3000", "3000"},
		{"5000g passa direto (BUG: deveria virar 0682)", "5000", "5000"},

		// Lookup cm pré-definido (4 chars com prefixo 0)
		// Nota: "0083", "0165", "0330", "0660" também têm 4 chars, então passam direto.
		// O lookup lengthMap nunca é alcançado pelo código atual.
		{"0083 passa direto (4 chars, NAO consulta lengthMap)", "0083", "0083"},
		{"0165 passa direto (4 chars, NAO consulta lengthMap)", "0165", "0165"},
		{"0330 passa direto (4 chars, NAO consulta lengthMap)", "0330", "0330"},
		{"0660 passa direto (4 chars, NAO consulta lengthMap)", "0660", "0660"},

		// Lookup gramas pré-definido (3 chars, alcança gramMap)
		{"250g vira 0053", "250", "0053"},
		{"500g vira 00A5", "500", "00A5"},

		// Custom: cm = gramas/3 → hex (entradas com 5+ chars)
		{"10000g vira 0D05 (3333cm)", "10000", "0D05"},
		{"15000g vira 1388 (5000cm)", "15000", "1388"},

		// Overflow: cap em 65535 cm
		{"200000g cap em FFFF", "200000", "FFFF"},
		{"196605g (65535*3) bate exato em FFFF", "196605", "FFFF"},

		// Inválidos: fallback "0053"
		{"vazio fallback (0 chars)", "", "0053"},
		{"texto invalido fallback (3 chars)", "abc", "0053"},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			resultado := convertLength(tt.entrada)
			if resultado != tt.esperado {
				t.Errorf("convertLength(%q) = %q, esperado %q", tt.entrada, resultado, tt.esperado)
			}
		})
	}
}

// TestConvertMaterial cobre os 3 caminhos: já em código (5 chars passa direto),
// nome conhecido (lookup), nome desconhecido (passa direto como fallback).
func TestConvertMaterial(t *testing.T) {
	testes := []struct {
		nome     string
		entrada  string
		esperado string
	}{
		{"codigo de 5 chars passa direto", "00001", "00001"},
		{"codigo eSUN passa direto", "E1001", "E1001"},
		{"codigo Polymaker passa direto", "P1001", "P1001"},
		{"PLA generico", "PLA", "00001"},
		{"Hyper PLA Creality", "Hyper PLA", "01001"},
		{"eSUN PLA+", "eSUN PLA+", "E1001"},
		{"PolySonic PLA", "PolySonic PLA", "P1004"},
		{"desconhecido passa direto", "Material X", "Material X"},
		{"vazio passa direto", "", ""},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			resultado := convertMaterial(tt.entrada)
			if resultado != tt.esperado {
				t.Errorf("convertMaterial(%q) = %q, esperado %q", tt.entrada, resultado, tt.esperado)
			}
		})
	}
}

// TestVendorToSupplier confirma que apenas "0000" mapeia para genérico;
// qualquer outro vendor (Creality, eSUN, Polymaker) vira "0276" no RFID.
func TestVendorToSupplier(t *testing.T) {
	testes := []struct {
		nome     string
		entrada  string
		esperado string
	}{
		{"generico", "0000", "0000"},
		{"Creality", "0276", "0276"},
		{"eSUN vira 0276 no RFID", "ESUN", "0276"},
		{"Polymaker vira 0276 no RFID", "POLY", "0276"},
		{"vazio vira 0276 (default)", "", "0276"},
		{"desconhecido vira 0276", "XYZ", "0276"},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			resultado := vendorToSupplier(tt.entrada)
			if resultado != tt.esperado {
				t.Errorf("vendorToSupplier(%q) = %q, esperado %q", tt.entrada, resultado, tt.esperado)
			}
		})
	}
}

// TestMaterialToVendor cobre os 4 caminhos:
// - prefixo E → "ESUN"
// - prefixo P → "POLY"
// - código numérico no range 01000-29999 → "0276" (Creality)
// - resto → "0000" (genérico)
func TestMaterialToVendor(t *testing.T) {
	testes := []struct {
		nome     string
		entrada  string
		esperado string
	}{
		{"prefixo E (eSUN)", "E1001", "ESUN"},
		{"prefixo P (Polymaker)", "P1001", "POLY"},
		{"limite inferior Creality (01000)", "01000", "0276"},
		{"meio do range Creality (15001)", "15001", "0276"},
		{"limite superior Creality (29999)", "29999", "0276"},
		{"abaixo do range (00001 generico)", "00001", "0000"},
		{"acima do range (30000 generico)", "30000", "0000"},
		{"vazio vira generico", "", "0000"},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			resultado := materialToVendor(tt.entrada)
			if resultado != tt.esperado {
				t.Errorf("materialToVendor(%q) = %q, esperado %q", tt.entrada, resultado, tt.esperado)
			}
		})
	}
}

// TestVendorName cobre os 4 vendors conhecidos e o fallback "(desconhecido)".
func TestVendorName(t *testing.T) {
	testes := []struct {
		nome     string
		entrada  string
		esperado string
	}{
		{"Creality", "0276", "Creality"},
		{"Generico", "0000", "Genérico"},
		{"eSUN", "ESUN", "eSUN"},
		{"Polymaker", "POLY", "Polymaker"},
		{"desconhecido", "XYZ", "XYZ (desconhecido)"},
		{"vazio", "", " (desconhecido)"},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			resultado := vendorName(tt.entrada)
			if resultado != tt.esperado {
				t.Errorf("vendorName(%q) = %q, esperado %q", tt.entrada, resultado, tt.esperado)
			}
		})
	}
}
