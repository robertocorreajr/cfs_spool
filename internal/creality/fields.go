package creality

import (
	"errors"
	"fmt"
	"strings"
)

// tamanhos em bytes ASCII
const (
	lenBatch    = 2  // A2 (fixo)
	lenDate     = 5
	lenSupplier = 4
	lenMaterial = 5
	lenColor    = 6  // 0 + 5 caracteres hex
	lenLength   = 4
	lenSerial   = 6
	lenReserve  = 6  // 000000 (fixo)
)

type Fields struct {
	Batch, Date, Supplier, Material, Color, Length, Serial, Reserve string
}

// NewFields cria uma nova instância de Fields com valores fixos para Batch e Reserve
func NewFields() Fields {
	return Fields{
		Batch:   "A2",     // Valor fixo
		Reserve: "000000", // Valor fixo
	}
}

// SetBatchFixed força o campo Batch para o valor fixo "A2"
func (f *Fields) SetBatchFixed() {
	f.Batch = "A2"
}

// SetReserveFixed força o campo Reserve para o valor fixo "000000"
func (f *Fields) SetReserveFixed() {
	f.Reserve = "000000"
}

// SetColor define o campo Color garantindo que sempre comece com "0"
// color deve ser uma string hex de 5 caracteres (sem #)
func (f *Fields) SetColor(color string) error {
	if len(color) != 5 {
		return errors.New("cor deve ter exatamente 5 caracteres hex")
	}
	f.Color = "0" + color
	return nil
}

// ValidateAndFix valida e corrige automaticamente os campos obrigatórios
func (f *Fields) ValidateAndFix() {
	// Força valores fixos
	f.SetBatchFixed()
	f.SetReserveFixed()
	
	// Garante que Color comece com "0" se não estiver vazio
	if f.Color != "" && len(f.Color) == 6 && f.Color[0] != '0' {
		// Se Color tem 6 caracteres mas não começa com 0, adiciona o 0
		if len(f.Color) == 5 {
			f.Color = "0" + f.Color
		}
	}
}

func (f Fields) ASCIIConcat() (string, error) {
	// Cria uma cópia para não modificar o original
	fields := f
	
	// Aplica validação e correções automáticas
	fields.ValidateAndFix()
	
	if len(fields.Batch) != lenBatch ||
		len(fields.Date) != lenDate ||
		len(fields.Supplier) != lenSupplier ||
		len(fields.Material) != lenMaterial ||
		len(fields.Color) != lenColor ||
		len(fields.Length) != lenLength ||
		len(fields.Serial) != lenSerial ||
		len(fields.Reserve) != lenReserve {
		return "", errors.New("algum campo está com tamanho incorreto")
	}
	return fmt.Sprintf("%s%s%s%s%s%s%s%s",
		fields.Batch, fields.Date, fields.Supplier, fields.Material,
		fields.Color, fields.Length, fields.Serial, fields.Reserve), nil
}

// ASCIIConcat48 gera payload de 48 bytes (38 dados + 10 padding) para compatibilidade
func (f Fields) ASCIIConcat48() (string, error) {
	// Obter payload de 38 bytes
	payload38, err := f.ASCIIConcat()
	if err != nil {
		return "", err
	}
	
	// Adicionar 10 bytes de padding com zeros
	return payload38 + "0000000000", nil
}

// ParseFields extrai os campos de uma string ASCII de 38 bytes
func ParseFields(ascii38 string) (Fields, error) {
	if len(ascii38) != 38 {
		return Fields{}, errors.New("string ASCII deve ter exatamente 38 bytes")
	}

	// Extrair cada campo baseado no tamanho fixo
	fields := Fields{
		Batch:    ascii38[0:2],                    // 2 bytes (A2)
		Date:     ascii38[2:7],                    // 5 bytes  
		Supplier: ascii38[7:11],                   // 4 bytes
		Material: ascii38[11:16],                  // 5 bytes
		Color:    ascii38[16:22],                  // 6 bytes (0 + 5 hex)
		Length:   ascii38[22:26],                  // 4 bytes
		Serial:   ascii38[26:32],                  // 6 bytes
		Reserve:  ascii38[32:38],                  // 6 bytes (000000)
	}

	return fields, nil
}

// ParseFieldsCompat extrai os campos com compatibilidade para formatos antigos
func ParseFieldsCompat(data string) (Fields, error) {
	// Se tem exatamente 38 bytes, usar formato novo
	if len(data) == 38 {
		return ParseFields(data)
	}
	
	// Se tem 48 bytes, é formato antigo - extrair apenas os primeiros 38 bytes
	if len(data) == 48 {
		return ParseFields(data[:38])
	}
	
	// Se tem menos de 38 bytes, completar com zeros
	if len(data) < 38 {
		padded := data + strings.Repeat("0", 38-len(data))
		return ParseFields(padded)
	}
	
	// Se tem mais de 48 bytes, truncar para 38
	if len(data) > 48 {
		return ParseFields(data[:38])
	}
	
	return Fields{}, fmt.Errorf("tamanho de dados não suportado: %d bytes", len(data))
}

// String retorna uma representação legível dos campos
func (f Fields) String() string {
	return fmt.Sprintf("Batch: %s, Date: %s, Supplier: %s, Material: %s, Color: %s, Length: %s, Serial: %s, Reserve: %s",
		f.Batch, f.Date, f.Supplier, f.Material, f.Color, f.Length, f.Serial, f.Reserve)
}

// FormatDate converte a data do formato YYMDD para formato legível
func (f Fields) FormatDate() string {
	if len(f.Date) != 5 {
		return f.Date + " (formato inválido)"
	}
	
	year := "20" + f.Date[0:2]
	month := f.Date[2:3]
	day := f.Date[3:5]
	
	monthNames := map[string]string{
		"1": "Janeiro", "2": "Fevereiro", "3": "Março", "4": "Abril",
		"5": "Maio", "6": "Junho", "7": "Julho", "8": "Agosto",
		"9": "Setembro", "10": "Outubro", "11": "Novembro", "12": "Dezembro",
	}
	
	monthName := monthNames[month]
	if monthName == "" {
		monthName = "Mês " + month
	}
	
	return fmt.Sprintf("%s de %s de %s", day, monthName, year)
}

// FormatColor converte a cor para formato legível
func (f Fields) FormatColor() string {
	if len(f.Color) == 6 && f.Color[0] == '0' {
		return "#" + f.Color[1:] + " (hex)"
	}
	return f.Color
}

// FormatLength interpreta o comprimento
func (f Fields) FormatLength() string {
	switch f.Length {
	case "0330":
		return "330cm (1kg de filamento)"
	case "0165":
		return "165cm (500g de filamento)"
	case "0083":
		return "83cm (250g de filamento)"
	default:
		return f.Length + "cm"
	}
}

// GetMaterialName retorna o nome do material baseado no código
func (f Fields) GetMaterialName() string {
	materials := map[string]string{
		"00001": "Generic PLA",
		"00002": "Generic PLA-Silk",
		"00003": "Generic PETG",
		"00004": "Generic ABS",
		"00005": "Generic TPU",
		"00006": "Generic PLA-CF",
		"00007": "Generic ASA",
		"00008": "PA Genérico",
		"00009": "PA-CF Genérico",
		"00010": "BVOH Genérico",
		"00012": "HIPS Genérico",
		"00013": "PET-CF Genérico",
		"00014": "PETG-CF Genérico",
		"00015": "PA6-CF Genérico",
		"00016": "PAHT-CF Genérico",
		"00020": "PET Genérico",
		"00021": "PC Genérico",
		"01001": "Hyper PLA",
		"02001": "Hyper PLA-CF",
		"03001": "Hyper ABS",
		"04001": "CR-PLA",
		"05001": "CR-Silk",
		"06001": "CR-PETG",
		"07001": "CR-ABS",
		"08001": "Ender-PLA",
		"09001": "EN-PLA+",
		"09002": "ENDERFASTPLA",
		"10001": "HP-TPU",
		"11001": "CR-Nylon",
		"13001": "CR-PLACarbon",
		"14001": "CR-PLAMatte",
		"15001": "CR-PLAFluo",
		"16001": "CR-TPU",
		"17001": "CR-Wood",
		"18001": "HPUltraPLA",
		"19001": "HP-ASA",
	}
	
	name := materials[f.Material]
	if name == "" {
		return f.Material + " (desconhecido)"
	}
	return name
}

// GetSupplierName retorna o nome do fornecedor baseado no código
func (f Fields) GetSupplierName() string {
	suppliers := map[string]string{
		"0276": "Creality",
		"0000": "Genérico",
	}
	
	name := suppliers[f.Supplier]
	if name == "" {
		return f.Supplier + " (desconhecido)"
	}
	return name
}

// TODO: EncryptPayloadToBlocks(payload string) (b4,b5,b6 string, error)
// (uso dos mesmos segredos AES/ECB do site)
