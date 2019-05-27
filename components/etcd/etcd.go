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
	e := new(Etcd)
	e.opts = common.NewCompOptions()
	e.opts.Set(OPT_CONNTIMEOUT, 5*time.Second)
	e.opts.Set(OPT_RWTIMEOUT, time.Second)
	return e
}

func (e *Etcd) Conn(urls string, root string) error {
	urlInfos := strings.Split(urls, ";")
	endpoints := make([]string, 0, len(urlInfos))
	for _, v := range urlInfos {
		endpoints = append(endpoints, v)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: e.opts.Get(OPT_CONNTIMEOUT).(time.Duration),
	})
	if err != nil {
		return err
	}
	e.cli = cli
	e.root = root
	return nil
}

func (e *Etcd) Close() {
	if e.cli != nil {
		e.cli.Close()
	}
}

func (e *Etcd) SetAuthor(user string, pass string) {

}

func (e *Etcd) SetOption(key string, val interface{}) {
	e.opts.Set(key, val)
}

func (e *Etcd) Set(path string, val string) error {
	if e.root != "" {
		path = e.root + "/" + path
	}
	ctx, cancel := context.WithTimeout(context.Background(), e.opts.Get(OPT_RWTIMEOUT).(time.Duration))
	_, err := e.cli.Put(ctx, path, val)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (e *Etcd) Get(path string) (string, error) {
	if e.root != "" {
		path = e.root + "/" + path
	}
	ctx, cancel := context.WithTimeout(context.Background(), e.opts.Get(OPT_RWTIMEOUT).(time.Duration))
	resp, err := e.cli.Get(ctx, path)
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

func (e *Etcd) Gets(path string) (map[string]string, error) {
	if e.root != "" {
		path = e.root + "/" + path
	}
	ctx, cancel := context.WithTimeout(context.Background(), e.opts.Get(OPT_RWTIMEOUT).(time.Duration))
	resp, err := e.cli.Get(ctx, path, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		if e.root != "" {
			key = strings.TrimPrefix(key, e.root+"/")
		}
		ret[key] = string(kv.Value)
	}
	return ret, nil
}

func (e *Etcd) Publish(path string, val string) error {
	return e.Set(path, val)
}

func (e *Etcd) Subscribe(path string) {
	if e.root != "" {
		path = e.root + "/" + path
	}
	ch := e.cli.Watch(context.Background(), path, clientv3.WithPrefix())
	e.subscriber.OnSubscribe(strings.TrimPrefix(path, e.root+"/"))
	for resp := range ch {
		for _, ev := range resp.Events {
			if ev.Type == mvccpb.PUT {
				e.subscriber.OnSubMessage(strings.TrimPrefix(string(ev.Kv.Key), e.root+"/"), string(ev.Kv.Value))
			} else if ev.Type == mvccpb.DELETE {
				e.subscriber.OnSubMessage(strings.TrimPrefix(string(ev.Kv.Key), e.root+"/"), "")
			}
		}
	}
}

func (e *Etcd) BindSubscriber(subscriber common.ISubscriber) {
	e.subscriber = subscriber
}

func (e *Etcd) Client() *clientv3.Client {
	return e.cli
}
