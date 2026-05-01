package updater

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestIsNewer cobre comparação semver de versões com e sem prefixo "v".
// Usa golang.org/x/mod/semver internamente — strings de qualquer comprimento
// são comparadas por componente (major.minor.patch), nunca lexicalmente.
func TestIsNewer(t *testing.T) {
	testes := []struct {
		nome     string
		latest   string
		current  string
		esperado bool
	}{
		{"versão maior", "v3.1.0", "v3.0.4", true},
		{"versão igual", "v3.0.4", "v3.0.4", false},
		{"versão menor", "v3.0.0", "v3.0.4", false},
		{"sem prefixo v", "3.1.0", "3.0.4", true},
		{"misturando prefixos", "v3.1.0", "3.0.4", true},
		{"semver lexical seria errado", "v3.10.0", "v3.2.0", true},
		{"current dev (não-semver) sempre desatualizado", "v3.0.0", "dev", true},
		{"latest inválido nunca atualiza", "garbage", "v3.0.0", false},
		{"pre-release < release", "v3.0.0-rc.1", "v3.0.0", false},
		{"release > pre-release current", "v3.0.0", "v3.0.0-rc.1", true},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			got := IsNewer(tt.latest, tt.current)
			if got != tt.esperado {
				t.Errorf("IsNewer(%q, %q) = %v, esperado %v",
					tt.latest, tt.current, got, tt.esperado)
			}
		})
	}
}

// TestCheckLatestRelease_Sucesso verifica que o checker faz parse correto
// do JSON do GitHub e retorna o Release populado.
func TestCheckLatestRelease_Sucesso(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/owner/repo/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"tag_name": "v3.1.0",
			"name": "v3.1.0 - Release de teste",
			"html_url": "https://github.com/owner/repo/releases/tag/v3.1.0",
			"published_at": "2026-05-01T12:00:00Z",
			"body": "## Novidades\n- Auto-update\n"
		}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := &Checker{Repo: "owner/repo", BaseURL: srv.URL}
	rel, err := c.CheckLatestRelease(context.Background())
	if err != nil {
		t.Fatalf("CheckLatestRelease erro inesperado: %v", err)
	}
	if rel.Version != "v3.1.0" {
		t.Errorf("Version = %q, esperado %q", rel.Version, "v3.1.0")
	}
	if rel.URL != "https://github.com/owner/repo/releases/tag/v3.1.0" {
		t.Errorf("URL inesperado: %q", rel.URL)
	}
	if rel.Body == "" {
		t.Errorf("Body vazio")
	}
	if rel.PublishedAt == "" {
		t.Errorf("PublishedAt vazio")
	}
}

// TestCheckLatestRelease_ParseiaAssets confirma que o checker captura a
// lista de binários publicados (DMG, ZIPs) — base do botão de download.
func TestCheckLatestRelease_ParseiaAssets(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/owner/repo/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{
			"tag_name": "v3.1.0",
			"html_url": "https://example.com/r/v3.1.0",
			"published_at": "2026-05-01T12:00:00Z",
			"body": "x",
			"assets": [
				{"name": "cfs-spool-darwin-universal.dmg", "browser_download_url": "https://example.com/dmg", "size": 12345},
				{"name": "cfs-spool-windows-amd64.zip", "browser_download_url": "https://example.com/win", "size": 9999}
			]
		}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := &Checker{Repo: "owner/repo", BaseURL: srv.URL}
	rel, err := c.CheckLatestRelease(context.Background())
	if err != nil {
		t.Fatalf("CheckLatestRelease erro: %v", err)
	}
	if len(rel.Assets) != 2 {
		t.Fatalf("Assets len = %d, esperado 2", len(rel.Assets))
	}
	if rel.Assets[0].Name != "cfs-spool-darwin-universal.dmg" {
		t.Errorf("Assets[0].Name = %q", rel.Assets[0].Name)
	}
	if rel.Assets[0].URL != "https://example.com/dmg" {
		t.Errorf("Assets[0].URL = %q", rel.Assets[0].URL)
	}
	if rel.Assets[1].Size != 9999 {
		t.Errorf("Assets[1].Size = %d", rel.Assets[1].Size)
	}
}

// TestPickAsset cobre a heurística de matching por SO usado pelo botão
// de download — escolhe o asset cujo nome contém o token do GOOS.
func TestPickAsset(t *testing.T) {
	assets := []ReleaseAsset{
		{Name: "cfs-spool-darwin-universal.dmg", URL: "u1"},
		{Name: "cfs-spool-linux-amd64.zip", URL: "u2"},
		{Name: "cfs-spool-windows-amd64.zip", URL: "u3"},
	}

	testes := []struct {
		goos     string
		esperado string
	}{
		{"darwin", "u1"},
		{"linux", "u2"},
		{"windows", "u3"},
		{"freebsd", ""}, // não suportado: nil
		{"", ""},
	}

	for _, tt := range testes {
		t.Run(tt.goos, func(t *testing.T) {
			a := PickAsset(assets, tt.goos)
			if tt.esperado == "" {
				if a != nil {
					t.Errorf("PickAsset(%q) = %+v, esperado nil", tt.goos, a)
				}
				return
			}
			if a == nil {
				t.Fatalf("PickAsset(%q) = nil, esperado URL %q", tt.goos, tt.esperado)
			}
			if a.URL != tt.esperado {
				t.Errorf("PickAsset(%q).URL = %q, esperado %q", tt.goos, a.URL, tt.esperado)
			}
		})
	}
}

// TestCheckLatestRelease_404 deve retornar erro silenciosamente, não panic.
func TestCheckLatestRelease_404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := &Checker{Repo: "owner/repo", BaseURL: srv.URL}
	_, err := c.CheckLatestRelease(context.Background())
	if err == nil {
		t.Errorf("CheckLatestRelease deveria retornar erro em 404")
	}
}

// TestCheckLatestRelease_RateLimit cobre 403 (rate limit do GitHub) — fallback silencioso.
func TestCheckLatestRelease_RateLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	c := &Checker{Repo: "owner/repo", BaseURL: srv.URL}
	_, err := c.CheckLatestRelease(context.Background())
	if err == nil {
		t.Errorf("CheckLatestRelease deveria retornar erro em 403")
	}
}

// TestCheckLatestRelease_DefaultBaseURL garante que sem BaseURL o checker
// usa api.github.com como default — não chamamos a API real, apenas verificamos
// a resolução do URL via uma requisição cancelada cedo.
func TestCheckLatestRelease_DefaultBaseURL(t *testing.T) {
	c := &Checker{Repo: "owner/repo"}
	if got := c.baseURL(); got != "https://api.github.com" {
		t.Errorf("baseURL default = %q, esperado %q", got, "https://api.github.com")
	}
}
