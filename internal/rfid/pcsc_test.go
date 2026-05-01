package rfid

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// fakePCSCContext implementa pcscContext para testes — controla
// determinísticamente o resultado de cada método chamado por Open()/Close().
type fakePCSCContext struct {
	listReadersResult []string
	listReadersErr    error
	connectErr        error
	connectFn         func(reader string) (pcscCard, error)
	releaseErr        error

	listReadersCalled int
	connectCalled     int
	releaseCalled     int

	connectedReaders []string // capturados para verificar qual reader foi escolhido
}

func (f *fakePCSCContext) ListReaders() ([]string, error) {
	f.listReadersCalled++
	return f.listReadersResult, f.listReadersErr
}

func (f *fakePCSCContext) Connect(reader string) (pcscCard, error) {
	f.connectCalled++
	f.connectedReaders = append(f.connectedReaders, reader)
	if f.connectFn != nil {
		return f.connectFn(reader)
	}
	if f.connectErr != nil {
		return nil, f.connectErr
	}
	return &fakePCSCCard{}, nil
}

func (f *fakePCSCContext) Release() error {
	f.releaseCalled++
	return f.releaseErr
}

// fakePCSCCard implementa pcscCard.
type fakePCSCCard struct {
	transmitResp     []byte
	transmitErr      error
	transmitCalls    [][]byte
	disconnectCalled int
	disconnectErr    error
}

func (f *fakePCSCCard) Transmit(cmd []byte) ([]byte, error) {
	f.transmitCalls = append(f.transmitCalls, append([]byte(nil), cmd...))
	return f.transmitResp, f.transmitErr
}

func (f *fakePCSCCard) Disconnect() error {
	f.disconnectCalled++
	return f.disconnectErr
}

// withFactory substitui defaultPCSCFactory pelo fake e devolve uma
// função para restaurar o original — uso típico via `defer`.
func withFactory(t *testing.T, f pcscFactory) func() {
	t.Helper()
	original := defaultPCSCFactory
	defaultPCSCFactory = f
	return func() { defaultPCSCFactory = original }
}

// TestOpen_FactoryFalha verifica propagação de erro da factory.
func TestOpen_FactoryFalha(t *testing.T) {
	want := errors.New("EstablishContext falhou")
	defer withFactory(t, func() (pcscContext, error) {
		return nil, want
	})()

	r, err := Open()
	if err == nil || err.Error() != want.Error() {
		t.Errorf("Open() erro = %v, esperado %v", err, want)
	}
	if r != nil {
		t.Error("Open() devolveu Reader não-nil em erro de factory")
	}
}

// TestOpen_SemReaders cobre o caso de PC/SC sem nenhum leitor anexado.
func TestOpen_SemReaders(t *testing.T) {
	ctx := &fakePCSCContext{listReadersResult: []string{}}
	defer withFactory(t, func() (pcscContext, error) { return ctx, nil })()

	r, err := Open()
	if err == nil || !strings.Contains(err.Error(), "nenhum leitor") {
		t.Errorf("Open() erro = %v, esperado conter 'nenhum leitor'", err)
	}
	if r != nil {
		t.Error("Open() devolveu Reader em ausência de leitores")
	}
	if ctx.releaseCalled != 1 {
		t.Errorf("Release() chamado %d vez(es), esperado 1 (libera contexto em erro)", ctx.releaseCalled)
	}
}

// TestOpen_ListReadersErro cobre falha ao listar leitores.
func TestOpen_ListReadersErro(t *testing.T) {
	ctx := &fakePCSCContext{listReadersErr: errors.New("driver indisponível")}
	defer withFactory(t, func() (pcscContext, error) { return ctx, nil })()

	r, err := Open()
	if err == nil {
		t.Fatal("Open() deveria falhar com erro de ListReaders")
	}
	if r != nil {
		t.Error("Open() devolveu Reader em erro de ListReaders")
	}
	if ctx.releaseCalled != 1 {
		t.Errorf("Release() chamado %d vez(es), esperado 1", ctx.releaseCalled)
	}
}

// TestOpen_ConnectFalha cobre falha ao conectar no leitor encontrado.
func TestOpen_ConnectFalha(t *testing.T) {
	ctx := &fakePCSCContext{
		listReadersResult: []string{"ACR122 PICC Interface"},
		connectErr:        errors.New("cartão removido"),
	}
	defer withFactory(t, func() (pcscContext, error) { return ctx, nil })()

	r, err := Open()
	if err == nil || !strings.Contains(err.Error(), "cartão removido") {
		t.Errorf("Open() erro = %v, esperado conter 'cartão removido'", err)
	}
	if r != nil {
		t.Error("Open() devolveu Reader em erro de Connect")
	}
	if ctx.releaseCalled != 1 {
		t.Errorf("Release() chamado %d vez(es), esperado 1", ctx.releaseCalled)
	}
}

// TestOpen_Sucesso cobre o caminho feliz: ListReaders devolve um leitor,
// Connect devolve um cartão; Reader fica utilizável e transmitFn delega
// para o cartão devolvido.
func TestOpen_Sucesso(t *testing.T) {
	card := &fakePCSCCard{transmitResp: []byte{0xAA, 0x90, 0x00}}
	ctx := &fakePCSCContext{
		listReadersResult: []string{"ACR122 PICC Interface"},
		connectFn:         func(string) (pcscCard, error) { return card, nil },
	}
	defer withFactory(t, func() (pcscContext, error) { return ctx, nil })()

	r, err := Open()
	if err != nil {
		t.Fatalf("Open() erro inesperado: %v", err)
	}
	if r == nil {
		t.Fatal("Open() devolveu Reader nil")
	}
	if ctx.connectCalled != 1 {
		t.Errorf("Connect chamado %d vez(es), esperado 1", ctx.connectCalled)
	}
	if len(ctx.connectedReaders) != 1 || ctx.connectedReaders[0] != "ACR122 PICC Interface" {
		t.Errorf("reader conectado = %v, esperado [ACR122 PICC Interface]", ctx.connectedReaders)
	}

	// transmitFn deve delegar para o card devolvido pela factory.
	resp, err := r.transmit([]byte{0xFF, 0xCA, 0x00, 0x00, 0x00})
	if err != nil {
		t.Fatalf("transmit falhou: %v", err)
	}
	if !bytes.Equal(resp, []byte{0xAA, 0x90, 0x00}) {
		t.Errorf("resposta transmit = %X, esperado AA9000", resp)
	}
	if len(card.transmitCalls) != 1 {
		t.Errorf("Card.Transmit chamado %d vez(es), esperado 1", len(card.transmitCalls))
	}
}

// TestClose_LiberaCardEContexto verifica que Close chama Disconnect e
// Release exatamente uma vez quando ambos foram inicializados.
func TestClose_LiberaCardEContexto(t *testing.T) {
	card := &fakePCSCCard{}
	ctx := &fakePCSCContext{}
	r := &Reader{ctx: ctx, card: card}

	r.Close()

	if card.disconnectCalled != 1 {
		t.Errorf("Disconnect chamado %d vez(es), esperado 1", card.disconnectCalled)
	}
	if ctx.releaseCalled != 1 {
		t.Errorf("Release chamado %d vez(es), esperado 1", ctx.releaseCalled)
	}
}

// TestClose_SemPanicComCamposNil garante que fechar um Reader meio-criado
// não panica (defesa contra Open() ter falhado em algum ponto).
func TestClose_SemPanicComCamposNil(t *testing.T) {
	r := &Reader{} // ctx e card nil
	defer func() {
		if rec := recover(); rec != nil {
			t.Errorf("Close() panicou com campos nil: %v", rec)
		}
	}()
	r.Close()
}
