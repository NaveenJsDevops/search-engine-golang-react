package daos

import (
	"github.com/FenixAra/go-util/log"
)

type PingCheckStruct struct {
	Count int32 `json:"count"`
}

type Ping struct {
	l                *log.Logger
	parquetDirectory string
}

func NewPing(l *log.Logger, parquetDirectory string) *Ping {
	return &Ping{
		l:                l,
		parquetDirectory: parquetDirectory,
	}
}

func (p *Ping) Ping() (bool, error) {
	pingModel := PingCheckStruct{}
	pingModel.Count = 1
	// Check Any issues In the directory
	return (pingModel.Count == 1), nil
}
