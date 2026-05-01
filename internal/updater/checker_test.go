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
