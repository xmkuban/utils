package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/xmkuban/logger"
)

type MysqlDBWrap struct {
	Host     string
	Username string
	Password string
	Schema   string
	Charset  string
	Loc      string

	MaxConnections  int
	MaxConnLifetime int64

	d     *sql.DB
	valid bool
}

func (db *MysqlDBWrap) genMysqlConnString() string {
	//check if contains 3306
	hostArr := strings.Split(db.Host, ":")
	if len(hostArr) == 1 {
		db.Host = fmt.Sprintf("%s:3306", db.Host)
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&loc=%s&parseTime=true", db.Username, db.Password, db.Host, db.Schema, db.Charset, db.Loc)
}

func NewMysql(Host string, username string, password string, schema string) (*MysqlDBWrap, error) {
	db := &MysqlDBWrap{
		Host:     Host,
		Username: username,
		Password: password,
		Schema:   schema,
	}
	var err error
	if db.Charset == "" {
		db.Charset = "utf8mb4"
	}

	if db.Loc == "" {
		db.Loc = "UTC"
	}

	if db.d == nil {
		connString := db.genMysqlConnString()
		db.d, err = sql.Open("mysql", connString)
		if err == nil {
			err = db.d.Ping()
			if err == nil {
				db.valid = true
				return db, nil
			}
		} else {
			return nil, err
		}
	}
	if db.MaxConnLifetime != 0 {
		db.d.SetConnMaxLifetime(time.Duration(db.MaxConnLifetime) * time.Second)
	} else { // default one hour
		db.d.SetConnMaxLifetime(1 * time.Hour)
	}

	if db.MaxConnections != 0 {
		db.d.SetMaxIdleConns(db.MaxConnections)
		db.d.SetMaxOpenConns(db.MaxConnections)
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			err := db.d.Ping()
			if err != nil {
				logger.Errorf("sql ping db fail:%s", err)
			}
		}
	}()

	return db, err
}

func (db *MysqlDBWrap) SetLoc(loc string) {
	db.Loc = url.QueryEscape(loc)
}

func (db *MysqlDBWrap) Valid() bool {
	return db.valid
}

func (db *MysqlDBWrap) DB() *sql.DB {
	if db.valid {
		return db.d
	}
	return nil
}
