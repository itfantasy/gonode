package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type SqlBuilder struct {
	db     *sql.DB
	table  string
	option int
	optval map[string]string
	selval string
	where  string
	etc    string
}

const opt_SELECT int = 0
const opt_INSERT int = 1
const opt_DELETE int = 2
const opt_UPDATE int = 3

func (bd *SqlBuilder) doTable(db *sql.DB, tab string) *SqlBuilder {
	bd.db = db
	bd.table = tab
	return bd
}

func (bd *SqlBuilder) Insert(val map[string]string) *SqlBuilder {
	bd.option = opt_INSERT
	bd.optval = val
	return bd
}

func (bd *SqlBuilder) Update(val map[string]string) *SqlBuilder {
	bd.option = opt_UPDATE
	bd.optval = val
	return bd
}

func (bd *SqlBuilder) Delete() *SqlBuilder {
	bd.option = opt_DELETE
	bd.optval = nil
	return bd
}

func (bd *SqlBuilder) Select(sel string) *SqlBuilder {
	bd.option = opt_SELECT
	bd.selval = sel
	return bd
}

func (bd *SqlBuilder) Where(field string, compare string, val string) *SqlBuilder {
	bd.where = fmt.Sprintf("`%s` %s '%s'", field, compare, val)
	return bd
}

func (bd *SqlBuilder) And(field string, compare string, val string) *SqlBuilder {
	bd.where += fmt.Sprintf(" and `%s` %s '%s'", field, compare, val)
	return bd
}

func (bd *SqlBuilder) Or(field string, compare string, val string) *SqlBuilder {
	bd.where += fmt.Sprintf(" or `%s` %s '%s'", field, compare, val)
	return bd
}

func (bd *SqlBuilder) Etc(etc string) {
	bd.etc += etc
}

func (bd *SqlBuilder) close() {
	bd.table = ""
	bd.option = 0
	bd.optval = nil
	bd.selval = ""
	bd.where = ""
	bd = nil
}

func (bd *SqlBuilder) buildSql() (string, error) {

	if bd.option == opt_INSERT {
		var (
			temp  string
			temp2 string
		)
		for k, v := range bd.optval {
			if temp != "" {
				temp += ", "
				temp2 += ", "
			}
			temp += "`" + k + "`"
			if strings.Index(v, "#") == 0 {
				temp2 += strings.Replace(temp2, "#", "", 1)
			} else {
				temp2 += "'" + v + "'"
			}
		}
		ret := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", bd.table, temp, temp2)
		return ret, nil
	} else if bd.option == opt_UPDATE {
		var (
			temp  string
			temp2 string
		)
		for k, v := range bd.optval {
			if temp != "" {
				temp += ", "
				temp2 += ", "
			}
			temp += "`" + k + "`"
			if strings.Index(v, "#") == 0 {
				temp2 += strings.Replace(temp2, "#", "", 1)
			} else {
				temp2 += "'" + v + "'"
			}
			temp += "=" + temp2
			temp2 = ""
		}
		ret := fmt.Sprintf("UPDATE `%s` SET '%s' WHERE %s", bd.table, temp, bd.where)
		return ret, nil

	} else if bd.option == opt_DELETE {

		ret := fmt.Sprintf("DELETE FROM `%s` WHERE %s", bd.table, bd.where)
		return ret, nil

	} else if bd.option == opt_SELECT {

		ret := fmt.Sprintf("SELECT %s FROM `%s`", bd.selval, bd.table)
		if bd.where != "" {
			ret += " WHERE " + bd.where
		}

		fmt.Println(ret)
		return ret, nil
	}

	return "", errors.New("未知的操作类型!")
}

func (bd *SqlBuilder) Fetch() (*sql.Rows, error) {

	if bd.db == nil {
		return nil, errors.New("数据库链接尚未初始化!")
	}

	str, err2 := bd.buildSql()
	if err2 != nil {
		return nil, err2
	}

	rows, err := bd.db.Query(str + bd.etc)
	bd.close()
	return rows, err
}

func (bd *SqlBuilder) Query() (sql.Result, error) {

	if bd.db == nil {
		return nil, errors.New("数据库链接尚未初始化!")
	}

	str, err2 := bd.buildSql()
	if err2 != nil {
		return nil, err2
	}

	ret, err := bd.db.Exec(str + bd.etc)
	bd.close()
	return ret, err
}
