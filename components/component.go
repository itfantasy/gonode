package components

import (
	"github.com/itfantasy/gonode/components/mysql"
	"github.com/itfantasy/gonode/components/rabbitmq"
	"github.com/itfantasy/gonode/components/redis"
)

const (
	Redis    string = "Redis"
	MySql           = "MySql"
	MongoDB         = "MongoDB"
	RabbitMQ        = "RabbitMQ"
)

type IComponent interface {
	Conn(string, string) error
	Close()
	SetAuther(string, string)
	SetOption(string, string)
}

func NewRedis() *redis.Redis {
	return redis.NewRedis()
}

func NewMySql() *mysql.MySql {
	return mysql.NewMySql()
}

func NewRabbitMQ() *rabbitmq.RabbitMQ {
	return rabbitmq.NewRabbitMQ()
}
