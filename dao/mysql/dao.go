package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"andals/gobox/log"
)

const (
	DRIVER_NAME_MYSQL = "mysql"
)

type Dao struct {
	db     *sql.DB
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
		db:     db,
		logger: logger,
	}, nil
}

func (this *Dao) Exec(query string, args ...interface{}) (sql.Result, error) {
	this.logQuery(query, args...)

	return this.db.Exec(query, args...)
}

func (this *Dao) Query(query string, args ...interface{}) (*sql.Rows, error) {
	this.logQuery(query, args...)

	return this.db.Query(query, args...)
}

func (this *Dao) QueryRow(query string, args ...interface{}) *sql.Row {
	this.logQuery(query, args...)

	return this.db.QueryRow(query, args...)
}

func (this *Dao) logQuery(query string, args ...interface{}) {
	query = strings.Replace(query, "?", "%s", -1)
	vs := make([]interface{}, len(args))

	for i, v := range args {
		s := fmt.Sprint(v)
		switch v.(type) {
		case int:
			vs[i] = s
		case string:
			vs[i] = "'" + s + "'"
		}
	}

	query = fmt.Sprintf(query, vs...)
	this.logger.Info([]byte(query))
}
