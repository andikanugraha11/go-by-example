package main

// docs: https://redis.io/docs/reference/patterns/distributed-locks/

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	redigolib "github.com/gomodule/redigo/redis"
)

// constant of time
const TimeFiveMinutes = time.Second * 300

// create redis pool connection
func newRedisPool(host, password string) redsyncredis.Pool {
	pool := redigo.NewPool(&redigolib.Pool{
		MaxIdle:     5,
		IdleTimeout: TimeFiveMinutes,
		Dial: func() (redigolib.Conn, error) {
			conn, err := redigolib.Dial("tcp", host)
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
	})

	return pool
}

func main() {
	// init redis pool connection
	pool := newRedisPool("localhost:6379", "my_pass")

	ctx := context.Background()
	c, err := pool.Get(ctx)
	if err != nil {
		log.Panicf("Failed to get connection redis: %s", err.Error())
	}
	defer c.Close()
	rs := redsync.New(pool)

	mutex := rs.NewMutex("test-redsync")

	// lock
	if err = mutex.Lock(); err != nil {
		log.Panic(err)
	}

	// store data using SET command
	_, err = c.Set("instagram_username", "andikaaa.nugraha")
	if err != nil {
		log.Printf("fail store data to redis: %s\n", err.Error())
	}

	// get data using GET command
	data, err := c.Get("instagram_username")
	if err != nil {
		log.Printf("fail get data from redis: %s\n", err.Error())
	}
	fmt.Println(data)

	// release lock
	if _, err = mutex.Unlock(); err != nil {
		log.Panic(err)
	}
}
