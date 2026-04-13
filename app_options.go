package main

// OptionsResponse resposta com opções para os dropdowns
type OptionsResponse struct {
	Materials []MaterialOption `json:"materials"`
	Vendors   []VendorOption   `json:"vendors"`
	Lengths   []LengthOption   `json:"lengths"`
}

// MaterialOption opção de material para dropdown
type MaterialOption struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Vendor string `json:"vendor"` // código do fornecedor na UI (filtragem)
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
	// Materiais Genéricos (códigos 00xxx) — Vendor "0000"
	{"00001", "PLA", "0000"},
	{"00002", "PLA-Silk", "0000"},
	{"00003", "PETG", "0000"},
	{"00004", "ABS", "0000"},
	{"00005", "TPU", "0000"},
	{"00006", "PLA-CF", "0000"},
	{"00007", "ASA", "0000"},
	{"00008", "PA", "0000"},
	{"00009", "PA-CF", "0000"},
	{"00010", "BVOH", "0000"},
	{"00011", "PVA", "0000"},
	{"00012", "HIPS", "0000"},
	{"00013", "PET-CF", "0000"},
	{"00014", "PETG-CF", "0000"},
	{"00015", "PA6-CF", "0000"},
	{"00016", "PAHT-CF", "0000"},
	{"00017", "PPS", "0000"},
	{"00018", "PPS-CF", "0000"},
	{"00019", "PP", "0000"},
	{"00020", "PET", "0000"},
	{"00021", "PC", "0000"},
	{"00022", "PA612-CF", "0000"},
	{"00023", "Support for PA", "0000"},
	{"00024", "Support for PLA", "0000"},
	{"00025", "PA12-CF", "0000"},
	{"00026", "TPU 64D", "0000"},
	{"00027", "PETG-GF", "0000"},
	{"00031", "PP-CF", "0000"},
	{"00032", "PCTG", "0000"},
	{"00033", "ASA-CF", "0000"},
	{"00034", "PA6-GF", "0000"},

	// Materiais Creality (códigos 01xxx-29xxx) — Vendor "0276"
	{"01001", "Hyper PLA", "0276"},
	{"01002", "Hyper L-W PLA", "0276"},
	{"01004", "Hyper Stardust", "0276"},
	{"01601", "Soleyin Ultra PLA", "0276"},
	{"02001", "Hyper PLA-CF", "0276"},
	{"03001", "Hyper ABS", "0276"},
	{"04001", "CR-PLA", "0276"},
	{"05001", "CR-Silk", "0276"},
	{"06001", "CR-PETG", "0276"},
	{"06002", "Hyper PETG", "0276"},
	{"06003", "Hyper PETG-CF", "0276"},
	{"06004", "Hyper PETG-GF", "0276"},
	{"07001", "CR-ABS", "0276"},
	{"07002", "Hyper PC", "0276"},
	{"08001", "Ender-PLA", "0276"},
	{"09001", "EN-PLA+", "0276"},
	{"09002", "ENDER FAST PLA", "0276"},
	{"10001", "HP-TPU", "0276"},
	{"11001", "CR-Nylon", "0276"},
	{"12002", "Hyper PPA-CF", "0276"},
	{"12003", "Hyper PAHT-CF", "0276"},
	{"12004", "Hyper PA612-CF", "0276"},
	{"12005", "Hyper PA6-CF", "0276"},
	{"13001", "CR-PLA Carbon", "0276"},
	{"14001", "CR-PLA Matte", "0276"},
	{"15001", "CR-PLA Fluo", "0276"},
	{"16001", "CR-TPU", "0276"},
	{"17001", "CR-Wood", "0276"},
	{"18001", "HP Ultra PLA", "0276"},
	{"19001", "HP-ASA", "0276"},
	{"29001", "Hyper Marble", "0276"},

	// Materiais eSUN — Vendor "ESUN"
	{"00035", "eSUN PLA-LW", "ESUN"},
	{"E1001", "eSUN PLA+", "ESUN"},
	{"E1002", "eSUN PLA-Silk", "ESUN"},
	{"E1003", "eSUN PLA-Matte", "ESUN"},
	{"E1004", "eSUN PLA-Lite", "ESUN"},
	{"E1005", "eSUN PLA-CF", "ESUN"},
	{"E1006", "eSUN PLA-HS", "ESUN"},
	{"E2001", "eSUN PETG", "ESUN"},
	{"E2002", "eSUN PETG+HS", "ESUN"},

	// Materiais Polymaker — Vendor "POLY"
	{"P1001", "Panchroma PLA Satin", "POLY"},
	{"P1002", "PolySonic PLA Pro", "POLY"},
	{"P1003", "Panchroma PLA Matte", "POLY"},
	{"P1004", "PolySonic PLA", "POLY"},
}

var vendors = []VendorOption{
	{"0276", "Creality"},
	{"0000", "Genérico"},
	{"ESUN", "eSUN"},
	{"POLY", "Polymaker"},
}

var lengths = []LengthOption{
	{"0083", "83cm (250g)", "250"},
	{"0165", "165cm (500g)", "500"},
	{"0330", "330cm (1kg)", "1000"},
	{"0660", "660cm (2kg)", "2000"},
	{"CUSTOM", "Personalizado", "0"},
}
