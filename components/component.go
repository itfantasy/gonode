package components

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
