package components

// an plus version for redigo
// 1.the funcs
// 2.pub/sub
// 3.conn pool for gorounte

import (
	"fmt"
	"time"

	red "github.com/garyburd/redigo/redis"
)

const (
	OPT_MAXPOOL string = "OPT_MAXPOOL"
)

type Redis struct {
	auth        string
	RedisClient *red.Pool
	pscDict     map[string]*red.PubSubConn
	subscriber  ISubscriber
	opts        *CompOptions
}

func NewRedis() *Redis {
	r := new(Redis)
	r.pscDict = make(map[string]*red.PubSubConn)
	r.opts = NewCompOptions()
	r.opts.Set(OPT_MAXPOOL, 100)
	return r
}

// ------------- com ----------------

func (r *Redis) Conn(url string, db string) error {
	// try to make a conn to redis
	c, err := red.Dial("tcp", url)
	if err != nil {
		return err
	}
	// make sure to dispose the temp conn
	defer c.Close()
	// enable the pool

	maxpool := r.opts.GetInt(OPT_MAXPOOL)
	// redis host
	REDIS_HOST := url
	// db
	REDIS_DB := db
	// build the pool
	r.RedisClient = &red.Pool{
		// set the maxidle and maxactive
		MaxIdle:     maxpool,
		MaxActive:   maxpool * 2,
		IdleTimeout: 15 * time.Second,
		Dial: func() (red.Conn, error) {
			c, err := red.Dial("tcp", REDIS_HOST)
			if err != nil {
				fmt.Println("[Redis]::create a new redis conn faild!!")
				return nil, err
			}
			if r.auth != "" {
				_, err2 := c.Do("AUTH", r.auth)
				if err2 != nil {
					fmt.Println("[Redis]::author faild!!")
					c.Close()
					return nil, err
				}
			}
			// select the db
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}

	return nil
}

func (r *Redis) SetAuthor(user string, pass string) {
	r.auth = user
	if pass != "" {
		r.auth += ":" + pass
	}
}

func (r *Redis) SetOption(key string, val interface{}) {
	r.opts.Set(key, val)
}

func (r *Redis) Close() {
	r.RedisClient.Close()
}

// -------------client--------------

func (r *Redis) GetClient() *RedisClient {
	rc := r.RedisClient.Get()
	return NewRedisClient(rc)
}

// -------------pub/sub----------------

func (r *Redis) BindSubscriber(subscriber ISubscriber) {
	r.subscriber = subscriber
}

func (r *Redis) Publish(channel string, msg string) error {
	rc := r.RedisClient.Get()
	rc.Do("PUBLISH", channel, msg)
	rc.Close()
	return nil
}

func (r *Redis) Subscribe(channel string) {
	_, exist := r.pscDict[channel]
	if !exist {
		// ps: use another conn for the pub/sub, otherwise there will be an error..
		psc := red.PubSubConn{r.RedisClient.Get()}
		err := psc.Subscribe(channel)
		if err != nil {
			if r.subscriber != nil {
				r.subscriber.OnSubError(channel, err)
			}
		}
		for {
			switch v := psc.Receive().(type) {
			case red.Message:
				if r.subscriber != nil {
					r.subscriber.OnSubMessage(v.Channel, string(v.Data))
				}
			case red.Subscription:
				//fmt.Println("%s: %s %d\n", v.Channel, v.Kind, v.Count)
				if r.subscriber != nil {
					r.subscriber.OnSubscribe(v.Channel)
				}
			case error:
				if r.subscriber != nil {
					r.subscriber.OnSubError(channel, v)
				}
				return
			}
		}
	}
}

// -------------string----------------

func (r *Redis) Get(key string) (string, error) {
	rc := r.RedisClient.Get()
	str, err := red.String(rc.Do("GET", key))
	rc.Close()
	return str, err
}

func (r *Redis) Set(key string, val string) (bool, error) {
	rc := r.RedisClient.Get()
	suc, err := red.String(rc.Do("SET", key, val))
	rc.Close()
	return suc == "OK", err
}

func (r *Redis) Exists(key string) (bool, error) {
	rc := r.RedisClient.Get()
	ret, err := red.Bool(rc.Do("EXISTS", key))
	rc.Close()
	return ret, err
}

func (r *Redis) Delete(key string) (int64, error) {
	rc := r.RedisClient.Get()
	suc, err := red.Int64(rc.Do("DEL", key))
	rc.Close()
	return suc, err
}

// -------------set----------------

func (r *Redis) SAdd(key string, member string) (bool, error) {
	rc := r.RedisClient.Get()
	suc, err := red.Bool(rc.Do("SADD", key, member))
	rc.Close()
	return suc, err
}

func (r *Redis) SMembers(key string) ([]string, error) {
	rc := r.RedisClient.Get()
	strs, err := red.Strings(rc.Do("SMEMBERS", key))
	rc.Close()
	return strs, err
}

func (r *Redis) SRem(key string, member string) (int64, error) {
	rc := r.RedisClient.Get()
	num, err := red.Int64(rc.Do("SREM", key))
	rc.Close()
	return num, err
}

// -------------zset----------------
// for the ranking, and only for setting; and the large datas for getting, maybe you can use php :P

func (r *Redis) ZAdd(key string, score float32, val string) (bool, error) {
	rc := r.RedisClient.Get()
	_, err := red.Bool(rc.Do("ZADD", key, score, val))
	rc.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Redis) ZCount(key string, start float32, end float32) (int, error) {
	rc := r.RedisClient.Get()
	ret, err := red.Int(rc.Do("ZCOUNT", key, start, end))
	rc.Close()
	return ret, err
}

func (r *Redis) ZSize(key string) (int, error) {
	rc := r.RedisClient.Get()
	ret, err := red.Int(rc.Do("ZCARD", key))
	rc.Close()
	return ret, err
}

func (r *Redis) ZRange(key string, start float32, end float32) ([]string, error) {
	rc := r.RedisClient.Get()
	ret, err := red.Strings(rc.Do("ZRANGE", key, start, end))
	rc.Close()
	return ret, err
}

// -------------hash----------------
// for the obj record

func (r *Redis) HSet(key string, hkey string, val string) (bool, error) {
	rc := r.RedisClient.Get()
	_, err := red.Bool(rc.Do("HSET", key, hkey, val))
	rc.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Redis) HSetNx(key string, hkey string, val string) (bool, error) {
	rc := r.RedisClient.Get()
	suc, err := red.Bool(rc.Do("HSETNX", key, hkey, val))
	rc.Close()
	return suc, err
}

func (r *Redis) HGet(key string, hkey string) (string, error) {
	rc := r.RedisClient.Get()
	str, err := red.String(rc.Do("HGET", key, hkey))
	rc.Close()
	return str, err
}

func (r *Redis) HDel(key string, hkey string) (bool, error) {
	rc := r.RedisClient.Get()
	suc, err := red.Bool(rc.Do("HDEL", key, hkey))
	rc.Close()
	return suc, err
}

func (r *Redis) HLen(key string) (int, error) {
	rc := r.RedisClient.Get()
	length, err := red.Int(rc.Do("HLEN", key))
	rc.Close()
	return length, err
}

func (r *Redis) HKeys(key string) ([]string, error) {
	rc := r.RedisClient.Get()
	dict, err := red.Strings(rc.Do("HKEYS", key))
	rc.Close()
	return dict, err
}

func (r *Redis) KVals(key string) ([]string, error) {
	rc := r.RedisClient.Get()
	dict, err := red.Strings(rc.Do("HVALS", key))
	rc.Close()
	return dict, err
}

func (r *Redis) HGetAll(key string) (map[string]string, error) {
	rc := r.RedisClient.Get()
	dict, err := red.StringMap(rc.Do("HGETALL", key))
	rc.Close()
	return dict, err
}

func (r *Redis) HExists(key string, hkey string) (bool, error) {
	rc := r.RedisClient.Get()
	ret, err := red.Bool(rc.Do("HEXISTS", key, hkey))
	rc.Close()
	return ret, err
}

func (r *Redis) HMSet(key string, dict map[string]string) (bool, error) {
	rc := r.RedisClient.Get()
	suc, err := red.String(rc.Do("HMSET", red.Args{}.Add(key).AddFlat(dict)...))
	rc.Close()
	return suc == "OK", err
}
