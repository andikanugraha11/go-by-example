package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type user struct {
	FirstName string `redis:"first_name"`
	LastName  string `redis:"last_name"`
	BirthYear int    `redis:"birth_year"`
	Gender    string `redis:"gender"`
}

// constant of time
const TimeFiveMinutes = time.Second * 300

// redis connection for main package variable
var c redis.Conn

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
	c = pool.Get()

	// grcefully close redis connection
	defer c.Close()

	// test ping to redis
	pong, err := c.Do("PING")
	if err != nil {
		log.Fatalf("fail ping to redis: %s\n", err.Error())
	}
	log.Println(pong)

	hashOperation()
	setOperation()

}

// HSET redis for store and get Hash datatype on redis
// doc for HSET: https://redis.io/commands/hset
func hashOperation() {
	fmt.Println("Hash Operation")
	// store data using HSET command
	_, err := c.Do("HSET", "user:1", "first_name", "andika", "last_name", "nugraha", "birth_year", 1995, "gender", "male")
	if err != nil {
		log.Printf("fail store data to redis: %s\n", err.Error())
	}

	// get single hash value to string using HGET
	firstName, err := redis.String(c.Do("HGET", "user:1", "first_name"))
	if err != nil {
		log.Printf("fail get data from redis: %s\n", err.Error())
	}
	fmt.Println(firstName)

	// get all hash value to map using HGETALL
	userMap, err := redis.StringMap(c.Do("HGETALL", "user:1"))
	if err != nil {
		log.Printf("fail get data from redis: %s\n", err.Error())
	}
	fmt.Println(userMap)

	// get all hash value to struct using HGETALL
	var userStruct user
	res, err := redis.Values(c.Do("HGETALL", "user:1"))
	if err != nil {
		log.Printf("fail get data from redis: %s\n", err.Error())
	}
	err = redis.ScanStruct(res, &userStruct)
	if err != nil {
		log.Printf("fail to scan value into struct: %s\n", err.Error())
	}
	fmt.Println(userStruct)
	fmt.Printf("%+v\n\n", userStruct)
}

// SET redis for store and get key-value datatype on redis
// doc for SET: https://redis.io/commands/set
func setOperation() {
	fmt.Println("Set Operation")

	// store data using SET command
	_, err := c.Do("SET", "instagram_username", "andikaaa.nugraha")
	if err != nil {
		log.Printf("fail store data to redis: %s\n", err.Error())
	}

	// get data using GET command
	data, err := redis.String(c.Do("GET", "instagram_username"))
	if err != nil {
		log.Printf("fail get data from redis: %s\n", err.Error())
	}
	fmt.Println(data)

	// set with struct data
	userStruct := user{
		FirstName: "andika",
		LastName:  "nugraha",
		BirthYear: 1995,
		Gender:    "male",
	}

	// marshaling data into []byte
	userBytes, err := json.Marshal(userStruct)
	if err != nil {
		log.Printf("fail to marshaling data: %s\n", err.Error())
	}

	// store data using SET
	_, err = c.Do("SET", "user:1", userBytes)
	if err != nil {
		log.Printf("fail store data to redis: %s\n", err.Error())
	}

	// get data using GET commnand
	var resultUserStruct user
	dataStructByte, err := redis.String(c.Do("GET", "user:1"))
	if err != nil {
		log.Printf("fail get data from redis: %s\n", err.Error())
	}
	err = json.Unmarshal([]byte(dataStructByte), &resultUserStruct)
	if err != nil {
		log.Printf("fail to unmarshaling data: %s\n", err.Error())
	}
	fmt.Println(resultUserStruct)
	fmt.Printf("%+v\n\n", resultUserStruct)
}
