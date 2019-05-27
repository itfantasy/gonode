package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itfantasy/gonode/components/common"
)

type MySql struct {
	user string
	pass string
	db   *sql.DB
	opts *common.CompOptions
}

func NewMySql() *MySql {
	m := new(MySql)
	m.user = "root"
	m.pass = ""
	m.opts = common.NewCompOptions()
	return m
}

func (m *MySql) Conn(url string, dbname string) error {
	if m.db != nil {
		m.Close()
	}
	connstr := m.user + ":" + m.pass + "@tcp(" + url + ")/" + dbname
	db, err := sql.Open("mysql", connstr)
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

func (m *MySql) SetAuthor(user string, pass string) {
	m.user = user
	m.pass = pass
}

func (m *MySql) SetOption(key string, val interface{}) {
	m.opts.Set(key, val)
}

func (m *MySql) Close() {
	if m.db != nil {
		m.db.Close()
		m.db = nil
	}
}

func (m *MySql) RawDB() *sql.DB {
	return m.db
}

func (m *MySql) Table(tab string) *SqlBuilder {
	bd := SqlBuilder{}
	bd.doTable(m.db, tab)
	return &bd
}
