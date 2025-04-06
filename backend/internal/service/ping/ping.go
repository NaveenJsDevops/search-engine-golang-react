package ping

import (
	"errors"

	"github.com/FenixAra/go-util/log"

	"go-parquet-read/internal/daos"
)

type Ping struct {
	l                *log.Logger
	ping             *daos.Ping
	parquetDirectory string
}

var (
	ErrUnableToPingDB = errors.New("unable to ping database")
)

func New(l *log.Logger, parquetDirectory string) *Ping {
	return &Ping{
		l:                l,
		parquetDirectory: parquetDirectory,
		ping:             daos.NewPing(l, parquetDirectory),
	}
}

func (p *Ping) Ping() error {
	ok, err := p.ping.Ping()
	if err != nil {
		return err
	}
	if !ok {
		return ErrUnableToPingDB
	}
	return nil
}
