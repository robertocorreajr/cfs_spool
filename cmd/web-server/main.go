package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/robertocorreajr/cfs_spool/internal/creality"
	"github.com/robertocorreajr/cfs_spool/internal/rfid"
)

// Estruturas para API
type ReadResponse struct {
	Success bool              `json:"success"`
	Data    *TagInfo          `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}

type WriteRequest struct {
	Date     string `json:"date"`     // Formato: YYYY-MM-DD
	Supplier string `json:"supplier"`
	Material string `json:"material"` // Código ou nome
	Color    string `json:"color"`
	Length   string `json:"length"`   // Código ou valor em gramas
	Serial   string `json:"serial"`
}

type WriteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type TagInfo struct {
	UID      string `json:"uid"`
	Date     string `json:"date"`
	Supplier string `json:"supplier"`
	Material string `json:"material"`
	Color    string `json:"color"`
	Length   string `json:"length"`
	Serial   string `json:"serial"`
}

type OptionsResponse struct {
	Materials []MaterialOption `json:"materials"`
	Vendors   []VendorOption   `json:"vendors"`
	Lengths   []LengthOption   `json:"lengths"`
}

type MaterialOption struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type VendorOption struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type LengthOption struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Grams  string `json:"grams"`
}

// CORS middleware
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Handler para obter opções de dropdowns
func getOptionsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	materials := []MaterialOption{
		{"00001", "Generic PLA"},
		{"00002", "Generic PLA-Silk"},
		{"00003", "Generic PETG"},
		{"00004", "Generic ABS"},
		{"00005", "Generic TPU"},
		{"00006", "Generic PLA-CF"},
		{"00007", "Generic ASA"},
		{"04001", "CR-PLA"},
		{"05001", "CR-Silk"},
		{"06001", "CR-PETG"},
		{"07001", "CR-ABS"},
		{"08001", "Ender-PLA"},
		{"09001", "EN-PLA+"},
		{"10001", "HP-TPU"},
		{"11001", "CR-Nylon"},
	}

	vendors := []VendorOption{
		{"1B3D", "Creality"},
		{"FFFF", "Genérico"},
	}

	lengths := []LengthOption{
		{"0083", "83cm (250g)", "250"},
		{"0165", "165cm (500g)", "500"},
		{"0330", "330cm (1kg)", "1000"},
		{"0660", "660cm (2kg)", "2000"},
		{"CUSTOM", "Personalizado", "0"},
	}

	response := OptionsResponse{
		Materials: materials,
		Vendors:   vendors,
		Lengths:   lengths,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler para ler tag
func readTagHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Abrir leitor RFID
	reader, err := rfid.Open()
	if err != nil {
		response := ReadResponse{
			Success: false,
			Error:   fmt.Sprintf("Erro ao conectar leitor: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	defer reader.Close()

	// Obter UID
	uid, err := reader.UID()
	if err != nil {
		response := ReadResponse{
			Success: false,
			Error:   fmt.Sprintf("Erro ao ler UID: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Ler dados da tag usando TryReadBlock para garantir autenticação correta
	var blocks []string
	
	// Tentar ler blocos 4, 5, 6 individualmente
	for block := byte(4); block <= 6; block++ {
		var data string
		var err error
		
		// Primeiro tentar key padrão para tags novas
		data, err = reader.TryReadBlock(block, rfid.KeyTypeA, "FFFFFFFFFFFF")
		if err != nil {
			// Se falhar, tentar key derivada do UID para tags usadas
			derivedKey := reader.DeriveKeyFromUID(uid)
			data, err = reader.TryReadBlock(block, rfid.KeyTypeA, derivedKey)
			if err != nil {
				response := ReadResponse{
					Success: false,
					Error:   fmt.Sprintf("Erro ao ler bloco %d: %v", block, err),
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
		}
		blocks = append(blocks, data)
	}

	// Descriptografar dados
	decrypted, err := creality.DecryptBlocks(strings.Join(blocks, ""))
	if err != nil {
		response := ReadResponse{
			Success: false,
			Error:   fmt.Sprintf("Erro na descriptografia: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parsear campos
	fields, err := creality.ParseFields(decrypted)
	if err != nil {
		response := ReadResponse{
			Success: false,
			Error:   fmt.Sprintf("Erro ao parsear dados: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Preparar resposta
	tagInfo := &TagInfo{
		UID:      uid,
		Date:     fields.FormatDate(),
		Supplier: fields.GetSupplierName(),
		Material: fields.GetMaterialName(),
		Color:    fields.FormatColor(),
		Length:   fields.FormatLength(),
		Serial:   fields.Serial,
	}

	response := ReadResponse{
		Success: true,
		Data:    tagInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler para escrever tag
func writeTagHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := WriteResponse{
			Success: false,
			Error:   "Dados inválidos",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Converter data YYYY-MM-DD para YYMDD
	date, err := convertDate(req.Date)
	if err != nil {
		response := WriteResponse{
			Success: false,
			Error:   fmt.Sprintf("Data inválida: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Converter material se necessário
	materialCode := convertMaterial(req.Material)
	
	// Converter comprimento se necessário
	lengthCode := convertLength(req.Length)

	// Abrir leitor RFID
	reader, err := rfid.Open()
	if err != nil {
		response := WriteResponse{
			Success: false,
			Error:   fmt.Sprintf("Erro ao conectar leitor: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	defer reader.Close()

	// Obter UID
	uid, err := reader.UID()
	if err != nil {
		response := WriteResponse{
			Success: false,
			Error:   fmt.Sprintf("Erro ao ler UID: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Preparar campos usando a função construtora
	fields := creality.NewFields()
	fields.Date = date
	fields.Supplier = req.Supplier
	fields.Material = materialCode
	fields.Length = lengthCode
	fields.Serial = req.Serial
	
	// Definir cor garantindo formato correto (0 + 5 caracteres hex)
	if req.Color != "" {
		err := fields.SetColor(req.Color)
		if err != nil {
			response := WriteResponse{
				Success: false,
				Error:   fmt.Sprintf("Erro no formato da cor: %v", err),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Validar e criar payload (valores fixos são aplicados automaticamente)
	payload, err := fields.ASCIIConcat()
	if err != nil {
		response := WriteResponse{
			Success: false,
			Error:   fmt.Sprintf("Erro na validação: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Criptografar dados
	b4, b5, b6, err := creality.EncryptPayloadToBlocks(payload)
	if err != nil {
		response := WriteResponse{
			Success: false,
			Error:   fmt.Sprintf("Erro na criptografia: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Escrever na tag
	blocksToWrite := []string{b4, b5, b6}
	err = reader.WriteTagCFS(uid, blocksToWrite, false)
	if err != nil {
		response := WriteResponse{
			Success: false,
			Error:   fmt.Sprintf("Erro na escrita: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := WriteResponse{
		Success: true,
		Message: "Tag gravada com sucesso!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Funções auxiliares
func convertDate(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}
	
	year := t.Year() % 100 // Últimos 2 dígitos do ano
	month := int(t.Month())
	day := t.Day()
	
	return fmt.Sprintf("%02d%01d%02d", year, month, day), nil
}

func convertMaterial(material string) string {
	// Se já é um código, retorna
	if len(material) == 5 {
		return material
	}
	
	// Mapear nomes para códigos
	materials := map[string]string{
		"Generic PLA": "00001",
		"Generic ABS": "00004",
		"Generic PETG": "00003",
		"CR-PLA": "04001",
		"CR-ABS": "07001",
	}
	
	if code, ok := materials[material]; ok {
		return code
	}
	
	return material
}

func convertLength(length string) string {
	// Se já é um código hexadecimal de 4 dígitos, retorna
	if len(length) == 4 {
		return length
	}
	
	// Mapear códigos para hexadecimal
	lengths := map[string]string{
		"0083": "0053", // 250g -> 83cm
		"0165": "00A5", // 500g -> 165cm
		"0330": "014A", // 1kg -> 330cm
		"0660": "0294", // 2kg -> 660cm
	}
	
	if code, ok := lengths[length]; ok {
		return code
	}
	
	// Mapear gramas diretamente para hexadecimal
	gramToHex := map[string]string{
		"250":  "0053",
		"500":  "00A5", 
		"1000": "014A",
		"2000": "0294",
	}
	
	if code, ok := gramToHex[length]; ok {
		return code
	}
	
	// Tentar converter valor em gramas para centímetros em hexadecimal
	if grams, err := strconv.Atoi(length); err == nil {
		cm := grams / 3 // Aproximação: 3g = 1cm
		if cm > 65535 { // Limite do uint16
			cm = 65535
		}
		return fmt.Sprintf("%04X", cm)
	}
	
	return "0053" // Default para 250g
}

func main() {
	// Servir arquivos estáticos
	http.Handle("/", http.FileServer(http.Dir("./web/")))
	
	// API endpoints
	http.HandleFunc("/api/options", getOptionsHandler)
	http.HandleFunc("/api/read", readTagHandler)
	http.HandleFunc("/api/write", writeTagHandler)

	fmt.Println("🌐 Servidor iniciado em http://localhost:8080")
	fmt.Println("📱 Interface web disponível no navegador")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
