package main

// OptionsResponse resposta com opções para os dropdowns
type OptionsResponse struct {
	Materials []MaterialOption `json:"materials"`
	Vendors   []VendorOption   `json:"vendors"`
	Lengths   []LengthOption   `json:"lengths"`
}

// MaterialOption opção de material para dropdown
type MaterialOption struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// VendorOption opção de fornecedor para dropdown
type VendorOption struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// LengthOption opção de comprimento para dropdown
type LengthOption struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Grams string `json:"grams"`
}

// Dados estáticos para os dropdowns

var materials = []MaterialOption{
	// Materiais Genéricos (códigos 00xxx)
	{"00001", "PLA"},
	{"00002", "PLA-Silk"},
	{"00003", "PETG"},
	{"00004", "ABS"},
	{"00005", "TPU"},
	{"00006", "PLA-CF"},
	{"00007", "ASA"},
	{"00008", "PA"},
	{"00009", "PA-CF"},
	{"00010", "BVOH"},
	{"00012", "HIPS"},
	{"00013", "PET-CF"},
	{"00014", "PETG-CF"},
	{"00015", "PA6-CF"},
	{"00016", "PAHT-CF"},
	{"00020", "PET"},
	{"00021", "PC"},

	// Materiais Hyper (códigos 01xxx-03xxx)
	{"01001", "Hyper PLA"},
	{"02001", "Hyper PLA-CF"},
	{"03001", "Hyper ABS"},

	// Materiais Creality (códigos 04xxx+)
	{"04001", "CR-PLA"},
	{"05001", "CR-Silk"},
	{"06001", "CR-PETG"},
	{"07001", "CR-ABS"},
	{"08001", "Ender-PLA"},
	{"09001", "EN-PLA+"},
	{"09002", "ENDERFASTPLA"},
	{"10001", "HP-TPU"},
	{"10100", "CR-PLA Especial"},
	{"11001", "CR-Nylon"},
	{"13001", "CR-PLACarbon"},
	{"14001", "CR-PLAMatte"},
	{"15001", "CR-PLAFluo"},
	{"16001", "CR-TPU"},
	{"17001", "CR-Wood"},
	{"18001", "HPUltraPLA"},
	{"19001", "HP-ASA"},
}

var vendors = []VendorOption{
	{"0276", "Creality"},
	{"0000", "Genérico"},
}

var lengths = []LengthOption{
	{"0083", "83cm (250g)", "250"},
	{"0165", "165cm (500g)", "500"},
	{"0330", "330cm (1kg)", "1000"},
	{"0660", "660cm (2kg)", "2000"},
	{"CUSTOM", "Personalizado", "0"},
}
