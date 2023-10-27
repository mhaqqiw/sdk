package qredis

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/mhaqqiw/sdk/go/qentity"
)

func CreateConn(init qentity.Redis) *redis.Client {
	redisConn := redis.NewClient(&redis.Options{
		Addr:     init.Host + ":" + init.Port,
		Password: init.Password,
		DB:       0,
	})
	return redisConn
}

func Set(conn *redis.Client, key string, data any, timeout int64) error {
	ttl := time.Duration(timeout) * time.Second
	var dataTxt string
	if fmt.Sprintf("%T", data) == "string" {
		dataTxt = fmt.Sprintf("%v", data)
	} else {
		res, err := json.Marshal(data)
		if err != nil {
			return errors.New("failed to marshal data")
		}
		dataTxt = string(res)
	}
	op1 := conn.Set(key, dataTxt, ttl)
	if err := op1.Err(); err != nil {
		return errors.New("failed to set data")
	}
	return nil
}

func Get(conn *redis.Client, key string) (string, error) {
	val, err := conn.Get(key).Result()
	if err != nil {
		return val, errors.New("failed to get data")
	}
	return val, nil
}

func Del(conn *redis.Client, key string) error {
	err := conn.Del(key).Err()
	if err != nil {
		return errors.New("failed to delete data")
	}
	return nil
}
