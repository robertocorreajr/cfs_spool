package creality

import "testing"

// TestFormatDate documenta o comportamento ATUAL de FormatDate.
//
// Atenção: a função tem um bug conhecido — extrai apenas 1 caractere do dia
// (f.Date[1]), não 2. Isso faz com que datas como "10524" (esperado "5 de
// Janeiro de 2024") sejam interpretadas como "0 de Janeiro de 2024".
// Os testes abaixo refletem o comportamento atual; quando a função for
// corrigida, atualizar os esperados aqui também.
func TestFormatDate(t *testing.T) {
	testes := []struct {
		nome     string
		date     string
		esperado string
	}{
		{"mes 1 (Janeiro), 2o char dia=0", "10524", "0 de Janeiro de 2024"},
		{"mes 9 (Setembro), 2o char dia=1", "91524", "1 de Setembro de 2024"},
		{"mes A (Outubro), 2o char dia=2", "A2024", "2 de Outubro de 2024"},
		{"mes B (Novembro), 2o char dia=B (=11)", "BB124", "11 de Novembro de 2024"},
		{"mes C (Dezembro), 2o char dia=2", "C2524", "2 de Dezembro de 2024"},
		{"mes desconhecido Z, 2o char dia=0", "Z0123", "0 de Mês Z de 2023"},
		{"dia base 36 (Z=35), mes 1", "1Z124", "35 de Janeiro de 2024"},
		{"date vazia", "", " (formato inválido)"},
		{"date curta", "12", "12 (formato inválido)"},
		{"date longa", "123456", "123456 (formato inválido)"},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			f := Fields{Date: tt.date}
			resultado := f.FormatDate()
			if resultado != tt.esperado {
				t.Errorf("FormatDate(%q) = %q, esperado %q", tt.date, resultado, tt.esperado)
			}
		})
	}
}

// TestFormatColor cobre cor válida (7 chars, começando com 0) e variações
// que devem retornar a cor crua sem formatação hex.
func TestFormatColor(t *testing.T) {
	testes := []struct {
		nome     string
		color    string
		esperado string
	}{
		{"cor valida", "0FF4010", "#FF4010 (hex)"},
		{"cor preta", "0000000", "#000000 (hex)"},
		{"cor branca", "0FFFFFF", "#FFFFFF (hex)"},
		{"sem o 0 prefixo", "FF40100", "FF40100"},
		{"6 chars sem prefixo", "FF4010", "FF4010"},
		{"vazio", "", ""},
		{"5 chars", "FF401", "FF401"},
		{"8 chars", "0FF40100", "0FF40100"},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			f := Fields{Color: tt.color}
			resultado := f.FormatColor()
			if resultado != tt.esperado {
				t.Errorf("FormatColor(%q) = %q, esperado %q", tt.color, resultado, tt.esperado)
			}
		})
	}
}

// TestFormatLength cobre os comprimentos pré-definidos e o fallback genérico
// para comprimentos customizados (que hoje exibe a string crua + "cm").
func TestFormatLength(t *testing.T) {
	testes := []struct {
		nome     string
		length   string
		esperado string
	}{
		{"330 (1kg)", "0330", "330cm (1kg de filamento)"},
		{"165 (500g)", "0165", "165cm (500g de filamento)"},
		{"83 (250g)", "0083", "83cm (250g de filamento)"},
		{"customizado hex 03E8 = 1000", "03E8", "03E8cm"},
		{"customizado decimal", "1234", "1234cm"},
		{"vazio", "", "cm"},
		{"660 (2kg) — não está no switch", "0660", "0660cm"},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			f := Fields{Length: tt.length}
			resultado := f.FormatLength()
			if resultado != tt.esperado {
				t.Errorf("FormatLength(%q) = %q, esperado %q", tt.length, resultado, tt.esperado)
			}
		})
	}
}

// TestGetMaterialName cobre lookups conhecidos (todas as marcas) e códigos
// inexistentes que devem cair no fallback "(desconhecido)".
func TestGetMaterialName(t *testing.T) {
	testes := []struct {
		nome     string
		material string
		esperado string
	}{
		{"PLA generico", "00001", "PLA"},
		{"Hyper PLA Creality", "01001", "Hyper PLA"},
		{"eSUN PLA+", "E1001", "eSUN PLA+"},
		{"Polymaker", "P1001", "Panchroma PLA Satin"},
		{"codigo desconhecido", "ZZZZZ", "ZZZZZ (desconhecido)"},
		{"vazio", "", " (desconhecido)"},
		{"codigo num intervalo nao mapeado", "99999", "99999 (desconhecido)"},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			f := Fields{Material: tt.material}
			resultado := f.GetMaterialName()
			if resultado != tt.esperado {
				t.Errorf("GetMaterialName(%q) = %q, esperado %q", tt.material, resultado, tt.esperado)
			}
		})
	}
}

// TestGetSupplierName cobre os dois códigos válidos no RFID (0276 e 0000)
// e códigos inválidos. O comentário no código deixa claro que marcas como
// eSUN/Polymaker são identificadas pelo material, não pelo supplier.
func TestGetSupplierName(t *testing.T) {
	testes := []struct {
		nome     string
		supplier string
		esperado string
	}{
		{"Creality", "0276", "Creality"},
		{"Generico", "0000", "Genérico"},
		{"ESUN nao deve aparecer no supplier", "ESUN", "ESUN (desconhecido)"},
		{"POLY nao deve aparecer no supplier", "POLY", "POLY (desconhecido)"},
		{"vazio", "", " (desconhecido)"},
		{"codigo numerico inexistente", "1234", "1234 (desconhecido)"},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			f := Fields{Supplier: tt.supplier}
			resultado := f.GetSupplierName()
			if resultado != tt.esperado {
				t.Errorf("GetSupplierName(%q) = %q, esperado %q", tt.supplier, resultado, tt.esperado)
			}
		})
	}
}

// tagValida retorna uma estrutura Fields completa que deve passar em todas
// as heurísticas de IsBlankTag.
func tagValida() Fields {
	return Fields{
		Batch:    "A2",
		Date:     "10524",
		Supplier: "0276",
		Material: "01001",
		Color:    "0FF4010",
		Length:   "0330",
		Serial:   "000001",
		Reserve:  "0000",
	}
}

// TestIsBlankTag cobre as heurísticas de detecção de tag virgem/inválida:
// batch, supplier, material (tamanho e charset), reserve e color.
func TestIsBlankTag(t *testing.T) {
	testes := []struct {
		nome     string
		ajuste   func(*Fields)
		esperado bool
	}{
		{"tag valida Creality", func(f *Fields) {}, false},
		{"tag valida Generico", func(f *Fields) { f.Supplier = "0000" }, false},
		{"tag valida material eSUN (E1001)", func(f *Fields) { f.Material = "E1001" }, false},
		{"tag valida material Polymaker (P1001)", func(f *Fields) { f.Material = "P1001" }, false},
		{"batch vazio", func(f *Fields) { f.Batch = "" }, true},
		{"batch errado", func(f *Fields) { f.Batch = "B1" }, true},
		{"supplier desconhecido", func(f *Fields) { f.Supplier = "9999" }, true},
		{"material com 4 chars", func(f *Fields) { f.Material = "0000" }, true},
		{"material com 6 chars", func(f *Fields) { f.Material = "000010" }, true},
		{"material com caractere invalido", func(f *Fields) { f.Material = "0000-" }, true},
		{"reserve diferente de 0000", func(f *Fields) { f.Reserve = "FFFF" }, true},
		{"color sem prefixo 0", func(f *Fields) { f.Color = "1FF4010" }, true},
		{"color com 6 chars", func(f *Fields) { f.Color = "FF4010" }, true},
		{"color com 8 chars", func(f *Fields) { f.Color = "0FF40100" }, true},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			f := tagValida()
			tt.ajuste(&f)
			resultado := f.IsBlankTag()
			if resultado != tt.esperado {
				t.Errorf("IsBlankTag() para tag %+v = %v, esperado %v", f, resultado, tt.esperado)
			}
		})
	}
}
