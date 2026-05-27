// Ferramenta local-only; não invocada por CI ou pelo binário Wails.
//
// Uso:
//
//	go run ./cmd/check-materials -src <caminho para materialList.json> [-format text|json]
//
// Compara a lista de materiais embutida em app_options.go com o catálogo
// oficial do Creality Print v7 (materialList.json) e reporta:
//   - MISSING: código presente no JSON oficial mas ausente em app_options.go
//   - STALE: código presente em app_options.go mas ausente no JSON oficial
//   - MISMATCH: código presente nos dois mas com nome ou vendor divergente
//
// Exit codes: 0 = sem divergências; 1 = divergências encontradas; 2 = erro de IO/parse/flag.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// officialMaterial representa uma entrada do materialList.json do Creality Print v7.
// Campos extras (density, minTemp, maxTemp, meterialType, rank, createTime) são ignorados.
type officialMaterial struct {
	ID    string `json:"id"`
	Brand string `json:"brand"`
	Name  string `json:"name"`
}

// localMaterial representa uma entrada da slice materials em app_options.go.
type localMaterial struct {
	Code   string
	Name   string
	Vendor string
}

// diffItem representa uma divergência encontrada entre as duas fontes.
type diffItem struct {
	Code     string   `json:"code"`
	Expected *matInfo `json:"expected,omitempty"`
	Actual   *matInfo `json:"actual,omitempty"`
	Reason   string   `json:"reason,omitempty"`
}

// matInfo agrega nome e vendor para exibição nas divergências.
type matInfo struct {
	Name   string `json:"name"`
	Vendor string `json:"vendor"`
}

// Report agrega as três categorias de divergência.
type Report struct {
	Missing  []diffItem
	Stale    []diffItem
	Mismatch []diffItem
}

// brandToVendor mapeia o campo brand do JSON oficial para o código de vendor da UI.
var brandToVendor = map[string]string{
	"Creality":  "0276",
	"eSUN":      "ESUN",
	"Polymaker": "POLY",
	"Generic":   "0000",
}

func main() {
	srcFlag := flag.String("src", "", "caminho para materialList.json (obrigatório)")
	formatFlag := flag.String("format", "text", "formato de saída: text ou json")
	flag.Parse()

	if *srcFlag == "" {
		fmt.Fprintln(os.Stderr, "erro: flag -src é obrigatória")
		fmt.Fprintln(os.Stderr, "uso: go run ./cmd/check-materials -src <caminho> [-format text|json]")
		os.Exit(2)
	}

	if *formatFlag != "text" && *formatFlag != "json" {
		fmt.Fprintf(os.Stderr, "erro: -format deve ser \"text\" ou \"json\", recebido %q\n", *formatFlag)
		os.Exit(2)
	}

	// Resolve caminho de app_options.go relativo a este arquivo fonte
	_, thisFile, _, _ := runtime.Caller(0)
	appOptionsPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "app_options.go")

	official, err := loadOfficial(*srcFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao carregar %s: %v\n", *srcFlag, err)
		os.Exit(2)
	}

	local, err := loadLocal(appOptionsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao parsear %s: %v\n", appOptionsPath, err)
		os.Exit(2)
	}

	report := diff(official, local)

	switch *formatFlag {
	case "json":
		data, err := renderJSON(report)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro ao serializar JSON: %v\n", err)
			os.Exit(2)
		}
		fmt.Println(string(data))
	default:
		fmt.Print(renderText(report))
	}

	if len(report.Missing) > 0 || len(report.Stale) > 0 || len(report.Mismatch) > 0 {
		os.Exit(1)
	}
}

// officialList é o wrapper de nível raiz do materialList.json do Creality Print v7.
// O arquivo real tem shape {"materials": [...]} — não um array bare.
type officialList struct {
	Materials []officialMaterial `json:"materials"`
}

// loadOfficial lê e parseia o materialList.json do Creality Print v7.
func loadOfficial(path string) ([]officialMaterial, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("leitura falhou: %w", err)
	}
	var wrapped officialList
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return nil, fmt.Errorf("JSON inválido: %w", err)
	}
	return wrapped.Materials, nil
}

// loadLocal parseia app_options.go via go/ast e extrai a slice literal `materials`.
func loadLocal(goFile string) ([]localMaterial, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, goFile, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse falhou: %w", err)
	}

	var result []localMaterial
	var found bool

	// Percorre declarações de nível de pacote procurando `var materials = []MaterialOption{...}`
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			valSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range valSpec.Names {
				if name.Name != "materials" {
					continue
				}
				found = true
				if i >= len(valSpec.Values) {
					continue
				}
				compLit, ok := valSpec.Values[i].(*ast.CompositeLit)
				if !ok {
					continue
				}
				for _, elt := range compLit.Elts {
					innerLit, ok := elt.(*ast.CompositeLit)
					if !ok {
						continue
					}
					if len(innerLit.Elts) < 3 {
						continue
					}
					code := extractStringLit(innerLit.Elts[0])
					matName := extractStringLit(innerLit.Elts[1])
					vendor := extractStringLit(innerLit.Elts[2])
					if code == "" || matName == "" || vendor == "" {
						continue
					}
					result = append(result, localMaterial{Code: code, Name: matName, Vendor: vendor})
				}
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("variável \"materials\" não encontrada em %s", goFile)
	}
	return result, nil
}

// extractStringLit extrai o valor de um nó AST BasicLit do tipo STRING.
func extractStringLit(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	// Remove as aspas duplas ao redor
	return strings.Trim(lit.Value, `"`)
}

// diff compara o catálogo oficial com a lista local e produz um Report.
func diff(official []officialMaterial, local []localMaterial) Report {
	// Índice por código para acesso rápido
	offByCode := make(map[string]officialMaterial, len(official))
	for _, o := range official {
		offByCode[o.ID] = o
	}

	locByCode := make(map[string]localMaterial, len(local))
	for _, l := range local {
		locByCode[l.Code] = l
	}

	var report Report

	// MISSING: no JSON mas não no local
	for _, o := range official {
		if _, exists := locByCode[o.ID]; !exists {
			expectedVendor, known := brandToVendor[o.Brand]
			reason := ""
			if !known {
				expectedVendor = "?"
				reason = fmt.Sprintf("brand desconhecida: %q", o.Brand)
			}
			item := diffItem{
				Code: o.ID,
				Expected: &matInfo{
					Name:   o.Name,
					Vendor: expectedVendor,
				},
			}
			if reason != "" {
				item.Reason = reason
			}
			report.Missing = append(report.Missing, item)
		}
	}

	// STALE: no local mas não no JSON oficial
	for _, l := range local {
		if _, exists := offByCode[l.Code]; !exists {
			report.Stale = append(report.Stale, diffItem{
				Code:   l.Code,
				Actual: &matInfo{Name: l.Name, Vendor: l.Vendor},
			})
		}
	}

	// MISMATCH: no ambos mas com divergências
	for _, l := range local {
		o, exists := offByCode[l.Code]
		if !exists {
			continue
		}

		expectedVendor, known := brandToVendor[o.Brand]
		reason := ""
		if !known {
			expectedVendor = "?"
			reason = fmt.Sprintf("brand desconhecida: %q", o.Brand)
		}

		nameDiverge := o.Name != l.Name
		vendorDiverge := known && expectedVendor != l.Vendor

		if nameDiverge || vendorDiverge || !known {
			item := diffItem{
				Code:     l.Code,
				Expected: &matInfo{Name: o.Name, Vendor: expectedVendor},
				Actual:   &matInfo{Name: l.Name, Vendor: l.Vendor},
			}
			if reason != "" {
				item.Reason = reason
			}
			report.Mismatch = append(report.Mismatch, item)
		}
	}

	return report
}

// renderText produz saída legível por humano com seções [MISSING], [STALE], [MISMATCH].
func renderText(r Report) string {
	var sb strings.Builder

	if len(r.Missing) > 0 {
		sb.WriteString("[MISSING] — presentes no JSON oficial mas ausentes em app_options.go\n")
		sb.WriteString(strings.Repeat("-", 70) + "\n")
		for _, item := range r.Missing {
			sb.WriteString(fmt.Sprintf("  code=%-10s name=%-35s vendor=%s\n",
				item.Code, item.Expected.Name, item.Expected.Vendor))
			if item.Reason != "" {
				sb.WriteString(fmt.Sprintf("             aviso: %s\n", item.Reason))
			}
		}
		sb.WriteString("\n")
	}

	if len(r.Stale) > 0 {
		sb.WriteString("[STALE] — presentes em app_options.go mas ausentes no JSON oficial\n")
		sb.WriteString(strings.Repeat("-", 70) + "\n")
		for _, item := range r.Stale {
			sb.WriteString(fmt.Sprintf("  code=%-10s name=%-35s vendor=%s\n",
				item.Code, item.Actual.Name, item.Actual.Vendor))
		}
		sb.WriteString("\n")
	}

	if len(r.Mismatch) > 0 {
		sb.WriteString("[MISMATCH] — código existe nos dois mas nome ou vendor diverge\n")
		sb.WriteString(strings.Repeat("-", 70) + "\n")
		for _, item := range r.Mismatch {
			sb.WriteString(fmt.Sprintf("  code=%s\n", item.Code))
			sb.WriteString(fmt.Sprintf("    esperado: name=%-35s vendor=%s\n",
				item.Expected.Name, item.Expected.Vendor))
			sb.WriteString(fmt.Sprintf("    atual:    name=%-35s vendor=%s\n",
				item.Actual.Name, item.Actual.Vendor))
			if item.Reason != "" {
				sb.WriteString(fmt.Sprintf("    aviso: %s\n", item.Reason))
			}
		}
		sb.WriteString("\n")
	}

	total := len(r.Missing) + len(r.Stale) + len(r.Mismatch)
	sb.WriteString(fmt.Sprintf("Resumo: %d missing, %d stale, %d mismatch\n",
		len(r.Missing), len(r.Stale), len(r.Mismatch)))
	if total == 0 {
		sb.WriteString("Nenhuma divergência encontrada.\n")
	}

	return sb.String()
}

// renderJSON produz saída em JSON com chaves missing, stale, mismatch.
func renderJSON(r Report) ([]byte, error) {
	// Garante arrays vazios em vez de null no JSON
	missing := r.Missing
	if missing == nil {
		missing = []diffItem{}
	}
	stale := r.Stale
	if stale == nil {
		stale = []diffItem{}
	}
	mismatch := r.Mismatch
	if mismatch == nil {
		mismatch = []diffItem{}
	}

	out := struct {
		Missing  []diffItem `json:"missing"`
		Stale    []diffItem `json:"stale"`
		Mismatch []diffItem `json:"mismatch"`
	}{missing, stale, mismatch}

	return json.MarshalIndent(out, "", "  ")
}
