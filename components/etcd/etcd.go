package etcd

import (
	"errors"
	"strings"
	"time"

	"github.com/itfantasy/gonode/components/common"
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
	subscriber common.ISubscriber
	opts       *common.CompOptions
	root       string
}

func NewEtcd() *Etcd {
	this := new(Etcd)
	this.opts = common.NewCompOptions()
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

func (this *Etcd) Set(path string, val string) error {
	if this.root != "" {
		path = this.root + "/" + path
	}
	ctx, cancel := context.WithTimeout(context.Background(), this.opts.Get(OPT_RWTIMEOUT).(time.Duration))
	_, err := this.cli.Put(ctx, path, val)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (this *Etcd) Get(path string) (string, error) {
	if this.root != "" {
		path = this.root + "/" + path
	}
	ctx, cancel := context.WithTimeout(context.Background(), this.opts.Get(OPT_RWTIMEOUT).(time.Duration))
	resp, err := this.cli.Get(ctx, path)
	cancel()
	if err != nil {
		return "", err
	}
	for _, kv := range resp.Kvs {
		if string(kv.Key) == path {
			return string(kv.Value), nil
		}
	}
	return "", errors.New("can not find the path! " + path)
}

func (this *Etcd) Gets(path string) (map[string]string, error) {
	if this.root != "" {
		path = this.root + "/" + path
	}
	ctx, cancel := context.WithTimeout(context.Background(), this.opts.Get(OPT_RWTIMEOUT).(time.Duration))
	resp, err := this.cli.Get(ctx, path, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		if this.root != "" {
			key = strings.TrimPrefix(key, this.root+"/")
		}
		ret[key] = string(kv.Value)
	}
	return ret, nil
}

func (this *Etcd) Publish(path string, val string) error {
	return this.Set(path, val)
}

func (this *Etcd) Subscribe(path string) {
	if this.root != "" {
		path = this.root + "/" + path
	}
	ch := this.cli.Watch(context.Background(), path, clientv3.WithPrefix())
	this.subscriber.OnSubscribe(strings.TrimPrefix(path, this.root+"/"))
	for resp := range ch {
		for _, ev := range resp.Events {
			if ev.Type == mvccpb.PUT {
				this.subscriber.OnSubMessage(strings.TrimPrefix(string(ev.Kv.Key), this.root+"/"), string(ev.Kv.Value))
			} else if ev.Type == mvccpb.DELETE {
				this.subscriber.OnSubMessage(strings.TrimPrefix(string(ev.Kv.Key), this.root+"/"), "")
			}
		}
	}
}

func (this *Etcd) BindSubscriber(subscriber common.ISubscriber) {
	this.subscriber = subscriber
}

func (this *Etcd) Client() *clientv3.Client {
	return this.cli
}
