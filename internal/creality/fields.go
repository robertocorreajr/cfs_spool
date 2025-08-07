package creality

import (
	"errors"
	"fmt"
)

// tamanhos em bytes ASCII
const (
	lenBatch    = 3
	lenDate     = 5
	lenSupplier = 4
	lenMaterial = 5
	lenColor    = 7
	lenLength   = 4
	lenSerial   = 6
	lenReserve  = 14
)

type Fields struct {
	Batch, Date, Supplier, Material, Color, Length, Serial, Reserve string
}

func (f Fields) ASCIIConcat() (string, error) {
	if len(f.Batch) != lenBatch ||
		len(f.Date) != lenDate ||
		len(f.Supplier) != lenSupplier ||
		len(f.Material) != lenMaterial ||
		len(f.Color) != lenColor ||
		len(f.Length) != lenLength ||
		len(f.Serial) != lenSerial ||
		len(f.Reserve) != lenReserve {
		return "", errors.New("algum campo está com tamanho incorreto")
	}
	return fmt.Sprintf("%s%s%s%s%s%s%s%s",
		f.Batch, f.Date, f.Supplier, f.Material,
		f.Color, f.Length, f.Serial, f.Reserve), nil
}

// ParseFields extrai os campos de uma string ASCII de 48 bytes
func ParseFields(ascii48 string) (Fields, error) {
	if len(ascii48) != 48 {
		return Fields{}, errors.New("string ASCII deve ter exatamente 48 bytes")
	}

	// Extrair cada campo baseado no tamanho fixo
	fields := Fields{
		Batch:    ascii48[0:3],                    // 3 bytes
		Date:     ascii48[3:8],                    // 5 bytes  
		Supplier: ascii48[8:12],                   // 4 bytes
		Material: ascii48[12:17],                  // 5 bytes
		Color:    ascii48[17:24],                  // 7 bytes
		Length:   ascii48[24:28],                  // 4 bytes
		Serial:   ascii48[28:34],                  // 6 bytes
		Reserve:  ascii48[34:48],                  // 14 bytes
	}

	return fields, nil
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
	default:
		return f.Length + "cm"
	}
}

// GetMaterialName retorna o nome do material baseado no código
func (f Fields) GetMaterialName() string {
	materials := map[string]string{
		"00003": "PLA Genérico",
		"00004": "ABS Genérico", 
		"00007": "ASA Genérico",
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
		"01001": "CR-PLA (padrão)",
	}
	
	name := materials[f.Material]
	if name == "" {
		return f.Material + " (desconhecido)"
	}
	return name
}

// TODO: EncryptPayloadToBlocks(payload string) (b4,b5,b6 string, error)
// (uso dos mesmos segredos AES/ECB do site)
