package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"andals/gobox/log"
)

const (
	DRIVER_NAME_MYSQL = "mysql"
)

type Dao struct {
	db  *sql.DB
	tx  *sql.Tx
	sqb *SimpleQueryBuilder

	logger log.ILogger
}

func DsnTcpIpV4(username, password, host, port, dbname string) string {
	return username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname
}

func NewDao(dsn string, logger log.ILogger) (*Dao, error) {
	db, err := sql.Open(DRIVER_NAME_MYSQL, dsn)
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = new(log.NoopLogger)
	}

	return &Dao{
		db:  db,
		tx:  nil,
		sqb: new(SimpleQueryBuilder),

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
	if this.tx != nil {
		this.logQuery("COMMIT")

		err := this.tx.Commit()
		this.tx = nil

		return err
	}

	return errors.New("Not in trans")
}

func (this *Dao) Rollback() error {
	if this.tx != nil {
		this.logQuery("ROLLBACK")

		err := this.tx.Rollback()
		this.tx = nil

		return err
	}

	return errors.New("Not in trans")
}

func (this *Dao) Insert(tableName string, colNames []string, colsValues ...[]interface{}) (sql.Result, error) {
	this.sqb.
		Insert(tableName, colNames...).
		Values(colsValues...)

	return this.Exec(this.sqb.Query(), this.sqb.Args()...)
}

func (this *Dao) DeleteById(tableName string, id uint64) (sql.Result, error) {
	this.sqb.
		Delete(tableName).
		WhereConditionAnd(NewColQueryItem("id", COND_EQUAL, id))

	return this.Exec(this.sqb.Query(), this.sqb.Args()...)
}

func (this *Dao) UpdateById(tableName string, id uint64, setItems ...*ColQueryItem) (sql.Result, error) {
	this.sqb.
		Update(tableName).
		Set(setItems...).
		WhereConditionAnd(NewColQueryItem("id", COND_EQUAL, id))

	return this.Exec(this.sqb.Query(), this.sqb.Args()...)
}

func (this *Dao) SelectById(what, tableName string, id uint64) *sql.Row {
	this.sqb.
		Select(what, tableName).
		WhereConditionAnd(NewColQueryItem("id", COND_EQUAL, id))

	return this.QueryRow(this.sqb.Query(), this.sqb.Args()...)
}

func (this *Dao) SelectByIds(what, tableName string, ids []uint64) (*sql.Rows, error) {
	is := make([]interface{}, len(ids))
	for k, v := range ids {
		is[k] = v
	}

	this.sqb.
		Select(what, tableName).
		WhereConditionAnd(NewColQueryItem("id", COND_IN, is...))

	return this.Query(this.sqb.Query(), this.sqb.Args()...)
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
