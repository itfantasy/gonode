package logger

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itfantasy/gonode/components/rabbitmq"

	log "github.com/jeanphorn/log4go"
)

// This log writer sends output to a socket
type RabbitMQLogWriter chan *log.LogRecord

// This is the SocketLogWriter's output method
func (w RabbitMQLogWriter) LogWrite(rec *log.LogRecord) {
	w <- rec
}

func (w RabbitMQLogWriter) Close() {
	close(w)
}

func NewRabbitMQLogWriter(url string, host string, logchan string, user string, pass string) RabbitMQLogWriter {

	rmq := rabbitmq.NewRabbitMQ()
	rmq.SetAuther(user, pass)
	err := rmq.Conn(url, host)

	if err != nil {
		fmt.Fprintf(os.Stderr, "NewRabbitMQLogWriter(%q): %s\n", url+"/"+host, err)
		return nil
	}

	w := RabbitMQLogWriter(make(chan *log.LogRecord, log.LogBufferLength))

	go func() {
		defer func() {
			rmq.Close()
		}()

		for rec := range w {
			// Marshall into JSON
			js, err := json.Marshal(rec)
			if err != nil {
				fmt.Fprint(os.Stderr, "RabbitMQLogWriter(%q): %s", url+"/"+host, err)
				return
			}

			err2 := rmq.Publish(logchan, string(js))
			if err2 != nil {
				fmt.Fprint(os.Stderr, "RabbitMQLogWriter(%q): %s", url+"/"+host, err2)
				return
			}
		}
	}()

	return w
}
