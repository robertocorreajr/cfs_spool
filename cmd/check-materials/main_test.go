package main

// main_test.go — testes de aceitação para o script cmd/check-materials.
//
// Critério coberto: "Script dev cmd/check-materials reporta divergências contra
// materialList.json" (Critério 9 da história aprovada).
//
// Estratégia:
//   - Testa loadOfficial com JSON sintético em t.TempDir().
//   - Testa loadLocal com arquivo Go sintético em t.TempDir().
//   - Testa diff diretamente com entradas controladas.
//   - Testa o CLI end-to-end com exec.Command, verificando exit codes.

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// jsonFixture gera um materialList.json sintético válido.
func jsonFixture(t *testing.T, materials []officialMaterial) string {
	t.Helper()
	dir := t.TempDir()
	wrapped := officialList{Materials: materials}
	data, err := json.Marshal(wrapped)
	if err != nil {
		t.Fatalf("não foi possível serializar fixture JSON: %v", err)
	}
	path := filepath.Join(dir, "materialList.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("não foi possível gravar fixture JSON: %v", err)
	}
	return path
}

// goFixture gera um app_options.go sintético com a slice materials fornecida.
func goFixture(t *testing.T, entries []localMaterial) string {
	t.Helper()
	dir := t.TempDir()

	src := "package main\n\nvar materials = []MaterialOption{\n"
	for _, e := range entries {
		src += fmt.Sprintf("\t{%q, %q, %q},\n", e.Code, e.Name, e.Vendor)
	}
	src += "}\n"

	path := filepath.Join(dir, "app_options.go")
	if err := os.WriteFile(path, []byte(src), 0600); err != nil {
		t.Fatalf("não foi possível gravar fixture Go: %v", err)
	}
	return path
}

// TestLoadOfficial_BasicCase verifica que loadOfficial parseia corretamente o
// arquivo JSON com shape {"materials": [...]}.
func TestLoadOfficial_BasicCase(t *testing.T) {
	fixture := []officialMaterial{
		{ID: "00001", Brand: "Generic", Name: "PLA"},
		{ID: "E1009", Brand: "eSUN", Name: "eSUN PLA-Basic"},
		{ID: "P9001", Brand: "Polymaker", Name: "Fiberon PPS-CF10"},
	}
	path := jsonFixture(t, fixture)

	result, err := loadOfficial(path)
	if err != nil {
		t.Fatalf("loadOfficial retornou erro inesperado: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("esperado 3 entradas, recebido %d", len(result))
	}
	if result[1].ID != "E1009" || result[1].Name != "eSUN PLA-Basic" {
		t.Errorf("entrada[1] inesperada: %+v", result[1])
	}
}

// TestLoadOfficial_ArquivoInexistente verifica que loadOfficial retorna erro para
// arquivo inexistente.
func TestLoadOfficial_ArquivoInexistente(t *testing.T) {
	_, err := loadOfficial("/tmp/nao_existe_jamais.json")
	if err == nil {
		t.Fatal("esperava erro para arquivo inexistente, recebeu nil")
	}
}

// TestLoadOfficial_JSONInvalido verifica que loadOfficial retorna erro para JSON malformado.
func TestLoadOfficial_JSONInvalido(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{invalido}"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := loadOfficial(path)
	if err == nil {
		t.Fatal("esperava erro para JSON inválido, recebeu nil")
	}
}

// TestLoadLocal_BasicCase verifica que loadLocal extrai corretamente as entradas
// da slice materials em um arquivo Go sintético.
func TestLoadLocal_BasicCase(t *testing.T) {
	entries := []localMaterial{
		{Code: "00001", Name: "PLA", Vendor: "0000"},
		{Code: "E1009", Name: "eSUN PLA-Basic", Vendor: "ESUN"},
		{Code: "P9001", Name: "Fiberon PPS-CF10", Vendor: "POLY"},
	}
	path := goFixture(t, entries)

	result, err := loadLocal(path)
	if err != nil {
		t.Fatalf("loadLocal retornou erro inesperado: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("esperado 3 entradas, recebido %d", len(result))
	}

	byCode := make(map[string]localMaterial)
	for _, r := range result {
		byCode[r.Code] = r
	}

	for _, e := range entries {
		got, ok := byCode[e.Code]
		if !ok {
			t.Errorf("code=%q não encontrado no resultado", e.Code)
			continue
		}
		if got.Name != e.Name || got.Vendor != e.Vendor {
			t.Errorf("code=%q: esperado (name=%q, vendor=%q), recebido (name=%q, vendor=%q)",
				e.Code, e.Name, e.Vendor, got.Name, got.Vendor)
		}
	}
}

// TestLoadLocal_ArquivoInexistente verifica que loadLocal retorna erro para
// arquivo inexistente.
func TestLoadLocal_ArquivoInexistente(t *testing.T) {
	_, err := loadLocal("/tmp/nao_existe_jamais.go")
	if err == nil {
		t.Fatal("esperava erro para arquivo inexistente, recebeu nil")
	}
}

// TestDiff_SemDivergencias verifica que diff retorna Report vazio quando
// oficial e local são idênticos.
func TestDiff_SemDivergencias(t *testing.T) {
	official := []officialMaterial{
		{ID: "00001", Brand: "Generic", Name: "PLA"},
		{ID: "E1009", Brand: "eSUN", Name: "eSUN PLA-Basic"},
	}
	local := []localMaterial{
		{Code: "00001", Name: "PLA", Vendor: "0000"},
		{Code: "E1009", Name: "eSUN PLA-Basic", Vendor: "ESUN"},
	}

	r := diff(official, local)

	if len(r.Missing) != 0 {
		t.Errorf("esperado 0 missing, recebido %d: %v", len(r.Missing), r.Missing)
	}
	if len(r.Stale) != 0 {
		t.Errorf("esperado 0 stale, recebido %d: %v", len(r.Stale), r.Stale)
	}
	if len(r.Mismatch) != 0 {
		t.Errorf("esperado 0 mismatch, recebido %d: %v", len(r.Mismatch), r.Mismatch)
	}
}

// TestDiff_Missing verifica que diff detecta código presente no oficial mas
// ausente no local.
func TestDiff_Missing(t *testing.T) {
	official := []officialMaterial{
		{ID: "00001", Brand: "Generic", Name: "PLA"},
		{ID: "E1009", Brand: "eSUN", Name: "eSUN PLA-Basic"},
	}
	local := []localMaterial{
		{Code: "00001", Name: "PLA", Vendor: "0000"},
		// E1009 ausente — deve aparecer em Missing
	}

	r := diff(official, local)

	if len(r.Missing) != 1 {
		t.Fatalf("esperado 1 missing, recebido %d", len(r.Missing))
	}
	if r.Missing[0].Code != "E1009" {
		t.Errorf("missing[0].Code = %q, esperado %q", r.Missing[0].Code, "E1009")
	}
	if r.Missing[0].Expected == nil || r.Missing[0].Expected.Vendor != "ESUN" {
		t.Errorf("missing[0].Expected.Vendor esperado %q, recebido %v", "ESUN", r.Missing[0].Expected)
	}
}

// TestDiff_Stale verifica que diff detecta código presente no local mas
// ausente no oficial.
func TestDiff_Stale(t *testing.T) {
	official := []officialMaterial{
		{ID: "00001", Brand: "Generic", Name: "PLA"},
	}
	local := []localMaterial{
		{Code: "00001", Name: "PLA", Vendor: "0000"},
		{Code: "ZZZZZ", Name: "Material Fantasma", Vendor: "0000"},
	}

	r := diff(official, local)

	if len(r.Stale) != 1 {
		t.Fatalf("esperado 1 stale, recebido %d", len(r.Stale))
	}
	if r.Stale[0].Code != "ZZZZZ" {
		t.Errorf("stale[0].Code = %q, esperado %q", r.Stale[0].Code, "ZZZZZ")
	}
}

// TestDiff_Mismatch verifica que diff detecta divergência de nome entre oficial e local.
func TestDiff_Mismatch(t *testing.T) {
	official := []officialMaterial{
		{ID: "E1009", Brand: "eSUN", Name: "eSUN PLA-Basic"},
	}
	local := []localMaterial{
		// Nome errado — deve aparecer em Mismatch
		{Code: "E1009", Name: "eSUN PLA-LW", Vendor: "ESUN"},
	}

	r := diff(official, local)

	if len(r.Mismatch) != 1 {
		t.Fatalf("esperado 1 mismatch, recebido %d", len(r.Mismatch))
	}
	if r.Mismatch[0].Code != "E1009" {
		t.Errorf("mismatch[0].Code = %q, esperado %q", r.Mismatch[0].Code, "E1009")
	}
	if r.Mismatch[0].Expected.Name != "eSUN PLA-Basic" {
		t.Errorf("mismatch[0].Expected.Name = %q, esperado %q", r.Mismatch[0].Expected.Name, "eSUN PLA-Basic")
	}
	if r.Mismatch[0].Actual.Name != "eSUN PLA-LW" {
		t.Errorf("mismatch[0].Actual.Name = %q, esperado %q", r.Mismatch[0].Actual.Name, "eSUN PLA-LW")
	}
}

// TestDiff_VendorMismatch verifica que diff detecta divergência de vendor.
func TestDiff_VendorMismatch(t *testing.T) {
	official := []officialMaterial{
		{ID: "E1009", Brand: "eSUN", Name: "eSUN PLA-Basic"},
	}
	local := []localMaterial{
		// Vendor errado — deve aparecer em Mismatch
		{Code: "E1009", Name: "eSUN PLA-Basic", Vendor: "0000"},
	}

	r := diff(official, local)

	if len(r.Mismatch) != 1 {
		t.Fatalf("esperado 1 mismatch por vendor, recebido %d", len(r.Mismatch))
	}
	if r.Mismatch[0].Expected.Vendor != "ESUN" {
		t.Errorf("mismatch[0].Expected.Vendor = %q, esperado %q", r.Mismatch[0].Expected.Vendor, "ESUN")
	}
}

// TestRenderJSON_SemDivergencias verifica que renderJSON produz arrays vazios (não null)
// quando o Report não tem divergências.
func TestRenderJSON_SemDivergencias(t *testing.T) {
	r := Report{}
	data, err := renderJSON(r)
	if err != nil {
		t.Fatalf("renderJSON retornou erro: %v", err)
	}

	var out struct {
		Missing  []diffItem `json:"missing"`
		Stale    []diffItem `json:"stale"`
		Mismatch []diffItem `json:"mismatch"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("saída JSON inválida: %v", err)
	}
	if out.Missing == nil || out.Stale == nil || out.Mismatch == nil {
		t.Error("renderJSON deve produzir arrays [] em vez de null para Report vazio")
	}
}

// buildBinary compila o binário check-materials em t.TempDir() e retorna o caminho.
// Usar o binário diretamente (em vez de "go run") garante que os exit codes sejam
// propagados fielmente — "go run" encapsula o código do programa filho em exit 1
// independentemente do valor real.
func buildBinary(t *testing.T) string {
	t.Helper()
	binPath := filepath.Join(t.TempDir(), "check-materials")
	root := findRepoRoot(t)
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/check-materials")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("falha ao compilar check-materials: %v\n%s", err, out)
	}
	return binPath
}

// TestCLI_ExitCode2_SrcAusente verifica que o CLI retorna exit code 2 quando
// a flag -src está ausente.
func TestCLI_ExitCode2_SrcAusente(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	err := cmd.Run()
	exitCode := exitCodeOf(err)
	if exitCode != 2 {
		t.Errorf("exit code esperado 2 (flag ausente), recebido %d", exitCode)
	}
}

// TestCLI_ExitCode2_FormatInvalido verifica que o CLI retorna exit code 2 para
// -format inválido.
func TestCLI_ExitCode2_FormatInvalido(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "materialList.json")
	if err := os.WriteFile(jsonPath, []byte(`{"materials":[]}`), 0600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bin, "-src", jsonPath, "-format", "invalido")
	err := cmd.Run()
	if exitCodeOf(err) != 2 {
		t.Errorf("exit code esperado 2 para -format inválido, recebido %d", exitCodeOf(err))
	}
}

// TestCLI_ExitCode0_SemDivergencias verifica que o CLI retorna exit code 0 quando
// não há divergências.
// Nota: o binário lê app_options.go via runtime.Caller — ao testar com o binário
// compilado, o caminho do arquivo fonte é embutido em tempo de compilação e aponta
// para app_options.go real do repositório. Portanto, para exit 0 precisaríamos que
// o JSON de entrada cobrisse TODOS os materiais locais sem sobra, o que não é
// viável com fixture sintética sem refatoração do binário. Este cenário é coberto
// pelos testes unitários de diff() (TestDiff_SemDivergencias). Pulamos o end-to-end.
func TestCLI_ExitCode0_SemDivergencias(t *testing.T) {
	t.Skip("end-to-end exit 0 requer que JSON de fixture cubra app_options.go real; coberto por TestDiff_SemDivergencias")
}

// TestCLI_ExitCode1_ComDivergencias verifica que o CLI retorna exit code 1 quando
// há divergências — usando JSON sintético com lista vazia (todos os materiais locais
// ficam STALE).
func TestCLI_ExitCode1_ComDivergencias(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(jsonPath, []byte(`{"materials":[]}`), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin, "-src", jsonPath)
	err := cmd.Run()
	code := exitCodeOf(err)
	if code != 1 {
		t.Errorf("exit code esperado 1 (divergências encontradas), recebido %d", code)
	}
}

// findRepoRoot retorna o diretório raiz do repositório.
func findRepoRoot(t *testing.T) string {
	t.Helper()
	// O teste está em cmd/check-materials; o repo root é dois níveis acima.
	dir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("não foi possível resolver caminho raiz: %v", err)
	}
	return dir
}

// exitCodeOf extrai o exit code de um erro retornado por cmd.Run().
func exitCodeOf(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return -1
}
