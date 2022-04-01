package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// constant of time
const TimeFiveMinutes = time.Second * 300

// create redis pool connection
func newRedisPool(host, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: TimeFiveMinutes,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := conn.Do("AUTH", password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		},
	}
}

func main() {
	// init redis pool connection
	pool := newRedisPool("localhost:6379", "my_pass")

	// get redis connection
	c := pool.Get()

	// grcefully close redis connection
	defer c.Close()

	// test ping to redis
	_, err := c.Do("PING")
	if err != nil {
		log.Fatalf("fail ping to redis: %s\n", err.Error())
	}

}
