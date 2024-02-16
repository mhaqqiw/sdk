package qredis

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/mhaqqiw/sdk/go/qentity"
)

var Prefix string
var Conn *redis.Client
var debug bool

func CreateConn(init qentity.Redis, useDebug bool) *redis.Client {
	redisConn := redis.NewClient(&redis.Options{
		Addr:     init.Host + ":" + init.Port,
		Password: init.Password,
		DB:       init.DB,
	})
	Conn = redisConn
	debug = useDebug
	return redisConn
}

func Set(module, key string, data any, timeout int64) error {
	key = concatKey(Prefix, module, key)
	ttl := time.Duration(timeout) * time.Second
	var dataTxt string
	if debug {
		fmt.Println("key", key)
		fmt.Println("ttl", ttl)
		fmt.Println("data", data)
	}
	if fmt.Sprintf("%T", data) == "string" {
		dataTxt = fmt.Sprintf("%v", data)
	} else {
		res, err := json.Marshal(data)
		if err != nil {
			return errors.New("failed to marshal data")
		}
		dataTxt = string(res)
	}
	op1 := Conn.Set(key, dataTxt, ttl)
	if err := op1.Err(); err != nil {
		return errors.New("failed to set data")
	}
	return nil
}

func Get(module, key string) (string, time.Duration, error) {
	key = concatKey(Prefix, module, key)
	if debug {
		fmt.Println("key", key)
	}
	val, err := Conn.Get(key).Result()
	if err != nil {
		//if key not found return empty string
		if err == redis.Nil {
			return "", 0, nil
		}
		return val, 0, err
	}
	ttlResult := Conn.TTL(key)
	if ttlResult.Err() != nil {
		fmt.Println("Error:", ttlResult.Err())
		return val, 0, errors.New("failed to get ttl")
	}

	return val, ttlResult.Val(), nil
}

func GetKeys(pattern string, overridePrefix bool) ([]string, error) {
	var res []string
	key := Prefix + ":" + pattern
	if overridePrefix {
		key = pattern
	}

	if debug {
		fmt.Println("key", key)
	}
	val, err := Conn.Keys(key).Result()
	if err != nil {
		return res, errors.New("failed to get data")
	}

	for _, v := range val {
		res = append(res, v)
	}
	return res, nil
}

func Del(module, key string) error {
	key = concatKey(Prefix, module, key)
	if debug {
		fmt.Println("key", key)
	}
	err := Conn.Del(key).Err()
	if err != nil {
		return errors.New("failed to delete data")
	}
	return nil
}

func concatKey(prefix, module, key string) string {
	if prefix != "" {
		prefix = prefix + ":"
	}
	if module != "" {
		module = module + ":"
	}
	return prefix + module + key
}
