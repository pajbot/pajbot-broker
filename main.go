package main

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot-broker/irc"
	"github.com/pajlada/pajbot-broker/pajbot"
)

func initRedis() *redis.Pool {
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

func publish(redisPool *redis.Pool, messageType, channel, message string) {
	c := redisPool.Get()
	c.Do("PUBLISH", messageType+"."+channel, message)
}

func main() {
	redisPool := initRedis()

	ircConnection, err := irc.Connect("pajlada")
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Print(ircConnection)

	for {
		select {
		case msg := <-ircConnection.ReadRaw:
			pMsg := pajbot.RawIRCMessage{
				Message: msg,
			}
			publish(redisPool, "twitch.raw", ircConnection.Channel, pMsg.String())
			log.Print(msg)
		case emotes := <-ircConnection.ReadEmotes:
			pMsg := pajbot.EmotesMessage{}
			publish(redisPool, "twitch.emotes", ircConnection.Channel, pMsg.String())
			log.Print(emotes)
		}
	}
}
