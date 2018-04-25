package redis

// an plus version for redigo
// 1.the funcs
// 2.pub/sub
// 3.conn pool for gorounte

import (
	"fmt"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	RedisClient *redis.Pool
	Psc         redis.PubSubConn
	lock        sync.RWMutex
}

// ------------- com ----------------

func (this *Redis) Conn(url string, db int, maxpool int, auth string) error {
	// try to make a conn to redis
	c, err := redis.Dial("tcp", url)
	if err != nil {
		return err
	}
	// make sure to dispose the temp conn
	defer c.Close()
	// enable the pool

	// redis host
	REDIS_HOST := url
	// db
	REDIS_DB := db
	// build the pool
	this.RedisClient = &redis.Pool{
		// set the maxidle and maxactive
		MaxIdle:     maxpool,
		MaxActive:   maxpool * 2,
		IdleTimeout: 15 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				fmt.Println("[Redis]::create a new redis conn faild!!")
				return nil, err
			}
			if auth != "" {
				_, err2 := c.Do("AUTH", auth)
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

	// ps: use another conn for the pub/sub, otherwise there will be an error..
	this.Psc = redis.PubSubConn{this.RedisClient.Get()}

	return nil
}

func (this *Redis) Close() {
	this.RedisClient.Close()
}

// -------------string----------------

func (this *Redis) Get(key string) (string, error) {
	rc := this.RedisClient.Get()
	str, err := redis.String(rc.Do("GET", key))
	rc.Close()
	return str, err
}

func (this *Redis) Set(key string, val string) (bool, error) {
	rc := this.RedisClient.Get()
	suc, err := redis.String(rc.Do("SET", key, val))
	rc.Close()
	return suc == "OK", err
}

func (this *Redis) Exists(key string) (bool, error) {
	rc := this.RedisClient.Get()
	ret, err := redis.Bool(rc.Do("EXISTS", key))
	rc.Close()
	return ret, err
}

func (this *Redis) Delete(key string) (int64, error) {
	rc := this.RedisClient.Get()
	suc, err := redis.Int64(rc.Do("DEL", key))
	rc.Close()
	return suc, err
}

// -------------pub/sub----------------

func (this *Redis) Publish(channel string, msg string) {
	rc := this.RedisClient.Get()
	rc.Do("PUBLISH", channel, msg)
	rc.Close()
}

func (this *Redis) Subscribe(channel string) {
	this.Psc.Subscribe(channel)
}

// -------------set----------------

func (this *Redis) SAdd(key string, member string) (bool, error) {
	rc := this.RedisClient.Get()
	suc, err := redis.Bool(rc.Do("SADD", key, member))
	rc.Close()
	return suc, err
}

func (this *Redis) SMembers(key string) ([]string, error) {
	rc := this.RedisClient.Get()
	strs, err := redis.Strings(rc.Do("SMEMBERS", key))
	rc.Close()
	return strs, err
}

// -------------zset----------------
// for the ranking, and only for setting; and the large datas for getting, maybe you can use php :P

func (this *Redis) ZAdd(key string, score float32, val string) (bool, error) {
	rc := this.RedisClient.Get()
	_, err := redis.Bool(rc.Do("ZADD", key, score, val))
	rc.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *Redis) ZCount(key string, start float32, end float32) (int, error) {
	rc := this.RedisClient.Get()
	ret, err := redis.Int(rc.Do("ZCOUNT", key, start, end))
	rc.Close()
	return ret, err
}

func (this *Redis) ZSize(key string) (int, error) {
	rc := this.RedisClient.Get()
	ret, err := redis.Int(rc.Do("ZCARD", key))
	rc.Close()
	return ret, err
}

func (this *Redis) ZRange(key string, start float32, end float32) ([]string, error) {
	rc := this.RedisClient.Get()
	ret, err := redis.Strings(rc.Do("ZRANGE", key, start, end))
	rc.Close()
	return ret, err
}

// -------------hash----------------
// for the obj record

func (this *Redis) HSet(key string, hkey string, val string) (bool, error) {
	rc := this.RedisClient.Get()
	_, err := redis.Bool(rc.Do("HSET", key, hkey, val))
	rc.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *Redis) HSetNx(key string, hkey string, val string) (bool, error) {
	rc := this.RedisClient.Get()
	suc, err := redis.Bool(rc.Do("HSETNX", key, hkey, val))
	rc.Close()
	return suc, err
}

func (this *Redis) HGet(key string, hkey string) (string, error) {
	rc := this.RedisClient.Get()
	str, err := redis.String(rc.Do("HGET", key, hkey))
	rc.Close()
	return str, err
}

func (this *Redis) HDel(key string, hkey string) (bool, error) {
	rc := this.RedisClient.Get()
	suc, err := redis.Bool(rc.Do("HDEL", key, hkey))
	rc.Close()
	return suc, err
}

func (this *Redis) HLen(key string) (int, error) {
	rc := this.RedisClient.Get()
	length, err := redis.Int(rc.Do("HLEN", key))
	rc.Close()
	return length, err
}

func (this *Redis) HGetAll(key string) (map[string]string, error) {
	rc := this.RedisClient.Get()
	dict, err := redis.StringMap(rc.Do("HGETALL", key))
	rc.Close()
	return dict, err
}

func (this *Redis) HExists(key string, hkey string) (bool, error) {
	rc := this.RedisClient.Get()
	ret, err := redis.Bool(rc.Do("HEXISTS", key, hkey))
	rc.Close()
	return ret, err
}

func (this *Redis) HMSet(key string, dict map[string]string) (bool, error) {
	rc := this.RedisClient.Get()
	suc, err := redis.String(rc.Do("HMSET", redis.Args{}.Add(key).AddFlat(dict)...))
	rc.Close()
	return suc == "OK", err
}
