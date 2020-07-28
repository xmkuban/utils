package db

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func RedisTestInit() (r *RedisWrap, err error) {
	redisDbKey := "test"
	InitRedis(redisDbKey, &RedisWrap{
		Host:        "127.0.0.1:6379",
		Db:          0,
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: 10 * time.Second,
		WaitIdle:    true,
	})
	r = R(redisDbKey)
	if r == nil {
		return nil, errors.New("connect to test db fail")
	}
	return r, err
}

func TestRedisGetConn(t *testing.T) {
	r, err := RedisTestInit()
	if err != nil {
		t.Error(err)
		return
	}

	conn := r.GetConn()
	if conn != nil {
		t.Log("redis get conn success")
	} else {
		t.Error("redis get conn nil")
	}
	conn.Close()
}

func TestRedisSet(t *testing.T) {
	r, err := RedisTestInit()
	if err != nil {
		t.Error(err)
		return
	}

	err = r.Set("test", "test")
	if err == nil {
		t.Log("redis set success")
	} else {
		t.Error("redis set fail")
	}

}

func TestRedisGet(t *testing.T) {
	r, err := RedisTestInit()
	if err != nil {
		t.Error(err)
		return
	}

	res, err := r.Get("test")
	if err != nil {
		t.Error("redis get fail")
		return
	}
	if res == "test" {
		t.Log("redis get success")
	} else {
		t.Error("redis get not match with set")
	}
}

func TestRedisDelete(t *testing.T) {
	r, err := RedisTestInit()
	if err != nil {
		t.Error(err)
		return
	}

	err = r.Del("test")
	if err != nil {
		t.Error("redis delete fail")
	} else {
		t.Log("redis delete success")
	}
}

func TestPub(t *testing.T) {
	r, err := RedisTestInit()
	if err != nil {
		t.Error(err)
		return
	}
	err = r.Publish("test", "fuck")
	if err != nil {
		t.Errorf("publish fail:%s", err)
		return
	}

}
func TestSub(t *testing.T) {
	r, err := RedisTestInit()
	if err != nil {
		t.Error(err)
		return
	}

	psc := r.GetSubConn()
	defer psc.Close()
	psc.Subscribe("test")

	defer psc.Close()
	for {
		switch n := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("Message: %s %s\n", n.Channel, n.Data)
		case redis.PMessage:
			fmt.Printf("PMessage: %s %s %s\n", n.Pattern, n.Channel, n.Data)
		case redis.Subscription:
			fmt.Printf("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count)
		case error:
			fmt.Printf("error: %v\n", n)
			return
		}
	}

}
