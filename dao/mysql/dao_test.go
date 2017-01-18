package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

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
	dsn := getTestDsn()
	dao = getTestDao(dsn)
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
		if err == sql.ErrNoRows {
			fmt.Println("no rows: " + err.Error())
		} else {
			fmt.Println("row scan error: " + err.Error())
		}
	} else {
		fmt.Println(item)
	}
}

func TestDaoTrans(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	dao.Begin()

	row := dao.QueryRow("SELECT * FROM test_mysql WHERE id = ?", 1)
	item := new(tableTestMysqlRowItem)
	err := row.Scan(&item.Id, &item.AddTime, &item.EditTime, &item.Name)
	if err != nil {
		fmt.Println("row scan error: " + err.Error())
	} else {
		fmt.Println(item)
	}

	dao.Commit()

	err = dao.Rollback()
	fmt.Println(err)
}

func TestDaoInsert(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	r, e := dao.Insert(
		TABLE_NAME,
		[]string{"name"},
		[]interface{}{"a"},
		[]interface{}{"b"},
		[]interface{}{"c"},
	)

	fmt.Println(r, e)
}

func TestDaoDeleteById(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	r, e := dao.DeleteById(TABLE_NAME, 1)

	fmt.Println(r, e)
}

func TestDaoUpdateById(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	r, e := dao.UpdateById(
		TABLE_NAME,
		7,
		NewColQueryItem("name", "", "e"),
	)

	fmt.Println(r, e)
}

func TestDaoSelectById(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	r := dao.SelectById(
		"*",
		TABLE_NAME,
		7,
	)

	fmt.Println(r)
}

func TestDaoSelectByIds(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	rows, err := dao.SelectByIds(
		"*",
		TABLE_NAME,
		[]interface{}{5, 7},
	)

	fmt.Println(err)
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

func TestDsnSetReadTimeout(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	dsn := getTestDsn()
	dsn.SetReadTimeout(1 * time.Microsecond)
	d := getTestDao(dsn)

	//should get io timeout
	_, err := d.Query("SELECT * FROM test_mysql")
	if err != nil {
		fmt.Println("get timeout successfully:" + err.Error())
	} else {
		t.Fatal("timeout not occur!")
	}

}

func TestDsnSetWriteTimeout(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	dsn := getTestDsn()
	dsn.SetWriteTimeout(1 * time.Microsecond)
	d := getTestDao(dsn)

	//should get io timeout
	_, err := d.Query("SELECT * FROM test_mysql")
	if err != nil {
		fmt.Println("get timeout successfully:" + err.Error())
	} else {
		t.Fatal("timeout not occur!")
	}
}

func TestDsnSetDialTimeout(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	dsn := getTestDsn()
	dsn.SetDialTimeout(1 * time.Microsecond)
	d := getTestDao(dsn)

	//should get dial timeout
	_, err := d.Query("SELECT * FROM test_mysql")
	if err != nil {
		fmt.Println("get timeout successfully:" + err.Error())
	} else {
		t.Fatal("timeout not occur!")
	}
}

func TestDsnInterpolateParams(t *testing.T) {
	misc.PrintCallerFuncNameForTest()

	dsn := getTestDsn()
	dsn.DisableInterpolateParams()
	status, ok := dsn.config.Params["interpolateParams"]
	if !ok || status != "false" {
		t.Fatal("disable interpolateParams error")
	}
	fmt.Println(dsn)
	dsn.EnableInterpolateParams()
	status, ok = dsn.config.Params["interpolateParams"]
	if !ok || status != "true" {
		t.Fatal("disable interpolateParams error")
	}
	fmt.Println(dsn)
}

func getTestDao(dsn *Dsn) *Dao {
	path := "/tmp/test_mysql.log"

	w, _ := writer.NewFileWriter(path)
	logger, _ := log.NewSimpleLogger(w, log.LEVEL_INFO, new(log.SimpleFormater))

	d, _ := NewDao(dsn, logger)

	return d
}

func getTestDsn() *Dsn {
	return NewDsn("root", "123", "127.0.0.1", "3306", "test")
}
