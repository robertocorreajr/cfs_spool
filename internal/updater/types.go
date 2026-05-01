// Package updater consulta o GitHub Releases API e expõe utilitários para
// detectar quando uma versão mais recente do CFS Spool foi publicada.
//
// Fluxo principal:
//
//	checker := updater.NewChecker("robertocorreajr/cfs_spool")
//	release, err := checker.CheckLatestRelease(ctx)
//	if err == nil && updater.IsNewer(release.Version, currentVersion) {
//	    // Notificar usuário
//	}
//
// O pacote é deliberadamente silencioso em erros de rede / rate limit:
// chamadas falham com erro e o caller decide ignorar (não há panics nem
// retries automáticos).
package updater

// Release representa o subset de campos do GitHub Releases API que o
// CFS Spool consome. Os nomes JSON correspondem ao schema em
// https://docs.github.com/en/rest/releases/releases#get-the-latest-release
type Release struct {
	// Version é o tag_name da release (ex.: "v3.1.0"). Pode incluir prefixo "v".
	Version string `json:"version"`
	// Name é o título legível da release (ex.: "v3.1.0 - Auto-update").
	Name string `json:"name"`
	// URL aponta para a página HTML da release no GitHub.
	URL string `json:"url"`
	// PublishedAt é o timestamp ISO-8601 (ex.: "2026-05-01T12:00:00Z").
	PublishedAt string `json:"publishedAt"`
	// Body contém o changelog em Markdown.
	Body string `json:"body"`
	// Assets é a lista de binários publicados na release (DMG, ZIP, etc.).
	Assets []ReleaseAsset `json:"assets"`
}

// ReleaseAsset é cada arquivo anexado a uma release no GitHub —
// tipicamente o instalador por SO (cfs-spool-darwin-universal.dmg,
// cfs-spool-linux-amd64.zip, etc.).
type ReleaseAsset struct {
	// Name é o nome do arquivo (ex.: "cfs-spool-darwin-universal.dmg").
	Name string `json:"name"`
	// URL é o link direto para download via browser_download_url.
	URL string `json:"url"`
	// Size em bytes (informativo — mostrado no botão se relevante).
	Size int64 `json:"size"`
}
