package main

import (
	"log"

	"github.com/pajlada/pajbot-broker/common"
	"github.com/pajlada/pajbot-broker/irc"
	"github.com/pajlada/pajbot-broker/pajbot"
)

func main() {
	redisPool := common.InitRedis()

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
			common.Publish(redisPool, "twitch.raw", ircConnection.Channel, pMsg.String())
			log.Print(msg)
		case emotes := <-ircConnection.ReadEmotes:
			pMsg := pajbot.EmotesMessage{}
			common.Publish(redisPool, "twitch.emotes", ircConnection.Channel, pMsg.String())
			log.Print(emotes)
		}
	}
}
