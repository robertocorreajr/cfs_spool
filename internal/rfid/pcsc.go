package rfid

// Camada fina sobre github.com/ebfe/scard para permitir mock em testes.
// Apenas os métodos efetivamente usados por Reader.Open()/Close() são
// expostos pelas interfaces — abstrair GetStatusChange/ReaderState (usado
// pelo watcher em app.go) ficou deliberadamente de fora desta camada para
// não inflar o blast radius.

import "github.com/ebfe/scard"

// pcscContext é a interface mínima exposta a Reader.Open(); em produção
// é um wrapper sobre *scard.Context (ver scardContext) e em testes é
// um fake controlável.
type pcscContext interface {
	ListReaders() ([]string, error)
	Connect(reader string) (pcscCard, error)
	Release() error
}

// pcscCard é a interface mínima exposta para um cartão conectado.
type pcscCard interface {
	Transmit(cmd []byte) ([]byte, error)
	Disconnect() error
}

// pcscFactory abre um contexto PC/SC. Sobrescritível em testes via
// `defaultPCSCFactory = ...`.
type pcscFactory func() (pcscContext, error)

// defaultPCSCFactory é o ponto de extensão usado por Reader.Open().
// Em produção fala com o serviço PC/SC real; em testes é trocado por
// uma factory que devolve um fake de pcscContext.
var defaultPCSCFactory pcscFactory = realPCSCFactory

// realPCSCFactory abre o contexto via scard.EstablishContext e o
// embrulha em scardContext.
func realPCSCFactory() (pcscContext, error) {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return nil, err
	}
	return &scardContext{ctx: ctx}, nil
}

// scardContext envolve *scard.Context para satisfazer pcscContext.
type scardContext struct{ ctx *scard.Context }

func (s *scardContext) ListReaders() ([]string, error) {
	return s.ctx.ListReaders()
}

func (s *scardContext) Connect(reader string) (pcscCard, error) {
	card, err := s.ctx.Connect(reader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return nil, err
	}
	return &scardCard{card: card}, nil
}

func (s *scardContext) Release() error {
	return s.ctx.Release()
}

// scardCard envolve *scard.Card para satisfazer pcscCard.
type scardCard struct{ card *scard.Card }

func (s *scardCard) Transmit(cmd []byte) ([]byte, error) {
	return s.card.Transmit(cmd)
}

func (s *scardCard) Disconnect() error {
	return s.card.Disconnect(scard.LeaveCard)
}
