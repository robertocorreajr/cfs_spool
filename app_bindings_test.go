package main

import (
	"regexp"
	"strings"
	"testing"
)

// TestGetVersion_Default cobre o comportamento padrão da binding GetVersion
// quando o binário roda sem ldflags (caso típico em `go test`).
func TestGetVersion_Default(t *testing.T) {
	app := NewApp()
	if got := app.GetVersion(); got != "dev" {
		t.Errorf("GetVersion() default = %q, esperado %q", got, "dev")
	}
}

// TestGetVersion_Stubbed verifica que GetVersion repassa qualquer valor
// injetado em `version` (ldflags em release) e que o formato semântico
// vX.Y.Z passa pelo regex documentado na issue #25.
func TestGetVersion_Stubbed(t *testing.T) {
	original := version
	defer func() { version = original }()

	semver := regexp.MustCompile(`^v\d+\.\d+\.\d+(-[A-Za-z0-9.-]+)?$`)
	for _, v := range []string{"v3.0.0", "v3.1.2", "v10.0.0-rc.1"} {
		version = v
		got := NewApp().GetVersion()
		if got != v {
			t.Errorf("GetVersion() = %q, esperado %q", got, v)
		}
		if !semver.MatchString(got) {
			t.Errorf("GetVersion() = %q não bate com formato semântico", got)
		}
	}
}

// TestGetOptions verifica o contrato consumido pelo frontend:
// dropdowns nunca vazios e entrada CUSTOM preservada em lengths.
func TestGetOptions(t *testing.T) {
	app := NewApp()
	opts := app.GetOptions()

	if len(opts.Materials) == 0 {
		t.Error("GetOptions().Materials vazio — frontend ficaria sem materiais")
	}
	if len(opts.Vendors) == 0 {
		t.Error("GetOptions().Vendors vazio — frontend ficaria sem vendors")
	}
	if len(opts.Lengths) == 0 {
		t.Error("GetOptions().Lengths vazio — frontend ficaria sem comprimentos")
	}

	// Materiais e vendors expostos pela binding precisam refletir o backing array.
	if len(opts.Materials) != len(materials) {
		t.Errorf("GetOptions().Materials len=%d, esperado %d (sincronizado com materials)", len(opts.Materials), len(materials))
	}
	if len(opts.Vendors) != len(vendors) {
		t.Errorf("GetOptions().Vendors len=%d, esperado %d (sincronizado com vendors)", len(opts.Vendors), len(vendors))
	}
	if len(opts.Lengths) != len(lengths) {
		t.Errorf("GetOptions().Lengths len=%d, esperado %d (sincronizado com lengths)", len(opts.Lengths), len(lengths))
	}

	// Entrada CUSTOM é contrato com SpoolForm/LengthSelect (issue #25 + UX de gramas).
	achouCustom := false
	for _, l := range opts.Lengths {
		if l.Code == "CUSTOM" {
			achouCustom = true
			break
		}
	}
	if !achouCustom {
		t.Error("GetOptions().Lengths sem entrada CUSTOM — quebra UX de gramas customizadas")
	}
}

// TestWriteTag_ValidacaoCorInvalida garante que a binding falha cedo,
// antes de tocar o leitor PC/SC, quando a cor está malformada. É o que
// permite ao frontend exibir mensagem amigável sem travar em hardware.
func TestWriteTag_ValidacaoCorInvalida(t *testing.T) {
	app := NewApp()
	err := app.WriteTag(WriteRequest{
		Date:     "2026-04-12",
		Supplier: "0276",
		Material: "04001",
		Color:    "GG0000", // hex inválido
		Length:   "0330",
		Serial:   "1",
	})
	if err == nil {
		t.Fatal("WriteTag deveria retornar erro para cor inválida")
	}
	if !strings.Contains(err.Error(), "Cor inválida") {
		t.Errorf("erro = %q, esperado conter %q", err.Error(), "Cor inválida")
	}
}

// TestWriteTag_ValidacaoDataInvalida verifica que data malformada também
// faz early-return antes de qualquer chamada PC/SC.
func TestWriteTag_ValidacaoDataInvalida(t *testing.T) {
	app := NewApp()
	err := app.WriteTag(WriteRequest{
		Date:     "data-quebrada", // não é YYYY-MM-DD
		Supplier: "0276",
		Material: "04001",
		Color:    "FF4010",
		Length:   "0330",
		Serial:   "1",
	})
	if err == nil {
		t.Fatal("WriteTag deveria retornar erro para data inválida")
	}
	if !strings.Contains(err.Error(), "Data inválida") {
		t.Errorf("erro = %q, esperado conter %q", err.Error(), "Data inválida")
	}
}

// TestStopTagWatcher_QuandoOcioso garante que parar um watcher que nunca
// foi iniciado é no-op seguro — usado pelo WriteTag em fluxos de erro
// para garantir que reentradas não bloqueiem em canais nil.
func TestStopTagWatcher_QuandoOcioso(t *testing.T) {
	app := NewApp()
	// Não chamar StartTagWatcher — stopWatch deve estar nil.
	app.StopTagWatcher() // se entrasse no ramo do close, panicaria com canal nil.
	if app.stopWatch != nil || app.watchDone != nil {
		t.Error("StopTagWatcher ocioso não deveria mexer nos canais")
	}
}

// Nota: ReadTag, WriteTag (caminho feliz), StartTagWatcher e os handlers
// handleTagPresent/handleTagRemoved dependem do leitor PC/SC e do runtime
// Wails (wailsRuntime.EventsEmit). Para cobri-los serão necessárias duas
// abstrações:
//
//  1. Interface mockável em torno de internal/rfid (rastreado na issue #27).
//  2. Wrapper de runtime Wails injetável (rastreado em TECH_DEBT.md).
//
// Sem isso, qualquer teste desses caminhos depende de hardware real
// (ACR122U + tag CFS) e fica fora do alcance de `go test ./...`.
