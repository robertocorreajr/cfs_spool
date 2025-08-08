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
	lenColor    = 7  // 0 + 6 caracteres hex
	lenLength   = 4
	lenSerial   = 6
	lenReserve  = 4  // 0000 (fixo) - 4 bytes na nova estrutura
)

type Fields struct {
	Batch, Date, Supplier, Material, Color, Length, Serial, Reserve string
}

// NewFields cria uma nova instância de Fields com valores fixos para Batch e Reserve
func NewFields() Fields {
	return Fields{
		Batch:   "A2",   // Valor fixo
		Reserve: "0000", // Valor fixo
	}
}

// SetBatchFixed força o campo Batch para o valor fixo "A2"
func (f *Fields) SetBatchFixed() {
	f.Batch = "A2"
}

// SetReserveFixed força o campo Reserve para o valor fixo "0000"
func (f *Fields) SetReserveFixed() {
	f.Reserve = "0000"
}

// SetColor define o campo Color garantindo que sempre tenha 7 bytes (0 + 6 hex)
// color deve ser uma string hex de 6 caracteres (sem #)
func (f *Fields) SetColor(color string) error {
	if len(color) != 6 {
		return errors.New("cor deve ter exatamente 6 caracteres hex")
	}
	f.Color = "0" + color
	return nil
}

// ValidateAndFix valida e corrige automaticamente os campos obrigatórios
func (f *Fields) ValidateAndFix() {
	// Força valores fixos
	f.SetBatchFixed()
	f.SetReserveFixed()
	
	// Garante que Color tenha 7 bytes se não estiver vazio
	if f.Color != "" && len(f.Color) == 6 && f.Color[0] != '0' {
		// Se Color tem 6 caracteres mas não começa com 0, adiciona o 0
		f.Color = "0" + f.Color
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
	return fmt.Sprintf("%s%s%s1%s%s%s%s%s",
		fields.Date, fields.Supplier, fields.Batch, fields.Material,
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
	if len(ascii38) < 38 {
		return Fields{}, errors.New("string ASCII deve ter pelo menos 38 bytes")
	}

	// Extrair cada campo baseado na tabela precisa do usuário:
	// BB1240276A2101001000000001650000010000
	// date(5) + venderId(4) + batch(2) + filamentId(6) + color(7) + filamentLen(4) + serialNum(6) + reserve(4) = 38 bytes
	if len(ascii38) < 38 {
		return Fields{}, fmt.Errorf("string ASCII deve ter pelo menos 38 bytes, recebido: %d", len(ascii38))
	}
	
	filamentId := ascii38[12:17]     // 5 bytes: 01001 (ignorando PRIMEIRO dígito do campo de 6 bytes)
	
	fields := Fields{
		Date:     ascii38[0:5],                    // 5 bytes - date: BB124
		Supplier: ascii38[5:9],                    // 4 bytes - venderId: 0276  
		Batch:    ascii38[9:11],                   // 2 bytes - batch: A2
		Material: filamentId,                      // 5 bytes - filamentId: 01001 (primeiro dígito ignorado)
		Color:    ascii38[17:24],                  // 7 bytes - color: 0000000 (0 fixo + 6 hex)
		Length:   ascii38[24:28],                  // 4 bytes - filamentLen: 0165
		Serial:   ascii38[28:34],                  // 6 bytes - serialNum: 000001
		Reserve:  ascii38[34:38],                  // 4 bytes - reserve: 0000
	}

	return fields, nil
}

// ParseFieldsCompat extrai os campos com compatibilidade para formatos antigos
func ParseFieldsCompat(data string) (Fields, error) {
	// Se tem exatamente 38 bytes, usar formato novo
	if len(data) == 38 {
		return ParseFields(data)
	}
	
	// Se tem 40 bytes, é um payload de dados criptografados real - usar apenas 38 bytes
	if len(data) == 40 {
		return ParseFields(data[:38])
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

// FormatDate converte a data do formato MDDYY para formato legível
// BB124 = B(mês base 36) + B(dia base 36) + 1(ano 1) + 24(ano 24) = ?
func (f Fields) FormatDate() string {
	if len(f.Date) != 5 {
		return f.Date + " (formato inválido)"
	}
	
	// Para BB124: preciso entender o formato
	// B em base 36 = 11
	// B em base 36 = 11  
	// 124 poderia ser 1/24 (janeiro 2024) ou 12/4 (dezembro 2004)?
	
	// Vou assumir que BB representa mês e dia em base 36
	// E 124 representa ano (2024)
	month := string(f.Date[0])  // B
	day := string(f.Date[1])    // B
	year := f.Date[2:5]         // 124
	
	// Conversão base 36 para decimal
	monthNames := map[string]string{
		"1": "Janeiro", "2": "Fevereiro", "3": "Março", "4": "Abril",
		"5": "Maio", "6": "Junho", "7": "Julho", "8": "Agosto",
		"9": "Setembro", "A": "Outubro", "B": "Novembro", "C": "Dezembro",
	}
	
	monthName := monthNames[month]
	if monthName == "" {
		monthName = "Mês " + month
	}
	
	// B em base 36 = 11
	dayNum := "11"
	if day >= "0" && day <= "9" {
		dayNum = day
	} else if day >= "A" && day <= "Z" {
		// Conversão base 36: A=10, B=11, C=12, etc.
		dayNum = fmt.Sprintf("%d", int(day[0])-int('A')+10)
	}
	
	// Interpretar ano 124 como 2024
	yearFormatted := "20" + year[1:]
	
	return fmt.Sprintf("%s de %s de %s", dayNum, monthName, yearFormatted)
}

// FormatColor converte a cor para formato legível
func (f Fields) FormatColor() string {
	if len(f.Color) == 7 && f.Color[0] == '0' {
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
		"00001": "PLA",
		"00002": "PLA-Silk",
		"00003": "PETG",
		"00004": "ABS",
		"00005": "TPU",
		"00006": "PLA-CF",
		"00007": "ASA",
		"00008": "PA",
		"00009": "PA-CF",
		"00010": "BVOH",
		"00012": "HIPS",
		"00013": "PET-CF",
		"00014": "PETG-CF",
		"00015": "PA6-CF",
		"00016": "PAHT-CF",
		"00020": "PET",
		"00021": "PC",
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
		"10100": "CR-PLA Especial", // Código antigo
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

// IsBlankTag verifica se a tag parece estar virgem ou com dados inválidos
func (f Fields) IsBlankTag() bool {
	// Verificar se os campos críticos têm valores válidos esperados
	
	// Batch deve ser "A2" em tags válidas
	if f.Batch != "A2" {
		return true
	}
	
	// Supplier deve ter um código conhecido
	knownSuppliers := []string{"0276", "0000"}
	validSupplier := false
	for _, supplier := range knownSuppliers {
		if f.Supplier == supplier {
			validSupplier = true
			break
		}
	}
	if !validSupplier {
		return true
	}
	
	// Material deve ter formato de código numérico
	if len(f.Material) != 5 {
		return true
	}
	
	// Reserve deve ser "0000" em tags válidas
	if f.Reserve != "0000" {
		return true
	}
	
	// Color deve ter formato válido (7 chars, começando com 0)
	if len(f.Color) != 7 || f.Color[0] != '0' {
		return true
	}
	
	return false
}

// TODO: EncryptPayloadToBlocks(payload string) (b4,b5,b6 string, error)
// (uso dos mesmos segredos AES/ECB do site)
