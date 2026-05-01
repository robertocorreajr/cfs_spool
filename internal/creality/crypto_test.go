package creality

import (
	"strings"
	"testing"
)

// payload38 é um payload ASCII válido de 38 bytes (formato real de tag CFS).
// Layout (slicing usado por ParseFields):
//   [0:5]   Date     "10524"
//   [5:9]   Supplier "0276"
//   [9:11]  Batch    "A2"
//   [11]    literal  "1"          (inserido por ASCIIConcat, ignorado pelo parse)
//   [12:17] Material "01001"
//   [17:24] Color    "0FF4010"
//   [24:28] Length   "0330"
//   [28:34] Serial   "000001"
//   [34:38] Reserve  "0000"
// Total: 5+4+2+1+5+7+4+6+4 = 38
const payload38 = "10524" + "0276" + "A2" + "1" + "01001" + "0FF4010" + "0330" + "000001" + "0000"

// TestDeriveS1KeyFromUID cobre o derivador da chave do setor 1 a partir do UID.
// O código atual aceita SOMENTE UIDs de 4 bytes (8 hex). UIDs de 7 ou 10 bytes
// retornam erro — esses casos estão na issue #22 mas hoje não são suportados.
func TestDeriveS1KeyFromUID(t *testing.T) {
	testes := []struct {
		nome     string
		uid      string
		ok       bool
		bytesOut int
	}{
		{"UID valido 4 bytes uppercase", "DEADBEEF", true, 12},
		{"UID valido 4 bytes lowercase", "deadbeef", true, 12},
		{"UID valido outro", "01020304", true, 12},
		{"UID vazio", "", false, 0},
		{"UID com 7 chars (impar)", "1234567", false, 0},
		{"UID com 9 chars", "DEADBEEF0", false, 0},
		{"UID 7 bytes (14 chars) nao suportado", "04A1B2C3D4E5F6", false, 0},
		{"UID 10 bytes (20 chars) nao suportado", "04A1B2C3D4E5F60718293A", false, 0},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			out, err := DeriveS1KeyFromUID(tt.uid)
			if tt.ok && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatalf("esperava erro para UID %q, recebeu %q", tt.uid, out)
			}
			if tt.ok {
				if len(out) != tt.bytesOut {
					t.Errorf("tamanho da chave = %d, esperado %d", len(out), tt.bytesOut)
				}
				if out != strings.ToUpper(out) {
					t.Errorf("chave deve ser uppercase, recebido %q", out)
				}
			}
		})
	}
}

// TestDeriveS1KeyDeterministic confirma que a mesma UID sempre produz a mesma
// chave (AES-ECB com chave fixa é determinístico) e que case do UID não afeta
// o resultado.
func TestDeriveS1KeyDeterministic(t *testing.T) {
	uid := "DEADBEEF"
	primeiro, err := DeriveS1KeyFromUID(uid)
	if err != nil {
		t.Fatalf("primeira chamada: %v", err)
	}
	segundo, err := DeriveS1KeyFromUID(uid)
	if err != nil {
		t.Fatalf("segunda chamada: %v", err)
	}
	if primeiro != segundo {
		t.Errorf("derivacao nao deterministica: %q != %q", primeiro, segundo)
	}

	upper, _ := DeriveS1KeyFromUID("DEADBEEF")
	lower, _ := DeriveS1KeyFromUID("deadbeef")
	if upper != lower {
		t.Errorf("case do UID afeta resultado: upper=%q lower=%q", upper, lower)
	}
}

// TestEncryptPayloadToBlocksTamanhos cobre os tamanhos aceitos (38 e 48 bytes)
// e tamanhos inválidos. Em todos os casos válidos a saída são 3 blocos de 32
// hex chars (96 hex total = 48 bytes cifrados).
func TestEncryptPayloadToBlocksTamanhos(t *testing.T) {
	testes := []struct {
		nome    string
		payload string
		ok      bool
	}{
		{"38 bytes (dados sem padding)", payload38, true},
		{"48 bytes (dados + padding manual)", payload38 + "0000000000", true},
		{"vazio", "", false},
		{"37 bytes", strings.Repeat("A", 37), false},
		{"39 bytes", strings.Repeat("A", 39), false},
		{"47 bytes", strings.Repeat("A", 47), false},
		{"49 bytes", strings.Repeat("A", 49), false},
		{"96 bytes (tamanho do hex cifrado por engano)", strings.Repeat("A", 96), false},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			b4, b5, b6, err := EncryptPayloadToBlocks(tt.payload)
			if tt.ok && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatalf("esperava erro, recebeu blocos")
			}
			if tt.ok {
				for i, b := range []string{b4, b5, b6} {
					if len(b) != 32 {
						t.Errorf("bloco %d tem %d chars, esperado 32", i+4, len(b))
					}
					if b != strings.ToUpper(b) {
						t.Errorf("bloco %d nao esta uppercase: %q", i+4, b)
					}
				}
			}
		})
	}
}

// TestEncryptDecryptRoundTrip valida que cifrar + decifrar retorna o payload
// original (com padding aplicado quando entrada tem 38 bytes).
func TestEncryptDecryptRoundTrip(t *testing.T) {
	t.Run("entrada 38 bytes vira 48 bytes apos round-trip", func(t *testing.T) {
		b4, b5, b6, err := EncryptPayloadToBlocks(payload38)
		if err != nil {
			t.Fatalf("encrypt: %v", err)
		}

		decifrado, err := DecryptBlocks(b4 + b5 + b6)
		if err != nil {
			t.Fatalf("decrypt: %v", err)
		}

		esperado := payload38 + "0000000000"
		if decifrado != esperado {
			t.Errorf("round-trip falhou\n  esperado: %q\n  recebido: %q", esperado, decifrado)
		}
	})

	t.Run("entrada 48 bytes preserva exato", func(t *testing.T) {
		original := payload38 + "ABCDEFGHIJ"
		b4, b5, b6, err := EncryptPayloadToBlocks(original)
		if err != nil {
			t.Fatalf("encrypt: %v", err)
		}

		decifrado, err := DecryptBlocks(b4 + b5 + b6)
		if err != nil {
			t.Fatalf("decrypt: %v", err)
		}

		if decifrado != original {
			t.Errorf("round-trip falhou\n  esperado: %q\n  recebido: %q", original, decifrado)
		}
	})
}

// TestEncryptDeterministic confirma que AES-ECB é determinístico — mesmo input
// produz mesmo output em chamadas repetidas (importante para a estabilidade da
// gravação na tag).
func TestEncryptDeterministic(t *testing.T) {
	b4a, b5a, b6a, err := EncryptPayloadToBlocks(payload38)
	if err != nil {
		t.Fatal(err)
	}
	b4b, b5b, b6b, err := EncryptPayloadToBlocks(payload38)
	if err != nil {
		t.Fatal(err)
	}
	if b4a != b4b || b5a != b5b || b6a != b6b {
		t.Errorf("AES-ECB nao deterministico para mesmo payload")
	}
}

// TestDecryptBlocks cobre tamanhos inválidos e hex malformado.
func TestDecryptBlocks(t *testing.T) {
	// Gerar um cifrado válido para o caso happy path
	b4, b5, b6, err := EncryptPayloadToBlocks(payload38)
	if err != nil {
		t.Fatal(err)
	}
	hexValido := b4 + b5 + b6

	testes := []struct {
		nome  string
		input string
		ok    bool
	}{
		{"96 hex chars uppercase", hexValido, true},
		{"96 hex chars lowercase", strings.ToLower(hexValido), true},
		{"vazio", "", false},
		{"95 chars", hexValido[:95], false},
		{"97 chars", hexValido + "F", false},
		{"32 chars (1 bloco apenas)", b4, false},
		{"hex invalido (com Z)", strings.Repeat("Z", 96), false},
		{"hex invalido (com espaco)", strings.Repeat("A", 95) + " ", false},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			out, err := DecryptBlocks(tt.input)
			if tt.ok && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatalf("esperava erro, recebeu %q", out)
			}
			if tt.ok && len(out) != 48 {
				t.Errorf("saida tem %d bytes, esperado 48", len(out))
			}
		})
	}
}

// TestParseFields valida o slicing dos campos a partir de um payload de 38
// bytes e os erros para entradas truncadas. Também faz round-trip via
// ASCIIConcat para garantir que ParseFields(ASCIIConcat(f)) preserva os campos.
func TestParseFields(t *testing.T) {
	t.Run("payload valido extrai todos os campos", func(t *testing.T) {
		f, err := ParseFields(payload38)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		// Layout: Date(5) Supplier(4) Batch(2) [1] Material(5) Color(7) Length(4) Serial(6) Reserve(4)
		// payload38 = "10524027620A210100100FF40100330000001A2A0"
		if f.Date != "10524" {
			t.Errorf("Date = %q, esperado %q", f.Date, "10524")
		}
		if f.Supplier != "0276" {
			t.Errorf("Supplier = %q, esperado %q", f.Supplier, "0276")
		}
		if f.Batch != "A2" {
			t.Errorf("Batch = %q, esperado %q", f.Batch, "A2")
		}
		if f.Material != "01001" {
			t.Errorf("Material = %q, esperado %q", f.Material, "01001")
		}
		if f.Color != "0FF4010" {
			t.Errorf("Color = %q, esperado %q", f.Color, "0FF4010")
		}
		if f.Length != "0330" {
			t.Errorf("Length = %q, esperado %q", f.Length, "0330")
		}
		if f.Serial != "000001" {
			t.Errorf("Serial = %q, esperado %q", f.Serial, "000001")
		}
		if f.Reserve != "0000" {
			t.Errorf("Reserve = %q, esperado %q", f.Reserve, "0000")
		}
	})

	t.Run("payload truncado retorna erro", func(t *testing.T) {
		testes := []string{
			"",
			"1",
			strings.Repeat("A", 37),
		}
		for _, in := range testes {
			if _, err := ParseFields(in); err == nil {
				t.Errorf("ParseFields(%q) deveria retornar erro", in)
			}
		}
	})

	t.Run("payload com 39+ bytes ignora excedente", func(t *testing.T) {
		// ParseFields aceita >= 38 bytes, ignorando o que vier depois
		_, err := ParseFields(payload38 + "EXTRADATA")
		if err != nil {
			t.Errorf("ParseFields aceita >38 bytes: %v", err)
		}
	})

	t.Run("round-trip Fields -> ASCIIConcat -> ParseFields", func(t *testing.T) {
		original := Fields{
			Batch:    "A2",
			Date:     "10524",
			Supplier: "0276",
			Material: "01001",
			Color:    "0FF4010",
			Length:   "0330",
			Serial:   "000001",
			Reserve:  "0000",
		}
		ascii, err := original.ASCIIConcat()
		if err != nil {
			t.Fatalf("ASCIIConcat: %v", err)
		}
		// ASCIIConcat insere "1" literal entre Batch e Material, então o slicing
		// de ParseFields pega Material a partir da posição 12 (após o "1"). O
		// campo Batch parseado vai ser "A2" e o Material "01001".
		volta, err := ParseFields(ascii)
		if err != nil {
			t.Fatalf("ParseFields: %v", err)
		}
		if volta.Date != original.Date ||
			volta.Supplier != original.Supplier ||
			volta.Material != original.Material ||
			volta.Color != original.Color ||
			volta.Length != original.Length ||
			volta.Serial != original.Serial ||
			volta.Reserve != original.Reserve {
			t.Errorf("round-trip alterou campos:\n  original: %+v\n  volta:    %+v", original, volta)
		}
	})
}

// TestSetColor cobre a validação de tamanho exato (6 chars) feita por SetColor
// e o prefixo "0" que ele adiciona ao campo Color.
func TestSetColor(t *testing.T) {
	testes := []struct {
		nome      string
		entrada   string
		ok        bool
		colorOut  string // só checado se ok
	}{
		{"6 chars hex maiusculo", "FF4010", true, "0FF4010"},
		{"6 chars hex minusculo (sem normalizacao)", "ff4010", true, "0ff4010"},
		{"6 chars zeros", "000000", true, "0000000"},
		{"6 chars todos F", "FFFFFF", true, "0FFFFFF"},
		{"vazio", "", false, ""},
		{"5 chars", "FF401", false, ""},
		{"7 chars", "FF40100", false, ""},
		{"8 chars", "FF401000", false, ""},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			f := Fields{}
			err := f.SetColor(tt.entrada)
			if tt.ok && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatalf("esperava erro para %q", tt.entrada)
			}
			if tt.ok && f.Color != tt.colorOut {
				t.Errorf("Color = %q, esperado %q", f.Color, tt.colorOut)
			}
		})
	}
}
