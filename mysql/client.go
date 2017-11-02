package mysql

import (
	"andals/golog"

	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Client struct {
	config *Config

	db *sql.DB
	tx *sql.Tx

	connClosed bool

	logger golog.ILogger
}

func NewClient(config *Config, logger golog.ILogger) (*Client, error) {
	if config.LogLevel == 0 {
		config.LogLevel = golog.LEVEL_INFO
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = new(golog.NoopLogger)
	}

	return &Client{
		config: config,

		db: db,
		tx: nil,

		logger: logger,
	}, nil
}

func (this *Client) SetLogger(logger golog.ILogger) {
	this.logger = logger
}

func (this *Client) Closed() bool {
	return this.connClosed
}

func (this *Client) Free() {
	this.db.Close()
	this.tx = nil
	this.connClosed = true
}

func (this *Client) Exec(query string, args ...interface{}) (sql.Result, error) {
	this.log(query, args...)

	if this.tx != nil {
		return this.tx.Exec(query, args...)
	} else {
		return this.db.Exec(query, args...)
	}
}

func (this *Client) Query(query string, args ...interface{}) (*sql.Rows, error) {
	this.log(query, args...)

	if this.tx != nil {
		return this.tx.Query(query, args...)
	} else {
		return this.db.Query(query, args...)
	}
}

func (this *Client) QueryRow(query string, args ...interface{}) *sql.Row {
	this.log(query, args...)

	if this.tx != nil {
		return this.tx.QueryRow(query, args...)
	} else {
		return this.db.QueryRow(query, args...)
	}
}

func (this *Client) Begin() error {
	tx, err := this.db.Begin()
	if err != nil {
		return err
	}

	this.log("BEGIN")
	this.tx = tx

	return nil
}

func (this *Client) Commit() error {
	defer func() {
		this.tx = nil
	}()

	if this.tx != nil {
		this.log("COMMIT")

		return this.tx.Commit()
	}

	return errors.New("Not in trans")
}

func (this *Client) Rollback() error {
	defer func() {
		this.tx = nil
	}()

	if this.tx != nil {
		this.log("ROLLBACK")

		return this.tx.Rollback()
	}

	return errors.New("Not in trans")
}

func (this *Client) log(query string, args ...interface{}) {
	query = strings.Replace(query, "?", "%s", -1)
	vs := make([]interface{}, len(args))

	for i, v := range args {
		s := fmt.Sprint(v)
		switch v.(type) {
		case string:
			vs[i] = "'" + s + "'"
		default:
			vs[i] = s
		}
	}

	query = fmt.Sprintf(query, vs...)
	this.logger.Log(this.config.LogLevel, []byte(query))
}
