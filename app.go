package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/robertocorreajr/cfs_spool/internal/creality"
	"github.com/robertocorreajr/cfs_spool/internal/rfid"
)

// App estrutura principal da aplicação Wails
type App struct {
	ctx context.Context
}

// NewApp cria uma nova instância da aplicação
func NewApp() *App {
	return &App{}
}

// startup é chamado quando a aplicação inicia
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetVersion retorna a versão da aplicação
func (a *App) GetVersion() string {
	return version
}

// --- Tipos para comunicação com o frontend ---

// TagData dados lidos de uma tag RFID
type TagData struct {
	UID          string `json:"uid"`
	Date         string `json:"date"`         // YYYY-MM-DD para input date
	DateDisplay  string `json:"dateDisplay"`   // formato legível pt-BR
	SupplierCode string `json:"supplierCode"`  // "0276"
	SupplierName string `json:"supplierName"`  // "Creality"
	MaterialCode string `json:"materialCode"`  // "04001"
	MaterialName string `json:"materialName"`  // "CR-PLA"
	Color        string `json:"color"`         // "77BB41" (6 chars hex, sem prefixo)
	LengthCode   string `json:"lengthCode"`    // "0330"
	LengthDisplay string `json:"lengthDisplay"` // "330cm (1kg)"
	Serial       string `json:"serial"`        // "000001"
}

// WriteRequest dados enviados pelo frontend para gravação
type WriteRequest struct {
	Date     string `json:"date"`     // YYYY-MM-DD
	Supplier string `json:"supplier"` // código 4 chars
	Material string `json:"material"` // código 5 chars
	Color    string `json:"color"`    // 6 chars hex (sem # ou prefixo 0)
	Length   string `json:"length"`   // código 4 chars ou gramas
	Serial   string `json:"serial"`   // até 6 dígitos
}

// --- Métodos expostos via Wails bindings ---

// ReadTag lê uma tag RFID e retorna os dados decodificados
func (a *App) ReadTag() (*TagData, error) {
	// Abrir leitor RFID
	reader, err := rfid.Open()
	if err != nil {
		return nil, fmt.Errorf("Erro ao conectar leitor: %v", err)
	}
	defer reader.Close()

	// Obter UID
	uid, err := reader.UID()
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler UID: %v", err)
	}

	// Ler blocos 4, 5, 6
	var blocks []string
	for block := byte(4); block <= 6; block++ {
		// Tentar chave padrão primeiro (tags novas)
		data, err := reader.TryReadBlock(block, rfid.KeyTypeA, "FFFFFFFFFFFF")
		if err != nil {
			// Se falhar, tentar chave derivada do UID (tags usadas)
			derivedKey := reader.DeriveKeyFromUID(uid)
			data, err = reader.TryReadBlock(block, rfid.KeyTypeA, derivedKey)
			if err != nil {
				return nil, fmt.Errorf("Erro ao ler bloco %d: %v", block, err)
			}
		}
		blocks = append(blocks, data)
	}

	// Descriptografar dados
	decrypted, err := creality.DecryptBlocks(strings.Join(blocks, ""))
	if err != nil {
		return nil, fmt.Errorf("Erro na descriptografia: %v", err)
	}

	// Parsear campos
	fields, err := creality.ParseFieldsCompat(decrypted)
	if err != nil {
		return nil, fmt.Errorf("Erro ao parsear dados: %v", err)
	}

	// Verificar tag virgem
	if fields.IsBlankTag() {
		return nil, fmt.Errorf("Tag virgem detectada. Esta tag não contém dados válidos ou nunca foi gravada.")
	}

	// Extrair cor sem prefixo "0"
	color := ""
	if len(fields.Color) == 7 && fields.Color[0] == '0' {
		color = strings.ToUpper(fields.Color[1:])
	}

	return &TagData{
		UID:          uid,
		Date:         parseDateToISO(fields.Date),
		DateDisplay:  fields.FormatDate(),
		SupplierCode: fields.Supplier,
		SupplierName: fields.GetSupplierName(),
		MaterialCode: fields.Material,
		MaterialName: fields.GetMaterialName(),
		Color:        color,
		LengthCode:   fields.Length,
		LengthDisplay: fields.FormatLength(),
		Serial:       fields.Serial,
	}, nil
}

// WriteTag grava dados em uma tag RFID
func (a *App) WriteTag(req WriteRequest) error {
	// Validar cor
	validatedColor, err := a.ValidateColor(req.Color)
	if err != nil {
		return fmt.Errorf("Cor inválida: %v", err)
	}

	// Converter data YYYY-MM-DD para YYMDD
	date, err := convertDate(req.Date)
	if err != nil {
		return fmt.Errorf("Data inválida: %v", err)
	}

	// Converter material se necessário
	materialCode := convertMaterial(req.Material)

	// Converter comprimento se necessário
	lengthCode := convertLength(req.Length)

	// Abrir leitor RFID
	reader, err := rfid.Open()
	if err != nil {
		return fmt.Errorf("Erro ao conectar leitor: %v", err)
	}
	defer reader.Close()

	// Obter UID
	uid, err := reader.UID()
	if err != nil {
		return fmt.Errorf("Erro ao ler UID: %v", err)
	}

	// Preparar campos
	fields := creality.NewFields()
	fields.Date = date
	fields.Supplier = req.Supplier
	fields.Material = materialCode
	fields.Length = lengthCode
	fields.Serial = padSerial(req.Serial)

	// Definir cor com validação
	if err := fields.SetColor(validatedColor); err != nil {
		return fmt.Errorf("Erro no formato da cor: %v", err)
	}

	// Gerar payload de 48 bytes
	payload, err := fields.ASCIIConcat48()
	if err != nil {
		return fmt.Errorf("Erro na validação: %v", err)
	}

	// Criptografar dados
	b4, b5, b6, err := creality.EncryptPayloadToBlocks(payload)
	if err != nil {
		return fmt.Errorf("Erro na criptografia: %v", err)
	}

	// Escrever na tag
	blocksToWrite := []string{b4, b5, b6}
	err = reader.WriteTagCFS(uid, blocksToWrite, false)
	if err != nil {
		return fmt.Errorf("Erro na escrita: %v", err)
	}

	return nil
}

// GetOptions retorna as opções para os dropdowns do formulário
func (a *App) GetOptions() OptionsResponse {
	return OptionsResponse{
		Materials: materials,
		Vendors:   vendors,
		Lengths:   lengths,
	}
}

// ValidateColor valida uma string hex de 6 caracteres e retorna uppercase
func (a *App) ValidateColor(hex string) (string, error) {
	hex = strings.TrimSpace(hex)
	hex = strings.TrimPrefix(hex, "#")

	validHex := regexp.MustCompile(`^[0-9A-Fa-f]{6}$`)
	if !validHex.MatchString(hex) {
		return "", fmt.Errorf("cor deve ter exatamente 6 caracteres hexadecimais válidos (0-9, A-F)")
	}

	return strings.ToUpper(hex), nil
}

// --- Helpers privados ---

// convertDate converte data de YYYY-MM-DD para formato interno YYMDD (5 chars)
// Formato: YY (2 dígitos) + M (1 char: 1-9 para Jan-Set, A=Out, B=Nov, C=Dez) + DD (2 dígitos)
func convertDate(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	year := t.Year() % 100
	month := int(t.Month())
	day := t.Day()

	// Mês como single char: 1-9 para Jan-Set, A/B/C para Out/Nov/Dez
	var monthChar string
	if month <= 9 {
		monthChar = fmt.Sprintf("%d", month)
	} else {
		monthChar = string(rune('A' + month - 10))
	}

	return fmt.Sprintf("%02d%s%02d", year, monthChar, day), nil
}

// parseDateToISO converte data do formato interno YYMDD para YYYY-MM-DD
// Formato interno: char[0-1]=ano (YY), char[2]=mês (1-9 para Jan-Set, A=Out, B=Nov, C=Dez), char[3-4]=dia (DD)
func parseDateToISO(date5 string) string {
	if len(date5) != 5 {
		return ""
	}

	// Extrair ano (primeiros 2 chars)
	yearStr := date5[0:2]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return ""
	}
	year += 2000

	// Extrair mês (char 2)
	monthChar := date5[2]
	var month int
	if monthChar >= '1' && monthChar <= '9' {
		month = int(monthChar - '0')
	} else if monthChar == 'A' || monthChar == 'a' {
		month = 10
	} else if monthChar == 'B' || monthChar == 'b' {
		month = 11
	} else if monthChar == 'C' || monthChar == 'c' {
		month = 12
	} else {
		return ""
	}

	// Extrair dia (chars 3-4)
	dayStr := date5[3:5]
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// convertMaterial converte nome do material para código
func convertMaterial(material string) string {
	if len(material) == 5 {
		return material
	}

	materialMap := map[string]string{
		"PLA": "00001", "PLA-Silk": "00002", "PETG": "00003",
		"ABS": "00004", "TPU": "00005", "PLA-CF": "00006",
		"ASA": "00007", "PA": "00008", "PA-CF": "00009",
		"BVOH": "00010", "HIPS": "00012", "PET-CF": "00013",
		"PETG-CF": "00014", "PA6-CF": "00015", "PAHT-CF": "00016",
		"PET": "00020", "PC": "00021",
		"Hyper PLA": "01001", "Hyper PLA-CF": "02001", "Hyper ABS": "03001",
		"CR-PLA": "04001", "CR-Silk": "05001", "CR-PETG": "06001",
		"CR-ABS": "07001", "Ender-PLA": "08001", "EN-PLA+": "09001",
		"ENDERFASTPLA": "09002", "HP-TPU": "10001", "CR-PLA Especial": "10100",
		"CR-Nylon": "11001", "CR-PLACarbon": "13001", "CR-PLAMatte": "14001",
		"CR-PLAFluo": "15001", "CR-TPU": "16001", "CR-Wood": "17001",
		"HPUltraPLA": "18001", "HP-ASA": "19001",
	}

	if code, ok := materialMap[material]; ok {
		return code
	}
	return material
}

// convertLength converte comprimento para código hex
func convertLength(length string) string {
	if len(length) == 4 {
		return length
	}

	lengthMap := map[string]string{
		"0083": "0053", "0165": "00A5", "0330": "014A", "0660": "0294",
	}
	if code, ok := lengthMap[length]; ok {
		return code
	}

	gramMap := map[string]string{
		"250": "0053", "500": "00A5", "1000": "014A", "2000": "0294",
	}
	if code, ok := gramMap[length]; ok {
		return code
	}

	if grams, err := strconv.Atoi(length); err == nil {
		cm := grams / 3
		if cm > 65535 {
			cm = 65535
		}
		return fmt.Sprintf("%04X", cm)
	}

	return "0053"
}

// padSerial preenche o serial com zeros à esquerda até 6 dígitos
func padSerial(serial string) string {
	serial = strings.TrimSpace(serial)
	if serial == "" {
		return "000001"
	}
	for len(serial) < 6 {
		serial = "0" + serial
	}
	if len(serial) > 6 {
		serial = serial[:6]
	}
	return serial
}
