package updater

import (
	"os"
	"path/filepath"
	"testing"
)

// TestConfigStore_RoundTrip cria um arquivo temporário, persiste uma versão
// ignorada e confirma que IsIgnored devolve true depois da releitura.
func TestConfigStore_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	store := &ConfigStore{Path: path}

	if store.IsIgnored("v3.1.0") {
		t.Error("IsIgnored = true antes de IgnoreVersion (deveria ser false)")
	}

	if err := store.IgnoreVersion("v3.1.0"); err != nil {
		t.Fatalf("IgnoreVersion erro: %v", err)
	}

	// Nova instância lendo o mesmo arquivo simula reabertura do app.
	store2 := &ConfigStore{Path: path}
	if !store2.IsIgnored("v3.1.0") {
		t.Error("IsIgnored = false após IgnoreVersion + reload")
	}
	if store2.IsIgnored("v3.2.0") {
		t.Error("IsIgnored = true para versão diferente — deveria ser false")
	}
}

// TestConfigStore_ArquivoInexistente garante que IsIgnored não panica
// nem retorna erro quando o arquivo de config ainda não existe (primeiro uso).
func TestConfigStore_ArquivoInexistente(t *testing.T) {
	path := filepath.Join(t.TempDir(), "subdir", "naoexiste.json")
	store := &ConfigStore{Path: path}

	if store.IsIgnored("v3.0.0") {
		t.Error("IsIgnored deveria retornar false em arquivo inexistente")
	}
}

// TestConfigStore_ArquivoCorrompido cobre o caso de JSON inválido — o store
// trata como vazio em vez de falhar (o usuário não deve ficar travado por
// causa de um arquivo de config danificado).
func TestConfigStore_ArquivoCorrompido(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := writeFile(t, path, []byte("isso não é json válido")); err != nil {
		t.Fatal(err)
	}

	store := &ConfigStore{Path: path}
	if store.IsIgnored("v3.0.0") {
		t.Error("IsIgnored deveria devolver false em config corrompido")
	}
}

// TestConfigStore_IgnoreDuasVezes verifica idempotência — chamar IgnoreVersion
// para a mesma versão duas vezes não deve duplicar entradas nem dar erro.
func TestConfigStore_IgnoreDuasVezes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	store := &ConfigStore{Path: path}

	for range 2 {
		if err := store.IgnoreVersion("v3.1.0"); err != nil {
			t.Fatalf("IgnoreVersion erro: %v", err)
		}
	}

	store2 := &ConfigStore{Path: path}
	if !store2.IsIgnored("v3.1.0") {
		t.Error("IsIgnored = false após dupla escrita")
	}
}

// writeFile delega para os.WriteFile com permissão padrão.
func writeFile(t *testing.T, path string, data []byte) error {
	t.Helper()
	return os.WriteFile(path, data, 0o644)
}
