package main

import (
	"context"
	"testing"
)

// emittedEvent registra uma chamada a emitEvent (event name + payload)
// para asserções determinísticas em testes.
type emittedEvent struct {
	name string
	data []interface{}
}

// withEventSpy substitui a variável de pacote emitEvent por um spy
// e devolve a fila capturada + uma função para restaurar o emit original.
func withEventSpy(t *testing.T) (*[]emittedEvent, func()) {
	t.Helper()
	original := emitEvent
	captured := make([]emittedEvent, 0)
	emitEvent = func(_ context.Context, name string, data ...interface{}) {
		captured = append(captured, emittedEvent{name: name, data: data})
	}
	restore := func() { emitEvent = original }
	return &captured, restore
}

// TestHandleTagRemoved verifica que ao remover a tag o handler limpa o
// lastUID e emite (\"tag:status\", \"waiting\") — UI volta a barra amarela.
func TestHandleTagRemoved(t *testing.T) {
	captured, restore := withEventSpy(t)
	defer restore()

	app := NewApp()
	app.lastUID = "DEADBEEF" // simula tag previamente lida

	app.handleTagRemoved()

	if app.lastUID != "" {
		t.Errorf("lastUID = %q, esperado vazio após remoção", app.lastUID)
	}
	if len(*captured) != 1 {
		t.Fatalf("esperado 1 evento emitido, obtido %d", len(*captured))
	}
	ev := (*captured)[0]
	if ev.name != "tag:status" {
		t.Errorf("evento emitido = %q, esperado %q", ev.name, "tag:status")
	}
	if len(ev.data) != 1 || ev.data[0] != "waiting" {
		t.Errorf("payload = %v, esperado [\"waiting\"]", ev.data)
	}
}

// TestHandleTagPresent_QuandoReadTagFalha cobre o caminho onde a leitura
// PC/SC falha (sem hardware): o handler deve retornar early sem emitir
// nada, mantendo a UI no estado anterior.
func TestHandleTagPresent_QuandoReadTagFalha(t *testing.T) {
	captured, restore := withEventSpy(t)
	defer restore()

	app := NewApp()
	app.lastUID = "PRECEDENTE"

	// ReadTag chama rfid.Open() que falha sem hardware/PC/SC; em CI/lá
	// devolve um erro que o handler engole. Após a chamada, lastUID
	// permanece e nenhum evento é emitido.
	app.handleTagPresent()

	if app.lastUID != "PRECEDENTE" {
		t.Errorf("lastUID alterado para %q após ReadTag falhar", app.lastUID)
	}
	if len(*captured) != 0 {
		t.Errorf("nenhum evento deveria ser emitido em falha de ReadTag, obtido %d: %+v", len(*captured), *captured)
	}
}

// TestEmitEvent_DefaultDelegaParaWailsRuntime confirma que a variável
// emitEvent existe e é a indireção esperada (não nil). A delegação real
// para wailsRuntime.EventsEmit não é exercida aqui (exigiria runtime
// Wails real); o smoke test garante que a substituição não quebra a
// invocação básica.
func TestEmitEvent_NaoEhNil(t *testing.T) {
	if emitEvent == nil {
		t.Fatal("emitEvent não deveria ser nil — substituição quebrou inicialização")
	}
}
