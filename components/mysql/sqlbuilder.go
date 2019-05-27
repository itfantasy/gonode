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

func (s *SqlBuilder) doTable(db *sql.DB, tab string) *SqlBuilder {
	s.db = db
	s.table = tab
	return s
}

func (s *SqlBuilder) Insert(val map[string]string) *SqlBuilder {
	s.option = opt_INSERT
	s.optval = val
	return s
}

func (s *SqlBuilder) Update(val map[string]string) *SqlBuilder {
	s.option = opt_UPDATE
	s.optval = val
	return s
}

func (s *SqlBuilder) Delete() *SqlBuilder {
	s.option = opt_DELETE
	s.optval = nil
	return s
}

func (s *SqlBuilder) Select(sel string) *SqlBuilder {
	s.option = opt_SELECT
	s.selval = sel
	return s
}

func (s *SqlBuilder) Where(field string, compare string, val string) *SqlBuilder {
	s.where = fmt.Sprintf("`%s` %s '%s'", field, compare, val)
	return s
}

func (s *SqlBuilder) And(field string, compare string, val string) *SqlBuilder {
	s.where += fmt.Sprintf(" and `%s` %s '%s'", field, compare, val)
	return s
}

func (s *SqlBuilder) Or(field string, compare string, val string) *SqlBuilder {
	s.where += fmt.Sprintf(" or `%s` %s '%s'", field, compare, val)
	return s
}

func (s *SqlBuilder) Etc(etc string) {
	s.etc += etc
}

func (s *SqlBuilder) close() {
	s.table = ""
	s.option = 0
	s.optval = nil
	s.selval = ""
	s.where = ""
	s = nil
}

func (s *SqlBuilder) buildSql() (string, error) {

	if s.option == opt_INSERT {
		var (
			temp  string
			temp2 string
		)
		for k, v := range s.optval {
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
		ret := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", s.table, temp, temp2)
		return ret, nil
	} else if s.option == opt_UPDATE {
		var (
			temp  string
			temp2 string
		)
		for k, v := range s.optval {
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
		ret := fmt.Sprintf("UPDATE `%s` SET '%s' WHERE %s", s.table, temp, s.where)
		return ret, nil

	} else if s.option == opt_DELETE {

		ret := fmt.Sprintf("DELETE FROM `%s` WHERE %s", s.table, s.where)
		return ret, nil

	} else if s.option == opt_SELECT {

		ret := fmt.Sprintf("SELECT %s FROM `%s`", s.selval, s.table)
		if s.where != "" {
			ret += " WHERE " + s.where
		}

		fmt.Println(ret)
		return ret, nil
	}

	return "", errors.New("未知的操作类型!")
}

func (s *SqlBuilder) Fetch() (*sql.Rows, error) {

	if s.db == nil {
		return nil, errors.New("数据库链接尚未初始化!")
	}

	str, err2 := s.buildSql()
	if err2 != nil {
		return nil, err2
	}

	rows, err := s.db.Query(str + s.etc)
	s.close()
	return rows, err
}

func (s *SqlBuilder) Query() (sql.Result, error) {

	if s.db == nil {
		return nil, errors.New("数据库链接尚未初始化!")
	}

	str, err2 := s.buildSql()
	if err2 != nil {
		return nil, err2
	}

	ret, err := s.db.Exec(str + s.etc)
	s.close()
	return ret, err
}
