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
