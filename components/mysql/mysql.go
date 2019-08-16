package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itfantasy/gonode/components/common"
)

const (
	OPT_MAXASYNC string = "OPT_MAXASYNC"
)

type MySql struct {
	user    string
	pass    string
	db      *sql.DB
	opts    *common.CompOptions
	sqlchan chan string
}

func NewMySql() *MySql {
	m := new(MySql)
	m.user = "root"
	m.pass = ""
	m.opts = common.NewCompOptions()
	m.opts.Set(OPT_MAXASYNC, 0)
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
	maxasync := m.opts.GetInt(OPT_MAXASYNC)
	if maxasync > 0 {
		m.sqlchan = make(chan string, maxasync)
		go func() {
			for sqlstr := range m.sqlchan {
				if sqlstr == "EOF" {
					break
				}
				_, err := m.db.Exec(sqlstr)
				if err != nil {
					fmt.Println("[MySql]::Sql Async Exec faild!!" + sqlstr)
				}
			}
		}()
	}
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
		m.sqlchan <- "EOF"
		m.db.Close()
		m.db = nil
	}
}

func (m *MySql) RawDB() *sql.DB {
	return m.db
}

func (m *MySql) Table(tab string) *SqlBuilder {
	bd := SqlBuilder{}
	bd.doTable(m.db, m.sqlchan, tab)
	return &bd
}
