package db

import (
	"github.com/xmkuban/logger"

	"github.com/go-sql-driver/mysql"
)

func init() {
	mysql.SetLogger(logger.GetMySQLLogger())
}

type BaseQuery struct {
	Key   string
	Value interface{}
}

type DocumentQuery struct {
	Match BaseQuery `json:"match,omitempty" `
	From  int       `json:"from"`
	Size  int       `json:"size"`
}

type DocumentRes struct {
	Count int64
	Res   []byte
}
