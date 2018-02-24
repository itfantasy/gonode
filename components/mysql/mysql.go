package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySql struct {
	db *sql.DB
}

func (this *MySql) Conn(url string, usr string, pass string, dbname string) error {
	if this.db != nil {
		this.Close()
	}
	connstr := usr + ":" + pass + "@tcp(" + url + ")/" + dbname
	fmt.Println(connstr)
	db, err := sql.Open("mysql", connstr)
	if err != nil {
		return err
	}
	this.db = db
	return nil
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
