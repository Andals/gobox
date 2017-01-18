package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"andals/gobox/log"
)

const (
	DRIVER_NAME_MYSQL = "mysql"

	DEFAULT_DAIL_TIMEOUT  = 10 * time.Second
	DEFAULT_READ_TIMEOUT  = 10 * time.Second
	DEFAULT_WRITE_TIMEOUT = 10 * time.Second
)

type Dao struct {
	db *sql.DB
	tx *sql.Tx

	logger log.ILogger
}

type Dsn struct {
	config *mysql.Config
}

var defaultParams = map[string]string{
	"interpolateParams": "true",
}

func NewDsn(username, password, host, port, dbname string) *Dsn {
	config := mysql.Config{
		User:         username,
		Passwd:       password,
		Net:          "tcp",
		Addr:         host + ":" + port,
		DBName:       dbname,
		Params:       defaultParams,
		Timeout:      DEFAULT_DAIL_TIMEOUT,
		ReadTimeout:  DEFAULT_READ_TIMEOUT,
		WriteTimeout: DEFAULT_WRITE_TIMEOUT,
	}

	return &Dsn{
		config: &config,
	}
}

func (this *Dsn) SetReadTimeout(timeout time.Duration) *Dsn {
	this.config.ReadTimeout = timeout

	return this
}

func (this *Dsn) SetWriteTimeout(timeout time.Duration) *Dsn {
	this.config.WriteTimeout = timeout

	return this
}

func (this *Dsn) SetDialTimeout(timeout time.Duration) *Dsn {
	this.config.Timeout = timeout

	return this
}

func (this *Dsn) SetParam(key, value string) *Dsn {
	this.config.Params[key] = value

	return this
}

func (this *Dsn) EnableInterpolateParams() *Dsn {
	this.config.Params["interpolateParams"] = "true"

	return this
}

func (this *Dsn) DisableInterpolateParams() *Dsn {
	this.config.Params["interpolateParams"] = "false"

	return this
}

func (this *Dsn) String() string {
	return this.config.FormatDSN()
}

func NewDao(dsn *Dsn, logger log.ILogger) (*Dao, error) {
	db, err := sql.Open(DRIVER_NAME_MYSQL, dsn.String())
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = new(log.NoopLogger)
	}

	return &Dao{
		db: db,
		tx: nil,

		logger: logger,
	}, nil
}

func (this *Dao) Exec(query string, args ...interface{}) (sql.Result, error) {
	this.logQuery(query, args...)

	if this.tx != nil {
		return this.tx.Exec(query, args...)
	} else {
		return this.db.Exec(query, args...)
	}
}

func (this *Dao) Query(query string, args ...interface{}) (*sql.Rows, error) {
	this.logQuery(query, args...)

	if this.tx != nil {
		return this.tx.Query(query, args...)
	} else {
		return this.db.Query(query, args...)
	}
}

func (this *Dao) QueryRow(query string, args ...interface{}) *sql.Row {
	this.logQuery(query, args...)

	if this.tx != nil {
		return this.tx.QueryRow(query, args...)
	} else {
		return this.db.QueryRow(query, args...)
	}
}

func (this *Dao) Begin() error {
	tx, err := this.db.Begin()
	if err != nil {
		return err
	}

	this.logQuery("BEGIN")
	this.tx = tx

	return nil
}

func (this *Dao) Commit() error {
	defer func() {
		this.tx = nil
	}()

	if this.tx != nil {
		this.logQuery("COMMIT")

		return this.tx.Commit()
	}

	return errors.New("Not in trans")
}

func (this *Dao) Rollback() error {
	defer func() {
		this.tx = nil
	}()

	if this.tx != nil {
		this.logQuery("ROLLBACK")

		return this.tx.Rollback()
	}

	return errors.New("Not in trans")
}

func (this *Dao) Insert(tableName string, colNames []string, colsValues ...[]interface{}) (sql.Result, error) {
	sqb := new(SimpleQueryBuilder)
	sqb.Insert(tableName, colNames...).
		Values(colsValues...)

	return this.Exec(sqb.Query(), sqb.Args()...)
}

func (this *Dao) DeleteById(tableName string, id interface{}) (sql.Result, error) {
	sqb := new(SimpleQueryBuilder)
	sqb.Delete(tableName).
		WhereConditionAnd(NewColQueryItem("id", COND_EQUAL, id))

	return this.Exec(sqb.Query(), sqb.Args()...)
}

func (this *Dao) UpdateById(tableName string, id interface{}, setItems ...*ColQueryItem) (sql.Result, error) {
	sqb := new(SimpleQueryBuilder)
	sqb.Update(tableName).
		Set(setItems...).
		WhereConditionAnd(NewColQueryItem("id", COND_EQUAL, id))

	return this.Exec(sqb.Query(), sqb.Args()...)
}

func (this *Dao) SelectById(what, tableName string, id interface{}) *sql.Row {
	sqb := new(SimpleQueryBuilder)
	sqb.Select(what, tableName).
		WhereConditionAnd(NewColQueryItem("id", COND_EQUAL, id))

	return this.QueryRow(sqb.Query(), sqb.Args()...)
}

func (this *Dao) SelectByIds(what, tableName string, ids []interface{}) (*sql.Rows, error) {
	is := make([]interface{}, len(ids))
	for k, v := range ids {
		is[k] = v
	}

	sqb := new(SimpleQueryBuilder)
	sqb.Select(what, tableName).
		WhereConditionAnd(NewColQueryItem("id", COND_IN, is...))

	return this.Query(sqb.Query(), sqb.Args()...)
}

func (this *Dao) logQuery(query string, args ...interface{}) {
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
	this.logger.Info([]byte(query))
}
