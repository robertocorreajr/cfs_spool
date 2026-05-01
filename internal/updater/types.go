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
}
