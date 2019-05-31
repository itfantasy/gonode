package mongodb

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"github.com/itfantasy/gonode/components/common"
)

type MongoDB struct {
	user    string
	pass    string
	session *mgo.Session
	db      *mgo.Database
	opts    *common.CompOptions
}

func NewMongoDB() *MongoDB {
	m := new(MongoDB)

	m.opts = common.NewCompOptions()
	return m
}

func (m *MongoDB) Conn(url string, dbname string) error {
	mongoUrl := "mongodb://"
	if m.user != "" {
		mongoUrl += m.user + ":" + m.pass + "@"
	}
	mongoUrl += url
	if dbname != "" {
		mongoUrl += "/" + dbname
	}
	fmt.Println(mongoUrl)
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return errors.New("session Dial failed!! " + err.Error())
	}
	m.session = session
	m.session.SetMode(mgo.Monotonic, true)
	db := session.DB(dbname)
	if db.Name != dbname {
		return errors.New("the database can not be found! " + dbname)
	}
	if m.user != "" {
		err2 := db.Login(m.user, m.pass)
		if err2 != nil {
			return errors.New("db author failed!! " + err2.Error())
		}
	}
	m.db = db
	return nil
}

func (m *MongoDB) SetAuthor(user string, pass string) {
	if user != "" {
		m.user = user
		m.pass = pass
	}
}

func (m *MongoDB) SetOption(key string, val interface{}) {

}

func (m *MongoDB) Close() {
	if m.session != nil {
		m.db.Logout()
		m.session.Close()
	}
}

func (m *MongoDB) RawSession() *mgo.Session {
	return m.session
}

func (m *MongoDB) Collect(name string) *mgo.Collection {
	return m.db.C(name)
}
