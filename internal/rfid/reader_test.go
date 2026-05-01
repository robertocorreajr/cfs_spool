package rfid

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// fakeAPDU é uma fila ordenada de respostas APDU usada para alimentar
// o Reader com cenários determinísticos. Cada chamada a transmit consome
// o primeiro item da fila e registra o comando recebido para asserções.
type fakeAPDU struct {
	queue []apduStep
	calls [][]byte
}

type apduStep struct {
	want []byte // se não-nulo, valida que o comando recebido começa com este prefixo
	resp []byte
	err  error
}

func (f *fakeAPDU) transmit(cmd []byte) ([]byte, error) {
	f.calls = append(f.calls, append([]byte(nil), cmd...))
	if len(f.queue) == 0 {
		return nil, fmt.Errorf("fakeAPDU: comando inesperado %X (fila vazia)", cmd)
	}
	step := f.queue[0]
	f.queue = f.queue[1:]
	if step.want != nil && !bytes.HasPrefix(cmd, step.want) {
		return nil, fmt.Errorf("fakeAPDU: cmd %X não começa com prefixo esperado %X", cmd, step.want)
	}
	return step.resp, step.err
}

// newReaderForTest cria um *Reader com transmit injetado — atalho para
// instanciar o tipo nos testes sem PC/SC real (Open() exige hardware).
func newReaderForTest(fake *fakeAPDU) *Reader {
	r := &Reader{}
	r.transmitFn = fake.transmit
	return r
}

// statusOK = sw1=0x90, sw2=0x00 (sucesso APDU). Helper para encadear bytes.
func statusOK() []byte { return []byte{0x90, 0x00} }

// statusErr = sw1=0x63, sw2=0x00 (falha genérica de operação).
func statusErr() []byte { return []byte{0x63, 0x00} }

// TestTransmit_DelegaParaTransmitFn confirma a indireção: transmit() do
// Reader chama o callback configurado em transmitFn (sem tocar r.card).
func TestTransmit_DelegaParaTransmitFn(t *testing.T) {
	resposta := []byte{0xDE, 0xAD, 0x90, 0x00}
	fake := &fakeAPDU{queue: []apduStep{{resp: resposta}}}
	r := newReaderForTest(fake)

	got, err := r.transmit([]byte{0x01, 0x02})
	if err != nil {
		t.Fatalf("transmit erro inesperado: %v", err)
	}
	if !bytes.Equal(got, resposta) {
		t.Errorf("transmit() = %X, esperado %X", got, resposta)
	}
	if len(fake.calls) != 1 || !bytes.Equal(fake.calls[0], []byte{0x01, 0x02}) {
		t.Errorf("transmitFn não recebeu o comando esperado, calls=%v", fake.calls)
	}
}

// TestTransmit_PropagaErro garante que erros do callback chegam ao chamador.
func TestTransmit_PropagaErro(t *testing.T) {
	want := errors.New("falha pcsc simulada")
	fake := &fakeAPDU{queue: []apduStep{{err: want}}}
	r := newReaderForTest(fake)

	if _, err := r.transmit([]byte{0xFF}); err == nil || err.Error() != want.Error() {
		t.Errorf("transmit erro = %v, esperado %v", err, want)
	}
}

// TestUID_Sucesso simula a resposta GET DATA do ACR122 (4 bytes + 9000).
func TestUID_Sucesso(t *testing.T) {
	uid := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	fake := &fakeAPDU{queue: []apduStep{
		{want: []byte{0xFF, 0xCA, 0x00, 0x00, 0x00}, resp: append(uid, statusOK()...)},
	}}
	r := newReaderForTest(fake)

	got, err := r.UID()
	if err != nil {
		t.Fatalf("UID erro: %v", err)
	}
	if got != "deadbeef" {
		t.Errorf("UID = %q, esperado %q", got, "deadbeef")
	}
}

// TestUID_FalhaSWAPDU cobre status word ≠ 0x90 (cartão respondeu mal).
func TestUID_FalhaSWAPDU(t *testing.T) {
	fake := &fakeAPDU{queue: []apduStep{
		{resp: []byte{0x6A, 0x82}}, // 6A82 = arquivo não encontrado
	}}
	r := newReaderForTest(fake)
	if _, err := r.UID(); err == nil {
		t.Error("UID deveria falhar quando SW != 0x9000")
	}
}

// TestUID_RespostaCurta cobre defesa contra resposta com menos de 2 bytes.
func TestUID_RespostaCurta(t *testing.T) {
	fake := &fakeAPDU{queue: []apduStep{{resp: []byte{0x90}}}}
	r := newReaderForTest(fake)
	if _, err := r.UID(); err == nil {
		t.Error("UID deveria falhar com resposta < 2 bytes")
	}
}

// TestUID_ErroTransmit confirma que erro de transporte é propagado.
func TestUID_ErroTransmit(t *testing.T) {
	fake := &fakeAPDU{queue: []apduStep{{err: errors.New("desconectado")}}}
	r := newReaderForTest(fake)
	if _, err := r.UID(); err == nil {
		t.Error("UID deveria propagar erro de transmit")
	}
}

// TestReadBlockDirect_Sucesso simula resposta de READ BINARY (16 bytes + 9000).
func TestReadBlockDirect_Sucesso(t *testing.T) {
	bloco := bytes.Repeat([]byte{0xAB}, 16)
	fake := &fakeAPDU{queue: []apduStep{
		{want: []byte{0xFF, 0xB0, 0x00, 0x04, 0x10}, resp: append(bloco, statusOK()...)},
	}}
	r := newReaderForTest(fake)

	got, err := r.ReadBlockDirect(4)
	if err != nil {
		t.Fatalf("ReadBlockDirect erro: %v", err)
	}
	if got != strings.ToUpper(hex.EncodeToString(bloco)) {
		t.Errorf("ReadBlockDirect = %q, esperado %q", got, strings.ToUpper(hex.EncodeToString(bloco)))
	}
}

// TestReadBlockDirect_RespostaCurta cobre defesa contra payload truncado.
func TestReadBlockDirect_RespostaCurta(t *testing.T) {
	fake := &fakeAPDU{queue: []apduStep{
		{resp: []byte{0x01, 0x02, 0x90, 0x00}}, // <16 bytes de dado
	}}
	r := newReaderForTest(fake)
	if _, err := r.ReadBlockDirect(4); err == nil {
		t.Error("ReadBlockDirect deveria falhar com resposta curta")
	}
}

// TestReadBlockDirect_StatusError simula falha do cartão (SW != 0x90).
func TestReadBlockDirect_StatusError(t *testing.T) {
	fake := &fakeAPDU{queue: []apduStep{
		{resp: append(bytes.Repeat([]byte{0x00}, 16), statusErr()...)},
	}}
	r := newReaderForTest(fake)
	if _, err := r.ReadBlockDirect(4); err == nil {
		t.Error("ReadBlockDirect deveria falhar com SW=0x6300")
	}
}

// TestWriteBlock_Sucesso encadeia 2 APDUs: AUTHENTICATE + WRITE.
func TestWriteBlock_Sucesso(t *testing.T) {
	fake := &fakeAPDU{queue: []apduStep{
		{want: []byte{0xFF, 0x86}, resp: statusOK()},                  // auth
		{want: []byte{0xFF, 0xD6, 0x00, 0x04, 0x10}, resp: statusOK()}, // write block 4
	}}
	r := newReaderForTest(fake)

	dados := strings.Repeat("AA", 16) // 32 hex
	if err := r.WriteBlock(4, KeyTypeA, "FFFFFFFFFFFF", dados); err != nil {
		t.Fatalf("WriteBlock erro: %v", err)
	}
	if len(fake.calls) != 2 {
		t.Errorf("esperado 2 chamadas APDU, obtido %d", len(fake.calls))
	}
}

// TestWriteBlock_AuthFalha aborta sem nunca chegar ao WRITE.
func TestWriteBlock_AuthFalha(t *testing.T) {
	fake := &fakeAPDU{queue: []apduStep{
		{want: []byte{0xFF, 0x86}, resp: statusErr()}, // auth falha
	}}
	r := newReaderForTest(fake)

	dados := strings.Repeat("AA", 16)
	if err := r.WriteBlock(4, KeyTypeA, "FFFFFFFFFFFF", dados); err == nil {
		t.Error("WriteBlock deveria falhar quando auth falha")
	}
	if len(fake.calls) != 1 {
		t.Errorf("WriteBlock deveria parar no auth, calls=%d", len(fake.calls))
	}
}

// TestWriteBlock_KeyMalformada não chega a chamar transmit (rejeita cedo).
func TestWriteBlock_KeyMalformada(t *testing.T) {
	fake := &fakeAPDU{}
	r := newReaderForTest(fake)
	if err := r.WriteBlock(4, KeyTypeA, "ZZ", strings.Repeat("AA", 16)); err == nil {
		t.Error("WriteBlock deveria falhar com key inválida")
	}
	if len(fake.calls) != 0 {
		t.Errorf("WriteBlock não deveria chamar transmit com key inválida, calls=%d", len(fake.calls))
	}
}

// TestWriteBlock_DataMalformada falha após o auth (validação após transmit).
// O comportamento atual gasta 1 chamada (auth) antes de detectar o erro;
// o teste documenta esse comportamento real.
func TestWriteBlock_DataMalformada(t *testing.T) {
	fake := &fakeAPDU{queue: []apduStep{
		{want: []byte{0xFF, 0x86}, resp: statusOK()}, // auth ok
	}}
	r := newReaderForTest(fake)
	if err := r.WriteBlock(4, KeyTypeA, "FFFFFFFFFFFF", "AABB"); err == nil {
		t.Error("WriteBlock deveria falhar com data <32 hex")
	}
}

// TestTryReadBlock_PrimeiroMetodoSucesso cobre o caminho feliz: AUTH + READ.
func TestTryReadBlock_PrimeiroMetodoSucesso(t *testing.T) {
	bloco := bytes.Repeat([]byte{0x42}, 16)
	fake := &fakeAPDU{queue: []apduStep{
		{want: []byte{0xFF, 0x86}, resp: statusOK()},                  // auth
		{want: []byte{0xFF, 0xB0, 0x00, 0x04, 0x10}, resp: append(bloco, statusOK()...)},
	}}
	r := newReaderForTest(fake)

	got, err := r.TryReadBlock(4, KeyTypeA, "FFFFFFFFFFFF")
	if err != nil {
		t.Fatalf("TryReadBlock erro: %v", err)
	}
	if got != strings.ToUpper(hex.EncodeToString(bloco)) {
		t.Errorf("TryReadBlock = %q, esperado %q", got, strings.ToUpper(hex.EncodeToString(bloco)))
	}
}

// TestTryReadBlock_FallbackLoadKey cobre o caminho alternativo:
// quando o método 1 (auth direto) falha, o Reader tenta LOAD KEY → AUTH → READ.
func TestTryReadBlock_FallbackLoadKey(t *testing.T) {
	bloco := bytes.Repeat([]byte{0x55}, 16)
	fake := &fakeAPDU{queue: []apduStep{
		// Método 1: auth falha
		{want: []byte{0xFF, 0x86}, resp: statusErr()},
		// Método 2: load key OK, auth OK, read OK
		{want: []byte{0xFF, 0x82, 0x00, 0x00, 0x06}, resp: statusOK()},
		{want: []byte{0xFF, 0x86}, resp: statusOK()},
		{want: []byte{0xFF, 0xB0, 0x00, 0x04, 0x10}, resp: append(bloco, statusOK()...)},
	}}
	r := newReaderForTest(fake)

	got, err := r.TryReadBlock(4, KeyTypeA, "FFFFFFFFFFFF")
	if err != nil {
		t.Fatalf("TryReadBlock fallback erro: %v", err)
	}
	if got != strings.ToUpper(hex.EncodeToString(bloco)) {
		t.Errorf("TryReadBlock fallback = %q, esperado %q", got, strings.ToUpper(hex.EncodeToString(bloco)))
	}
}

// TestTryReadBlock_TudoFalha confirma que esgotar todos os métodos
// retorna erro com mensagem contendo o número do bloco.
func TestTryReadBlock_TudoFalha(t *testing.T) {
	fake := &fakeAPDU{queue: []apduStep{
		{resp: statusErr()}, // método 1: auth falha
		{resp: statusOK()},  // método 2: load key ok
		{resp: statusErr()}, // método 2: auth falha
	}}
	r := newReaderForTest(fake)

	_, err := r.TryReadBlock(4, KeyTypeA, "FFFFFFFFFFFF")
	if err == nil {
		t.Fatal("TryReadBlock deveria falhar quando todos os métodos falham")
	}
	if !strings.Contains(err.Error(), "bloco 4") {
		t.Errorf("erro = %q, esperado conter %q", err.Error(), "bloco 4")
	}
}

// TestDeriveKeyFromUID_FallbackEmErro garante que o método em *Reader
// devolve a chave padrão quando o UID é malformado (não panica).
func TestDeriveKeyFromUID_FallbackEmErro(t *testing.T) {
	r := &Reader{}
	got := r.DeriveKeyFromUID("xyz") // UID inválido
	if got != "FFFFFFFFFFFF" {
		t.Errorf("DeriveKeyFromUID com UID inválido = %q, esperado fallback FFFFFFFFFFFF", got)
	}
}

// Nota: Open() abre contexto PC/SC real e não é testável sem hardware
// (ou sem refatoração mais agressiva expondo um factory de scard.Context).
// Os testes acima cobrem a lógica APDU de todos os métodos públicos do
// Reader; cobrir Open()/Close() ficará para uma futura abstração.
