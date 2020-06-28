package components

import (
	"github.com/streadway/amqp"
)

const (
	OPT_DURABLE    string = "OPT_DURABLE"
	OPT_AUTODELETE        = "OPT_AUTODELETE"
	OPT_EXCLUSIVE         = "OPT_EXCLUSIVE"
	OPT_NOWAIT            = "OPT_NOWAIT"
	OPT_ARGS              = "OPT_ARGS"
	OPT_EXCHANGE          = "OPT_EXCHANGE"
	OPT_MANDATORY         = "OPT_MANDATORY"
	OPT_IMMEDIATE         = "OPT_IMMEDIATE"
	OPT_CONSUMER          = "OPT_CONSUMER"
	OPT_AUTOACK           = "OPT_AUTOACK"
	OPT_NOLOCAL           = "OPT_NOLOCAL"
)

type RabbitMQ struct {
	user       string
	pass       string
	conn       *amqp.Connection
	ch         *amqp.Channel
	queDict    map[string]*amqp.Queue
	opts       *CompOptions
	subscriber ISubscriber
}

func NewRabbitMQ() *RabbitMQ {
	r := new(RabbitMQ)
	r.user = "guest"
	r.pass = "guest"
	r.queDict = make(map[string]*amqp.Queue)
	r.opts = NewCompOptions()
	r.opts.Set(OPT_AUTOACK, true)
	return r
}

func (r *RabbitMQ) Conn(url string, host string) error {
	connStr := "amqp://" + r.user + ":" + r.pass + "@" + url + "/" + host
	conn, err := amqp.Dial(connStr)
	if err != nil {
		r.Close()
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		r.Close()
		return err
	}
	r.conn = conn
	r.ch = ch
	return nil
}

func (r *RabbitMQ) autoQueueDeclare(name string) (*amqp.Queue, error) {
	_, exist := r.queDict[name]
	if exist {
		return r.queDict[name], nil
	}
	q, err := r.ch.QueueDeclare(
		name,
		r.opts.GetBool(OPT_DURABLE),
		r.opts.GetBool(OPT_AUTODELETE),
		r.opts.GetBool(OPT_EXCLUSIVE),
		r.opts.GetBool(OPT_NOWAIT),
		r.opts.GetArgs(OPT_ARGS),
	)
	if err != nil {
		return nil, err
	}
	r.queDict[name] = &q
	return r.queDict[name], nil
}

func (r *RabbitMQ) SetAuthor(user string, pass string) {
	r.user = user
	r.pass = pass
}

func (r *RabbitMQ) SetOption(key string, val interface{}) {
	r.opts.Set(key, val)
}

func (r *RabbitMQ) BindSubscriber(subscriber ISubscriber) {
	r.subscriber = subscriber
}

func (r *RabbitMQ) Publish(que string, msg string) error {
	_, _err := r.autoQueueDeclare(que)
	if _err != nil {
		return _err
	}
	err := r.ch.Publish(
		r.opts.GetString(OPT_EXCHANGE), // exchange
		que, // routing key
		r.opts.GetBool(OPT_MANDATORY), // mandatory
		r.opts.GetBool(OPT_IMMEDIATE), // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	return err
}

func (r *RabbitMQ) Subscribe(que string) {
	_, err := r.autoQueueDeclare(que)
	if err != nil {
		if r.subscriber != nil {
			r.subscriber.OnSubError(que, err)
		}
	}
	if r.subscriber != nil {
		r.subscriber.OnSubscribe(que)
	}
	msgs, err := r.ch.Consume(
		que,
		r.opts.GetString(OPT_CONSUMER),
		r.opts.GetBool(OPT_AUTOACK),
		r.opts.GetBool(OPT_EXCLUSIVE),
		r.opts.GetBool(OPT_NOLOCAL),
		r.opts.GetBool(OPT_NOWAIT),
		r.opts.GetArgs(OPT_ARGS),
	)
	if err != nil {
		if r.subscriber != nil {
			r.subscriber.OnSubError(que, err)
		}
	}
	for d := range msgs {
		if r.subscriber != nil {
			r.subscriber.OnSubMessage(que, string(d.Body))
		}
	}
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
	if r.ch != nil {
		r.ch.Close()
	}
}
