package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const defaultBaseURL = "https://api.github.com"

// Checker consulta o GitHub Releases API. BaseURL pode ser sobrescrita em
// testes (httptest.Server.URL) — em produção fica vazia e cai no default.
type Checker struct {
	// Repo no formato "owner/name" (ex.: "robertocorreajr/cfs_spool").
	Repo string
	// BaseURL é a raiz da API. Em produção: "https://api.github.com" (default).
	BaseURL string
	// HTTPClient permite injetar timeout/transport customizado. Default: 3s.
	HTTPClient *http.Client
}

// NewChecker cria um Checker com timeout de 3s — curto o suficiente para
// não atrasar o startup do app numa rede ruim.
func NewChecker(repo string) *Checker {
	return &Checker{
		Repo: repo,
		HTTPClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// baseURL devolve a URL configurada ou o default público.
func (c *Checker) baseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return defaultBaseURL
}

// CheckLatestRelease consulta GET /repos/{repo}/releases/latest e retorna
// o Release decodificado. Erros de rede, rate limit (403) ou repo não
// encontrado (404) viram erros normais — o caller engole silenciosamente.
func (c *Checker) CheckLatestRelease(ctx context.Context) (*Release, error) {
	if c.Repo == "" {
		return nil, fmt.Errorf("updater: repo não configurado")
	}

	url := fmt.Sprintf("%s/repos/%s/releases/latest", c.baseURL(), c.Repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("updater: criar request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 3 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("updater: HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("updater: status %d do GitHub", resp.StatusCode)
	}

	// Schema mínimo do GitHub — evita acoplar struct externa ao response oficial.
	var raw struct {
		TagName     string `json:"tag_name"`
		Name        string `json:"name"`
		HTMLURL     string `json:"html_url"`
		PublishedAt string `json:"published_at"`
		Body        string `json:"body"`
		Assets      []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("updater: decodificar JSON: %w", err)
	}

	assets := make([]ReleaseAsset, 0, len(raw.Assets))
	for _, a := range raw.Assets {
		assets = append(assets, ReleaseAsset{
			Name: a.Name,
			URL:  a.BrowserDownloadURL,
			Size: a.Size,
		})
	}

	return &Release{
		Version:     raw.TagName,
		Name:        raw.Name,
		URL:         raw.HTMLURL,
		PublishedAt: raw.PublishedAt,
		Body:        raw.Body,
		Assets:      assets,
	}, nil
}

// PickAsset escolhe o asset que casa com o SO recebido.
//
// Convenção dos releases do CFS Spool (ver auto-tag.yml + build.yml):
//   - darwin → contém "darwin" no nome (ex.: cfs-spool-darwin-universal.dmg)
//   - windows → contém "windows" no nome (ex.: cfs-spool-windows-amd64.zip)
//   - linux → contém "linux" no nome (ex.: cfs-spool-linux-amd64.zip)
//
// Devolve nil se não achar — caller decide cair pra release page no browser.
func PickAsset(assets []ReleaseAsset, goos string) *ReleaseAsset {
	keyword := goosKeyword(goos)
	if keyword == "" {
		return nil
	}
	lk := strings.ToLower(keyword)
	for i, a := range assets {
		if strings.Contains(strings.ToLower(a.Name), lk) {
			return &assets[i]
		}
	}
	return nil
}

// goosKeyword traduz o GOOS para o token usado nos nomes dos assets.
// Mantemos uma função separada para facilitar adicionar aliases futuros
// (ex.: "macos" virando sinônimo de "darwin").
func goosKeyword(goos string) string {
	switch strings.ToLower(goos) {
	case "darwin":
		return "darwin"
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	default:
		return ""
	}
}

// IsNewer compara duas versões usando semver e retorna true se latest > current.
// Aceita strings com ou sem prefixo "v" — internamente normaliza para "vX.Y.Z".
//
// Casos especiais:
//   - current não-semver (ex.: "dev"): considera latest mais novo (encoraja update).
//   - latest não-semver: retorna false (não dá para comparar; não notifica).
func IsNewer(latest, current string) bool {
	latest = normalizeSemver(latest)
	current = normalizeSemver(current)

	if !semver.IsValid(latest) {
		return false
	}
	if !semver.IsValid(current) {
		// Builds locais ("dev") sempre vêem o último release como mais novo —
		// o usuário ainda assim pode ignorar a notificação.
		return true
	}
	return semver.Compare(latest, current) > 0
}

// normalizeSemver garante prefixo "v" para que golang.org/x/mod/semver aceite.
func normalizeSemver(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return v
	}
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}
