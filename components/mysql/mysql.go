package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itfantasy/gonode/components/etc"
)

type MySql struct {
	user string
	pass string
	db   *sql.DB
	opts *etc.CompOptions
}

func NewMySql() *MySql {
	this := new(MySql)
	this.user = "root"
	this.pass = ""
	return this
}

func (this *MySql) Conn(url string, dbname string) error {
	if this.db != nil {
		this.Close()
	}
	connstr := this.user + ":" + this.pass + "@tcp(" + url + ")/" + dbname
	db, err := sql.Open("mysql", connstr)
	if err != nil {
		return err
	}
	this.db = db
	return nil
}

func (this *MySql) SetAuthor(user string, pass string) {
	this.user = user
	this.pass = pass
}

func (this *MySql) SetOption(key string, val interface{}) {
	this.opts.Set(key, val)
}

func (this *MySql) Close() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *MySql) RawDB() *sql.DB {
	return this.db
}

func (this *MySql) Table(tab string) *SqlBuilder {
	bd := SqlBuilder{}
	bd.doTable(this.db, tab)
	return &bd
}
