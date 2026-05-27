package main

// acceptance_changelog_test.go — teste de aceitação para o critério 7 da história:
// CHANGELOG [Unreleased] ganha entrada Added/Changed cobrindo a expansão do catálogo v7.
//
// Este teste protege contra regressão futura onde alguém remove ou esquece a entrada
// de changelog ao fazer merge da feature.

import (
	"os"
	"strings"
	"testing"
)

// TestChangelog_UnreleasedContainsCatalogV7 verifica que CHANGELOG.md possui:
//   - Uma seção [Unreleased] ativa.
//   - Menção ao catálogo Creality Print v7 (expansão de materiais).
//   - Menção à ferramenta cmd/check-materials (critério 9).
//   - Menção à ordenação alfabética do dropdown (critério 5).
//   - Menção à correção de 00035 (critério 2 — entrada corrigida).
func TestChangelog_UnreleasedContainsCatalogV7(t *testing.T) {
	data, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		t.Fatalf("não foi possível ler CHANGELOG.md: %v", err)
	}
	content := string(data)

	// Extrai somente a seção [Unreleased] para não contar entradas de releases antigas.
	_, afterHeader, found := strings.Cut(content, "## [Unreleased]")
	if !found {
		t.Fatal("CHANGELOG.md não contém seção ## [Unreleased]")
	}

	// Seção termina no próximo ## ou no fim do arquivo.
	unreleased, _, _ := strings.Cut(afterHeader, "\n## [")

	provas := []struct {
		descricao string
		needle    string
	}{
		{
			"menção à expansão do catálogo v7 (Creality Print v7)",
			"Creality Print v7",
		},
		{
			"menção à ferramenta cmd/check-materials",
			"check-materials",
		},
		{
			"menção à ordenação alfabética",
			"alphabetical",
		},
		{
			"menção à correção do code 00035",
			"00035",
		},
	}

	for _, p := range provas {
		if !strings.Contains(unreleased, p.needle) {
			t.Errorf("CHANGELOG.md [Unreleased] não contém %s (buscado: %q)", p.descricao, p.needle)
		}
	}
}
