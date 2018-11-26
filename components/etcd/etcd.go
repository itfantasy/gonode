package etcd

import (
	"strings"
	"time"

	"github.com/itfantasy/gonode/components/etc"
	"github.com/itfantasy/gonode/components/pubsub"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
)

const (
	OPT_CONNTIMEOUT string = "OPT_CONNTIMEOUT"
	OPT_RWTIMEOUT          = "OPT_RWTIMEOUT"
)

type Etcd struct {
	cli        *clientv3.Client
	subscriber pubsub.ISubscriber
	opts       *etc.CompOptions
}

func NewEtcd() *Etcd {
	this := new(Etcd)
	this.opts = etc.NewCompOptions()
	this.opts.Set(OPT_CONNTIMEOUT, 5*time.Second)
	this.opts.Set(OPT_RWTIMEOUT, time.Second)
	return this
}

func (this *Etcd) Conn(urls string, other string) error {
	urlInfos := strings.Split(urls, ";")
	endpoints := make([]string, 0, len(urlInfos))
	for _, v := range urlInfos {
		endpoints = append(endpoints, v)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: this.opts.Get(OPT_CONNTIMEOUT).(time.Duration),
	})
	if err != nil {
		return err
	}
	this.cli = cli
	return nil
}

func (this *Etcd) Close() {
	if this.cli != nil {
		this.cli.Close()
	}
}

func (this *Etcd) SetAuther(user string, pass string) {

}

func (this *Etcd) SetOption(key string, val string) {
	this.opts.Set(key, val)
}

func (this *Etcd) Set(key string, val string) error {
	ctx, cancel := context.WithTimeout(context.Background(), this.opts.Get(OPT_RWTIMEOUT).(time.Duration))
	_, err := this.cli.Put(ctx, key, val)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (this *Etcd) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), this.opts.Get(OPT_RWTIMEOUT).(time.Duration))
	resp, err := this.cli.Get(ctx, key)
	cancel()
	if err != nil {
		return "", err
	}
	val := ""
	for _, kv := range resp.Kvs {
		if string(kv.Key) == key {
			val = string(kv.Value)
		}
	}
	return val, nil
}

func (this *Etcd) Subscribe(key string) {
	ch := this.cli.Watch(context.Background(), key, clientv3.WithPrefix())
	this.subscriber.OnSubscribe(key)
	for resp := range ch {
		for _, ev := range resp.Events {
			if ev.Type == mvccpb.PUT {
				this.subscriber.OnSubMessage(string(ev.Kv.Key), string(ev.Kv.Value))
			} else if ev.Type == mvccpb.DELETE {
				this.subscriber.OnSubMessage(string(ev.Kv.Key), "")
			}
		}
	}
}

func (this *Etcd) BindSubscriber(subscriber pubsub.ISubscriber) {
	this.subscriber = subscriber
}

func (this *Etcd) Client() *clientv3.Client {
	return this.cli
}
