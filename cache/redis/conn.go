package redis

import (
	"Cloud-Go/config"
	"fmt"
	"github.com/garyburd/redigo/redis"

	"time"
)

var pool *redis.Pool

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			c, err := redis.Dial("tcp", config.RedisHost)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			if _, err = c.Do("AUTH", config.RedisPass); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil

		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
