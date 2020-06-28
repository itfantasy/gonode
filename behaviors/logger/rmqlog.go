package logger

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itfantasy/gonode/components"
)

type RabbitMQLogWriter struct {
	rmq     *components.RabbitMQ
	logchan chan *LogInfo
}

func (r *RabbitMQLogWriter) LogWrite(info *LogInfo) {
	if info == nil {
		return
	}
	info.Println()
	r.logchan <- info
}

func (r *RabbitMQLogWriter) Close() {
	r.logchan <- nil
}

func (r *RabbitMQLogWriter) dispose() {
	r.rmq.Close()
	close(r.logchan)
}

func NewRabbitMQLogWriter(rmq *components.RabbitMQ, logchan string) (*RabbitMQLogWriter, error) {
	r := new(RabbitMQLogWriter)
	r.rmq = rmq
	r.logchan = make(chan *LogInfo, 1024)
	go func() {
		defer r.dispose()
		for info := range r.logchan {
			if info == nil {
				break
			}
			js, err := json.Marshal(info)
			if err != nil {
				fmt.Fprint(os.Stderr, "RabbitMQLogWriter: %s", err)
				return
			}
			if err := r.rmq.Publish(logchan, string(js)); err != nil {
				fmt.Fprint(os.Stderr, "RabbitMQLogWriter: %s", err)
				return
			}
		}
	}()
	return r, nil
}
