package common

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// InitRedis connects to redis and returns a redis pool
func InitRedis() *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", "localhost:6379")
		if err != nil {
			log.Fatal("An error occured while connecting to redis:", err)
			return nil, err
		}
		_, err = c.Do("SELECT", 8)
		if err != nil {
			log.Fatal("Error while selceting redis db:", err)
			return nil, err
		}
		return c, nil
	}, 69)
}

// Publish is a helper method that simplifies publishing a message to redis pubsub
func Publish(redisPool *redis.Pool, messageType, channel, message string) {
	c := redisPool.Get()
	c.Do("PUBLISH", messageType+"."+channel, message)
}
