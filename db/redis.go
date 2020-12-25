package db

import (
	"time"

	"errors"
	"reflect"

	"github.com/gomodule/redigo/redis"
)

type RedisWrap struct {
	Host     string
	Password string
	Db       int
	Pool     *redis.Pool
}

func NewDefaultRedis(host string, password string, db int) (*RedisWrap, error) {
	return NewRedis(host, password, db, 0, 0, 0, false)
}

func NewRedis(host string, password string, db int, maxIdle int, maxActive int, idleTimeout time.Duration, waitIdle bool) (*RedisWrap, error) {
	r := &RedisWrap{
		Host:     host,
		Password: password,
		Db:       db,
	}
	r.Pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Wait:        waitIdle,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", r.Host)
			if err != nil {
				return nil, err
			}
			if r.Password != "" {
				if _, err := c.Do("AUTH", r.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if r.Db > 0 {
				if _, err := c.Do("SELECT", r.Db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return err
			}
			return nil
		},
	}
	return r, nil
}

func (r *RedisWrap) GetConn() redis.Conn {
	if r.Pool == nil {
		return nil
	}
	conn := r.Pool.Get()
	return conn
}

func (r *RedisWrap) Get(key string) (string, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.String(conn.Do("get", key))
}

func (r *RedisWrap) GetInt(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.Int64(conn.Do("get", key))
}

func (r *RedisWrap) Set(key string, val string) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	return err
}

func (r *RedisWrap) Lpush(key string, val []byte) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("LPUSH", key, val)
	return err
}

func (r *RedisWrap) BRPOP(key string, timeout int) ([]byte, error) {
	conn := r.GetConn()
	defer conn.Close()
	reply, err := conn.Do("BRPOP", key, timeout)
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, redis.ErrNil
	}
	replyType := reflect.TypeOf(reply).String()
	if replyType == "[]interface {}" {
		_reply := reply.([]interface{})

		for k, v := range _reply {
			if k == 0 {
				continue
			}
			if reflect.TypeOf(v).String() == "[]uint8" {
				return v.([]byte), nil
			}
		}
	} else {
		return nil, errors.New("data type is " + replyType)
	}
	return nil, redis.ErrNil
}

func (r *RedisWrap) RPOP(key string) ([]byte, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.Bytes(conn.Do("RPOP", key))
}

func (r *RedisWrap) Setex(key string, val string, time int64) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, time, val)
	return err
}

func (r *RedisWrap) SetInt(key string, val int) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	return err
}

func (r *RedisWrap) Del(key string) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (r *RedisWrap) Incr(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.Int64(conn.Do("INCR", key))
}

func (r *RedisWrap) IncrBy(key string, step int64) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Int64(conn.Do("INCRBY", key, step))
}

func (r *RedisWrap) Decr(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.Int64(conn.Do("DECR", key))
}

func (r *RedisWrap) DecrBy(key string, step int64) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Int64(conn.Do("DECRBY", key, step))
}

func (r *RedisWrap) Exists(key string) (bool, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", key))
}

func (r *RedisWrap) Expire(key string, expire int) error {
	conn := r.GetConn()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, expire)
	return err
}

func (r *RedisWrap) GetTTL(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.Int64(conn.Do("ttl", key))
}

func (r *RedisWrap) HSet(key string, field string, value string) error {
	conn := r.GetConn()
	defer conn.Close()
	_, err := conn.Do("HSET", key, field, value)
	return err
}

func (r *RedisWrap) HDel(key string, field string) error {
	conn := r.GetConn()
	defer conn.Close()
	_, err := conn.Do("HDEL", key, field)
	return err
}

func (r *RedisWrap) Hlen(key string) (int, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Int(conn.Do("HLEN", key))
}

func (r *RedisWrap) HGetAll(key string) (map[string]string, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.StringMap(conn.Do("HGETALL", key))
}

func (r *RedisWrap) HGet(key string, field string) (string, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.String(conn.Do("HGET", key, field))
}

func (r *RedisWrap) HKeys(key string) ([]string, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Strings(conn.Do("HKEYS", key))
}

func (r *RedisWrap) HVALS(key string) ([]string, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Strings(conn.Do("HVALS", key))
}

func (r *RedisWrap) HScanStringMap(key string, f func(k string, v string) error) error {
	conn := r.GetConn()
	defer conn.Close()

	iter := 0
	for {
		if arr, err := redis.MultiBulk(conn.Do("HSCAN", key, iter)); err != nil {
			return err
		} else {
			iter, _ = redis.Int(arr[0], nil)
			maps, _ := redis.StringMap(arr[1], nil)
			for k, v := range maps {
				err := f(k, v)
				if err != nil {
					return err
				}
			}
		}
		if iter == 0 {
			break
		}
	}
	return nil
}

func (r *RedisWrap) Publish(channel string, msg string) error {
	conn := r.GetConn()
	defer conn.Close()
	err := conn.Send("PUBLISH", channel, msg)
	if err != nil {
		return err
	}
	return conn.Flush()
}

func (r *RedisWrap) GetSubConn() redis.PubSubConn {
	conn := redis.PubSubConn{
		Conn: r.GetConn(),
	}
	return conn
}

func (r *RedisWrap) RPush(key string, value []byte) error {
	conn := r.GetConn()
	defer conn.Close()
	_, err := conn.Do("RPUSH", key, value)
	return err
}

func (r *RedisWrap) BLPOP(key string, timeout int) ([]byte, error) {
	conn := r.GetConn()
	defer conn.Close()
	reply, err := conn.Do("BLPOP", key, timeout)
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, redis.ErrNil
	}
	replyType := reflect.TypeOf(reply).String()
	if replyType == "[]interface {}" {
		_reply := reply.([]interface{})

		for k, v := range _reply {
			if k == 0 {
				continue
			}
			if reflect.TypeOf(v).String() == "[]uint8" {
				return v.([]byte), nil
			}
		}
	} else {
		return nil, errors.New("data type is " + replyType)
	}
	return nil, redis.ErrNil
}

func (r *RedisWrap) LPOP(key string) ([]byte, error) {
	conn := r.GetConn()
	defer conn.Close()

	return redis.Bytes(conn.Do("LPOP", key))
}

func (r *RedisWrap) LLEN(key string) (int, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Int(conn.Do("LLEN", key))
}

func (r *RedisWrap) Lindex(key string, index int) ([]byte, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Bytes(conn.Do("lindex", key, index))
}

func (r *RedisWrap) Lrange(key string, start int, end int) ([]interface{}, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Values(conn.Do("LRANGE", key, start, end))
}

func (r *RedisWrap) Sadd(key string, member string) error {
	conn := r.GetConn()
	defer conn.Close()
	_, err := conn.Do("SADD", key, member)
	return err
}

func (r *RedisWrap) Scard(key string) (int, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Int(conn.Do("SCARD", key))
}

func (r *RedisWrap) Smembers(key string) ([]string, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Strings(conn.Do("SMEMBERS", key))
}

func (r *RedisWrap) Sismember(key string, member string) (bool, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Bool(conn.Do("SISMEMBER", key, member))
}

func (r *RedisWrap) Srem(key string, member string) error {
	conn := r.GetConn()
	defer conn.Close()
	_, err := conn.Do("SREM", key, member)
	return err
}

func (r *RedisWrap) SPop(key string) (string, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.String(conn.Do("SPOP", key))
}

//给有序集合添加成员
func (r *RedisWrap) ZAdd(key string, member string, score string) error {
	conn := r.GetConn()
	defer conn.Close()
	_, err := conn.Do("ZADD", key, score, member)
	return err
}

//获取有序集合的成员数
func (r *RedisWrap) ZCard(key string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Int64(conn.Do("ZCARD", key))
}

//获取指定成员的排名
func (r *RedisWrap) ZreVRank(key string, member string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Int64(conn.Do("ZREVRANK", key, member))
}

//计算在有序集合中指定区间分数的成员数
func (r *RedisWrap) ZCount(key string, min string, max string) (int64, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Int64(conn.Do("ZCOUNT", key, min, max))
}

//获取给定区间的排名 (从高到低)
func (r *RedisWrap) ZreVRange(key string, start int, end int, withscore bool) ([]string, error) {
	conn := r.GetConn()
	defer conn.Close()
	if !withscore {
		return redis.Strings(conn.Do("ZREVRANGE", key, start, end))
	}
	return redis.Strings(conn.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
}

//获取给定区间的排名 (从低到搞)
func (r *RedisWrap) ZRange(key string, start int, end int, withscore bool) ([]string, error) {
	conn := r.GetConn()
	defer conn.Close()
	if !withscore {
		return redis.Strings(conn.Do("ZRANGE", key, start, end))
	}
	return redis.Strings(conn.Do("ZRANGE", key, start, end, "WITHSCORES"))
}

//获取自己的分数
func (r *RedisWrap) ZScore(key string, member string) (string, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.String(conn.Do("ZSCORE", key, member))
}

func (r *RedisWrap) ZinCrBy(key string, member string, increment int) (int, error) {
	conn := r.GetConn()
	defer conn.Close()
	return redis.Int(conn.Do("ZINCRBY", key, increment, member))
}

func (r *RedisWrap) Zrem(key string, member string) error {
	conn := r.GetConn()
	defer conn.Close()
	_, err := conn.Do("ZREM", key, member)
	return err
}

func (r *RedisWrap) Zrevrangebyscore(key string, max int64, min int64, withscore bool) ([]string, error) {
	conn := r.GetConn()
	defer conn.Close()
	if !withscore {
		return redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, max, min))
	}
	return redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, max, min, "WITHSCORES"))
}

func (r *RedisWrap) Zrangebyscore(key string, max int64, min int64, withscore bool) ([]string, error) {
	conn := r.GetConn()
	defer conn.Close()
	if !withscore {
		return redis.Strings(conn.Do("ZRANGEBYSCORE", key, max, min))
	}
	return redis.Strings(conn.Do("ZRANGEBYSCORE", key, max, min, "WITHSCORES"))
}
