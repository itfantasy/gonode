package etcd

import (
	"errors"
	"strings"
	"time"

	"github.com/itfantasy/gonode/components/other"
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
	opts       *other.CompOptions
	root       string
}

func NewEtcd() *Etcd {
	this := new(Etcd)
	this.opts = other.NewCompOptions()
	this.opts.Set(OPT_CONNTIMEOUT, 5*time.Second)
	this.opts.Set(OPT_RWTIMEOUT, time.Second)
	return this
}

func (this *Etcd) Conn(urls string, root string) error {
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
	this.root = root
	return nil
}

func (this *Etcd) Close() {
	if this.cli != nil {
		this.cli.Close()
	}
}

func (this *Etcd) SetAuthor(user string, pass string) {

}

func (this *Etcd) SetOption(key string, val interface{}) {
	this.opts.Set(key, val)
}

func (this *Etcd) Set(key string, val string) error {
	if this.root != "" {
		key = this.root + "/" + key
	}
	ctx, cancel := context.WithTimeout(context.Background(), this.opts.Get(OPT_RWTIMEOUT).(time.Duration))
	_, err := this.cli.Put(ctx, key, val)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (this *Etcd) Get(key string) (string, error) {
	dict, err := this.Gets(key)
	if err != nil {
		return "", err
	}
	val, exist := dict[key]
	if !exist {
		return "", errors.New("the key does not exist! " + key)
	}
	return val, nil
}

func (this *Etcd) Gets(key string) (map[string]string, error) {
	if this.root != "" {
		key = this.root + "/" + key
	}
	ctx, cancel := context.WithTimeout(context.Background(), this.opts.Get(OPT_RWTIMEOUT).(time.Duration))
	resp, err := this.cli.Get(ctx, key)
	cancel()
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	for _, kv := range resp.Kvs {
		ret[string(kv.Key)] = string(kv.Value)
	}
	return ret, nil
}

func (this *Etcd) Subscribe(key string) {
	if this.root != "" {
		key = this.root + "/" + key
	}
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
