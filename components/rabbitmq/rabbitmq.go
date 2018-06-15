package rabbitmq

import (
	"github.com/itfantasy/gonode/components/etc"
	"github.com/itfantasy/gonode/components/pubsub"
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
	opts       *etc.CompOptions
	subscriber pubsub.ISubscriber
}

func NewRabbitMQ() *RabbitMQ {
	this := new(RabbitMQ)
	this.user = "guest"
	this.pass = "guest"
	this.queDict = make(map[string]*amqp.Queue)
	this.opts = etc.NewCompOptions()
	this.opts.Set(OPT_AUTOACK, true)
	return this
}

func (this *RabbitMQ) Conn(url string, host string) error {
	connStr := "amqp://" + this.user + ":" + this.pass + "@" + url + "/" + host
	conn, err := amqp.Dial(connStr)
	if err != nil {
		this.Close()
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		this.Close()
		return err
	}
	this.conn = conn
	this.ch = ch
	return nil
}

func (this *RabbitMQ) autoQueueDeclare(name string) (*amqp.Queue, error) {
	_, exist := this.queDict[name]
	if exist {
		return this.queDict[name], nil
	}
	q, err := this.ch.QueueDeclare(
		name,
		this.opts.GetBool(OPT_DURABLE),
		this.opts.GetBool(OPT_AUTODELETE),
		this.opts.GetBool(OPT_EXCLUSIVE),
		this.opts.GetBool(OPT_NOWAIT),
		this.opts.GetArgs(OPT_ARGS),
	)
	if err != nil {
		return nil, err
	}
	this.queDict[name] = &q
	return this.queDict[name], nil
}

func (this *RabbitMQ) SetAuther(user string, pass string) {
	this.user = user
	this.pass = pass
}

func (this *RabbitMQ) SetOption(key string, val interface{}) {
	this.opts.Set(key, val)
}

func (this *RabbitMQ) BindSubscriber(subscriber pubsub.ISubscriber) {
	this.subscriber = subscriber
}

func (this *RabbitMQ) Publish(que string, msg string) error {
	_, _err := this.autoQueueDeclare(que)
	if _err != nil {
		return _err
	}
	err := this.ch.Publish(
		this.opts.GetString(OPT_EXCHANGE), // exchange
		que, // routing key
		this.opts.GetBool(OPT_MANDATORY), // mandatory
		this.opts.GetBool(OPT_IMMEDIATE), // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	return err
}

func (this *RabbitMQ) Subscribe(que string) {
	_, err := this.autoQueueDeclare(que)
	if err != nil {
		if this.subscriber != nil {
			this.subscriber.OnSubError(que, err)
		}
	}
	if this.subscriber != nil {
		this.subscriber.OnSubscribe(que)
	}
	msgs, err := this.ch.Consume(
		que,
		this.opts.GetString(OPT_CONSUMER),
		this.opts.GetBool(OPT_AUTOACK),
		this.opts.GetBool(OPT_EXCLUSIVE),
		this.opts.GetBool(OPT_NOLOCAL),
		this.opts.GetBool(OPT_NOWAIT),
		this.opts.GetArgs(OPT_ARGS),
	)
	if err != nil {
		if this.subscriber != nil {
			this.subscriber.OnSubError(que, err)
		}
	}
	for d := range msgs {
		if this.subscriber != nil {
			this.subscriber.OnSubMessage(que, string(d.Body))
		}
	}
}

func (this *RabbitMQ) Close() {
	if this.conn != nil {
		this.conn.Close()
	}
	if this.ch != nil {
		this.ch.Close()
	}
}
