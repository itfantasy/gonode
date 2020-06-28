package components

import (
	red "github.com/garyburd/redigo/redis"
)

type RedisClient struct {
	conn red.Conn
}

func NewRedisClient(conn red.Conn) *RedisClient {
	rc := new(RedisClient)
	rc.conn = conn
	return rc
}

func (rc *RedisClient) Dispose() {
	rc.conn.Close()
	rc = nil
}

// -------------string----------------

func (rc *RedisClient) Get(key string) (string, error) {
	str, err := red.String(rc.conn.Do("GET", key))
	return str, err
}

func (rc *RedisClient) Set(key string, val string) (bool, error) {
	suc, err := red.String(rc.conn.Do("SET", key, val))
	return suc == "OK", err
}

func (rc *RedisClient) Exists(key string) (bool, error) {
	ret, err := red.Bool(rc.conn.Do("EXISTS", key))
	return ret, err
}

func (rc *RedisClient) Delete(key string) (int64, error) {
	suc, err := red.Int64(rc.conn.Do("DEL", key))
	return suc, err
}

// -------------set----------------

func (rc *RedisClient) SAdd(key string, member string) (bool, error) {
	suc, err := red.Bool(rc.conn.Do("SADD", key, member))
	return suc, err
}

func (rc *RedisClient) SMembers(key string) ([]string, error) {
	strs, err := red.Strings(rc.conn.Do("SMEMBERS", key))
	return strs, err
}

func (rc *RedisClient) SRem(key string, member string) (int64, error) {
	num, err := red.Int64(rc.conn.Do("SREM", key))
	return num, err
}

// -------------zset----------------
// for the ranking, and only for setting; and the large datas for getting, maybe you can use php :P

func (rc *RedisClient) ZAdd(key string, score float32, val string) (bool, error) {
	_, err := red.Bool(rc.conn.Do("ZADD", key, score, val))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rc *RedisClient) ZCount(key string, start float32, end float32) (int, error) {
	ret, err := red.Int(rc.conn.Do("ZCOUNT", key, start, end))
	return ret, err
}

func (rc *RedisClient) ZSize(key string) (int, error) {
	ret, err := red.Int(rc.conn.Do("ZCARD", key))
	return ret, err
}

func (rc *RedisClient) ZRange(key string, start float32, end float32) ([]string, error) {
	ret, err := red.Strings(rc.conn.Do("ZRANGE", key, start, end))
	return ret, err
}

// -------------hash----------------
// for the obj record

func (rc *RedisClient) HSet(key string, hkey string, val string) (bool, error) {
	_, err := red.Bool(rc.conn.Do("HSET", key, hkey, val))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rc *RedisClient) HSetNx(key string, hkey string, val string) (bool, error) {
	suc, err := red.Bool(rc.conn.Do("HSETNX", key, hkey, val))
	return suc, err
}

func (rc *RedisClient) HGet(key string, hkey string) (string, error) {
	str, err := red.String(rc.conn.Do("HGET", key, hkey))
	return str, err
}

func (rc *RedisClient) HDel(key string, hkey string) (bool, error) {
	suc, err := red.Bool(rc.conn.Do("HDEL", key, hkey))
	return suc, err
}

func (rc *RedisClient) HLen(key string) (int, error) {
	length, err := red.Int(rc.conn.Do("HLEN", key))
	return length, err
}

func (rc *RedisClient) HKeys(key string) ([]string, error) {
	dict, err := red.Strings(rc.conn.Do("HKEYS", key))
	return dict, err
}

func (rc *RedisClient) KVals(key string) ([]string, error) {
	dict, err := red.Strings(rc.conn.Do("HVALS", key))
	return dict, err
}

func (rc *RedisClient) HGetAll(key string) (map[string]string, error) {
	dict, err := red.StringMap(rc.conn.Do("HGETALL", key))
	return dict, err
}

func (rc *RedisClient) HExists(key string, hkey string) (bool, error) {
	ret, err := red.Bool(rc.conn.Do("HEXISTS", key, hkey))
	return ret, err
}

func (rc *RedisClient) HMSet(key string, dict map[string]string) (bool, error) {
	suc, err := red.String(rc.conn.Do("HMSET", red.Args{}.Add(key).AddFlat(dict)...))
	return suc == "OK", err
}
