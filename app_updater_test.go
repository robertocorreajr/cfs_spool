package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/robertocorreajr/cfs_spool/internal/updater"
)

// stubReleasesServer monta um httptest server que devolve um JSON
// fixo de release. Útil para evitar bater na API real do GitHub em testes.
func stubReleasesServer(t *testing.T, payload string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/owner/repo/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(payload))
	})
	return httptest.NewServer(mux)
}

// withUpdaterStub substitui o checker e o config store da App por instâncias
// de teste apontando para um servidor httptest e arquivo temporário. Devolve
// uma função para restaurar os originais.
func withUpdaterStub(t *testing.T, app *App, baseURL, configPath string) func() {
	t.Helper()
	origChecker := app.updateChecker
	origStore := app.updateConfig

	app.updateChecker = &updater.Checker{
		Repo:    "owner/repo",
		BaseURL: baseURL,
	}
	app.updateConfig = &updater.ConfigStore{Path: configPath}

	return func() {
		app.updateChecker = origChecker
		app.updateConfig = origStore
	}
}

// TestCheckForUpdate_PopulaDownloadParaSO confirma que UpdateInfo carrega
// DownloadURL/DownloadName escolhidos pelo SO atual via PickAsset.
func TestCheckForUpdate_PopulaDownloadParaSO(t *testing.T) {
	original := version
	defer func() { version = original }()
	version = "v3.0.0"

	originalGOOS := platformGOOS
	defer func() { platformGOOS = originalGOOS }()
	platformGOOS = "darwin"

	srv := stubReleasesServer(t, `{
		"tag_name": "v3.1.0",
		"html_url": "https://example.com/r/v3.1.0",
		"published_at": "2026-05-01T12:00:00Z",
		"body": "x",
		"assets": [
			{"name": "cfs-spool-darwin-universal.dmg", "browser_download_url": "https://example.com/dmg", "size": 1},
			{"name": "cfs-spool-windows-amd64.zip", "browser_download_url": "https://example.com/win", "size": 1}
		]
	}`)
	defer srv.Close()

	app := NewApp()
	app.ctx = context.Background()
	restore := withUpdaterStub(t, app, srv.URL, filepath.Join(t.TempDir(), "cfg.json"))
	defer restore()

	info, err := app.CheckForUpdate()
	if err != nil {
		t.Fatalf("CheckForUpdate erro: %v", err)
	}
	if info.OS != "darwin" {
		t.Errorf("OS = %q, esperado %q", info.OS, "darwin")
	}
	if info.DownloadURL != "https://example.com/dmg" {
		t.Errorf("DownloadURL = %q, esperado URL do .dmg", info.DownloadURL)
	}
	if info.DownloadName != "cfs-spool-darwin-universal.dmg" {
		t.Errorf("DownloadName = %q", info.DownloadName)
	}
}

// TestCheckForUpdate_SemAssetParaSO garante que DownloadURL fica vazia
// quando nenhum asset bate com o SO atual — frontend cai no botão "Abrir
// release no GitHub" como fallback.
func TestCheckForUpdate_SemAssetParaSO(t *testing.T) {
	original := version
	defer func() { version = original }()
	version = "v3.0.0"

	originalGOOS := platformGOOS
	defer func() { platformGOOS = originalGOOS }()
	platformGOOS = "freebsd" // não há asset .freebsd

	srv := stubReleasesServer(t, `{
		"tag_name": "v3.1.0",
		"html_url": "https://example.com/r/v3.1.0",
		"published_at": "2026-05-01T12:00:00Z",
		"body": "x",
		"assets": [
			{"name": "cfs-spool-darwin-universal.dmg", "browser_download_url": "https://example.com/dmg", "size": 1}
		]
	}`)
	defer srv.Close()

	app := NewApp()
	app.ctx = context.Background()
	restore := withUpdaterStub(t, app, srv.URL, filepath.Join(t.TempDir(), "cfg.json"))
	defer restore()

	info, _ := app.CheckForUpdate()
	if info.DownloadURL != "" {
		t.Errorf("DownloadURL = %q, esperado vazio (sem asset freebsd)", info.DownloadURL)
	}
}

// TestCheckForUpdate_NovaVersao simula uma release nova e verifica que
// a binding devolve UpdateInfo.Available = true com URL e changelog.
func TestCheckForUpdate_NovaVersao(t *testing.T) {
	original := version
	defer func() { version = original }()
	version = "v3.0.0"

	srv := stubReleasesServer(t, `{
		"tag_name": "v3.1.0",
		"name": "v3.1.0",
		"html_url": "https://example.com/r/v3.1.0",
		"published_at": "2026-05-01T12:00:00Z",
		"body": "novidades"
	}`)
	defer srv.Close()

	app := NewApp()
	app.ctx = context.Background()
	restore := withUpdaterStub(t, app, srv.URL, filepath.Join(t.TempDir(), "cfg.json"))
	defer restore()

	info, err := app.CheckForUpdate()
	if err != nil {
		t.Fatalf("CheckForUpdate erro: %v", err)
	}
	if !info.Available {
		t.Error("Available = false, esperado true (v3.1.0 > v3.0.0)")
	}
	if info.Version != "v3.1.0" {
		t.Errorf("Version = %q, esperado %q", info.Version, "v3.1.0")
	}
	if info.URL == "" {
		t.Error("URL vazia")
	}
	if info.Ignored {
		t.Error("Ignored = true mas não foi ignorado ainda")
	}
}

// TestCheckForUpdate_VersaoAtualizada confirma que Available = false
// quando current >= latest — frontend vai exibir "você está atualizado".
func TestCheckForUpdate_VersaoAtualizada(t *testing.T) {
	original := version
	defer func() { version = original }()
	version = "v3.5.0"

	srv := stubReleasesServer(t, `{
		"tag_name": "v3.1.0",
		"name": "v3.1.0",
		"html_url": "https://example.com/r/v3.1.0",
		"published_at": "2026-05-01T12:00:00Z",
		"body": ""
	}`)
	defer srv.Close()

	app := NewApp()
	app.ctx = context.Background()
	restore := withUpdaterStub(t, app, srv.URL, filepath.Join(t.TempDir(), "cfg.json"))
	defer restore()

	info, err := app.CheckForUpdate()
	if err != nil {
		t.Fatalf("CheckForUpdate erro: %v", err)
	}
	if info.Available {
		t.Errorf("Available = true, esperado false (v3.5.0 >= v3.1.0)")
	}
}

// TestCheckForUpdate_VersaoIgnorada verifica que após IgnoreUpdateVersion,
// CheckForUpdate ainda devolve a release mas com Ignored = true para o
// frontend não exibir o toast.
func TestCheckForUpdate_VersaoIgnorada(t *testing.T) {
	original := version
	defer func() { version = original }()
	version = "v3.0.0"

	srv := stubReleasesServer(t, `{
		"tag_name": "v3.1.0",
		"name": "v3.1.0",
		"html_url": "https://example.com/r/v3.1.0",
		"published_at": "2026-05-01T12:00:00Z",
		"body": ""
	}`)
	defer srv.Close()

	app := NewApp()
	app.ctx = context.Background()
	restore := withUpdaterStub(t, app, srv.URL, filepath.Join(t.TempDir(), "cfg.json"))
	defer restore()

	if err := app.IgnoreUpdateVersion("v3.1.0"); err != nil {
		t.Fatalf("IgnoreUpdateVersion erro: %v", err)
	}

	info, err := app.CheckForUpdate()
	if err != nil {
		t.Fatalf("CheckForUpdate erro: %v", err)
	}
	if !info.Available {
		t.Errorf("Available = false (deveria continuar refletindo o estado real)")
	}
	if !info.Ignored {
		t.Error("Ignored = false, esperado true após IgnoreUpdateVersion")
	}
}

// TestCheckAndEmitUpdate_EmiteEvento_QuandoNovaVersao garante que o helper
// chamado no startup emite "update:available" com o payload correto, mas
// só quando há versão nova e não-ignorada.
func TestCheckAndEmitUpdate_EmiteEvento_QuandoNovaVersao(t *testing.T) {
	original := version
	defer func() { version = original }()
	version = "v3.0.0"

	srv := stubReleasesServer(t, `{
		"tag_name": "v3.1.0",
		"name": "v3.1.0",
		"html_url": "https://example.com/r/v3.1.0",
		"published_at": "2026-05-01T12:00:00Z",
		"body": "changelog"
	}`)
	defer srv.Close()

	captured, restoreSpy := withEventSpy(t)
	defer restoreSpy()

	app := NewApp()
	app.ctx = context.Background()
	restore := withUpdaterStub(t, app, srv.URL, filepath.Join(t.TempDir(), "cfg.json"))
	defer restore()

	app.checkAndEmitUpdate()

	if len(*captured) != 1 {
		t.Fatalf("esperado 1 evento, obtido %d: %+v", len(*captured), *captured)
	}
	ev := (*captured)[0]
	if ev.name != "update:available" {
		t.Errorf("evento = %q, esperado %q", ev.name, "update:available")
	}
}

// TestCheckAndEmitUpdate_NaoEmite_QuandoIgnorada cobre o caso "usuário
// dispensou esta versão" — o evento não dispara, então a UI fica silenciosa.
func TestCheckAndEmitUpdate_NaoEmite_QuandoIgnorada(t *testing.T) {
	original := version
	defer func() { version = original }()
	version = "v3.0.0"

	srv := stubReleasesServer(t, `{
		"tag_name": "v3.1.0",
		"name": "v3.1.0",
		"html_url": "https://example.com/r/v3.1.0",
		"published_at": "2026-05-01T12:00:00Z",
		"body": ""
	}`)
	defer srv.Close()

	captured, restoreSpy := withEventSpy(t)
	defer restoreSpy()

	app := NewApp()
	app.ctx = context.Background()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	restore := withUpdaterStub(t, app, srv.URL, cfgPath)
	defer restore()

	if err := app.IgnoreUpdateVersion("v3.1.0"); err != nil {
		t.Fatalf("IgnoreUpdateVersion erro: %v", err)
	}

	app.checkAndEmitUpdate()

	for _, ev := range *captured {
		if ev.name == "update:available" {
			t.Errorf("update:available emitido apesar da versão estar ignorada")
		}
	}
}

// TestCheckAndEmitUpdate_NaoEmite_QuandoCheckerFalha cobre o fallback
// silencioso em caso de erro de rede / 404 / rate limit.
func TestCheckAndEmitUpdate_NaoEmite_QuandoCheckerFalha(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden) // simula rate limit
	}))
	defer srv.Close()

	captured, restoreSpy := withEventSpy(t)
	defer restoreSpy()

	app := NewApp()
	app.ctx = context.Background()
	restore := withUpdaterStub(t, app, srv.URL, filepath.Join(t.TempDir(), "cfg.json"))
	defer restore()

	app.checkAndEmitUpdate()

	for _, ev := range *captured {
		if ev.name == "update:available" {
			t.Errorf("update:available emitido apesar do erro 403")
		}
	}
}
