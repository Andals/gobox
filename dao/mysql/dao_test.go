package mysql

import (
	"fmt"
	"strconv"
	"testing"

	"andals/gobox/log"
	"andals/gobox/log/writer"
	"andals/gobox/misc"
)

type tableTestMysqlRowItem struct {
	Id       int
	AddTime  string
	EditTime string
	Name     string
}

var dao *Dao

func init() {
	dsn := DsnTcpIpV4("root", "123", "127.0.0.1", "3306", "test")

	path := "/tmp/test_mysql.log"

	w, _ := writer.NewFileWriter(path)
	logger, _ := log.NewSimpleLogger(w, log.LEVEL_INFO)

	dao, _ = NewDao(dsn, logger)
}

func TestDaoExec(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	result, err := dao.Exec("INSERT INTO test_mysql (name) VALUES (?)", "a")
	if err != nil {
		fmt.Println("exec error: " + err.Error())
	} else {
		li, err := result.LastInsertId()
		if err != nil {
			fmt.Println("lastInsertId error: " + err.Error())
		} else {
			fmt.Println("lastInsertId: " + strconv.FormatInt(li, 10))
		}

		rf, err := result.RowsAffected()
		if err != nil {
			fmt.Println("rowsAffected error: " + err.Error())
		} else {
			fmt.Println("rowsAffected: " + strconv.FormatInt(rf, 10))
		}
	}
}

func TestDaoQuery(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	rows, err := dao.Query("SELECT * FROM test_mysql WHERE id IN (?,?)", 1, 5)
	if err != nil {
		fmt.Println("query error: " + err.Error())
	} else {
		for rows.Next() {
			item := new(tableTestMysqlRowItem)
			err = rows.Scan(&item.Id, &item.AddTime, &item.EditTime, &item.Name)
			if err != nil {
				fmt.Println("rows scan error: " + err.Error())
			} else {
				fmt.Println(item)
			}
		}
	}
}

func TestDaoQueryRow(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	row := dao.QueryRow("SELECT * FROM test_mysql WHERE id = ?", 5)
	item := new(tableTestMysqlRowItem)
	err := row.Scan(&item.Id, &item.AddTime, &item.EditTime, &item.Name)
	if err != nil {
		fmt.Println("row scan error: " + err.Error())
	} else {
		fmt.Println(item)
	}
}
