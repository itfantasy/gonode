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
func (r RabbitMQLogWriter) LogWrite(rec *log.LogRecord) {
	r <- rec
}

func (r RabbitMQLogWriter) Close() {
	close(r)
}

func NewRabbitMQLogWriter(rmq *rabbitmq.RabbitMQ, logchan string) RabbitMQLogWriter {

	r := RabbitMQLogWriter(make(chan *log.LogRecord, log.LogBufferLength))

	go func() {
		defer func() {
			rmq.Close()
		}()

		for rec := range r {
			// Marshall into JSON
			js, err := json.Marshal(rec)
			if err != nil {
				fmt.Fprint(os.Stderr, "RabbitMQLogWriter: %s", err)
				return
			}

			err2 := rmq.Publish(logchan, string(js))
			if err2 != nil {
				fmt.Fprint(os.Stderr, "RabbitMQLogWriter: %s", err2)
				return
			}
		}
	}()

	return r
}
