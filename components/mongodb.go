package components

import (
	"errors"
	urllib "net/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type MongoDB struct {
	user string
	pass string

	client *mongo.Client
	db     *mongo.Database

	subscriber ISubscriber

	opts *CompOptions
}

func NewMongoDB() *MongoDB {
	m := new(MongoDB)
	m.opts = NewCompOptions()
	return m
}

func (m *MongoDB) Conn(url string, dbname string) error {
	mongoUrl := "mongodb://"
	mongoUrl += url
	if dbname != "" {
		mongoUrl += "/" + dbname
	}
	_, err := urllib.Parse(mongoUrl)
	if err != nil {
		return err
	}
	opts := &options.ClientOptions{}
	opts.ApplyURI(mongoUrl)
	opts.SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    dbname,
		Username:      m.user,
		Password:      m.pass})
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return err
	}
	m.client = client
	m.db = m.client.Database(dbname)
	if m.db.Name() != dbname {
		return errors.New("the database can not be found! " + dbname)
	}
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
	if m.client != nil {
		m.client.Disconnect(context.Background())
	}
}

func (m *MongoDB) RawDB() *mongo.Database {
	return m.db
}

func (m *MongoDB) Collect(colName string) *mongo.Collection {
	return m.db.Collection(colName)
}

func (m *MongoDB) Subscribe(colName string) error {
	col := m.Collect(colName)
	ctx := context.Background()
	pipeline := make([]interface{}, 0, 1024)
	cur, err := col.Watch(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)
	m.subscriber.OnSubscribe(colName)
	for cur.Next(ctx) {
		elem := &bson.RawElement{}
		if err := cur.Decode(elem); err != nil {
			m.subscriber.OnSubError(colName, err)
		}
		m.subscriber.OnSubMessage(elem.Key(), elem.String())
	}
	if err := cur.Err(); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) BindSubscriber(subscriber ISubscriber) {
	m.subscriber = subscriber
}

func (m *MongoDB) NewObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}
