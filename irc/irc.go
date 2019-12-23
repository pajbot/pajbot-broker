package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"

	"github.com/pajlada/pajbot2/parser"
)

type Connection struct {
	conn       net.Conn
	ReadRaw    chan string
	ReadEmotes chan string
	Connected  bool
	Channel    string
}

func Connect(channel string) (Connection, error) {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		return Connection{}, err
	}
	ircConnection := Connection{}
	ircConnection.ReadRaw = make(chan string)
	ircConnection.ReadEmotes = make(chan string)
	ircConnection.conn = conn
	ircConnection.Channel = channel

	ircConnection.send("NICK justinfan123")
	ircConnection.send("CAP REQ :twitch.tv/tags")
	ircConnection.send("CAP REQ :twitch.tv/commands")

	go ircConnection.startReading()
	return ircConnection, nil
}

func (c *Connection) send(message string) {
	fmt.Fprint(c.conn, message+"\r\n")
}

func (c *Connection) startReading() {
	reader := bufio.NewReader(c.conn)
	tp := textproto.NewReader(reader)

	for {
		line, err := tp.ReadLine()
		if err != nil {
			log.Fatal("Error reading from connection:", err)
			return
		}

		// log.Print(line)

		if strings.HasPrefix(line, ":tmi.twitch.tv 376 ") {
			// CONNECTED
			c.Connected = true

			c.send("JOIN #" + c.Channel)
		}

		if c.Connected {
			// log.Print("sending:", line)
			c.ReadRaw <- line

			heh := parser.Parse(line)
			fmt.Println("heh", heh)

			c.test()
		}
	}
}

var i = 0

func (c *Connection) test() {
	c.ReadEmotes <- "xd"
}
