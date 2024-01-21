package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/huggo-42/twitch-irc/internal/parser"
	"github.com/joho/godotenv"
)

var addr = flag.String("addr", "irc-ws.chat.twitch.tv:80", "http service address")

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Have you created one?")
		return
	}

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	token := os.Getenv("TOKEN")
	nick := os.Getenv("NICK")
	channel := os.Getenv("CHANNEL")
	if token == "" || nick == "" || channel == "" {
		log.Fatal("Error loading .env file. Have you created one as the env.example?")
		return
	}
	log.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
	log.Println("Logging in using the following credentials")
	log.Println("token: ", token)
	log.Println("nick: ", nick)
	log.Println("channel: ", channel)
	log.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")

	c.WriteMessage(websocket.TextMessage, []byte(parser.PASS+" oauth:"+token))
	c.WriteMessage(websocket.TextMessage, []byte(parser.NICK+" "+nick))
	c.WriteMessage(websocket.TextMessage, []byte(parser.JOIN+" #"+channel))
	c.WriteMessage(websocket.TextMessage, []byte(parser.PRIVMSG+" #"+channel+" :The bot is ON!"))

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Printf("Twitch message received at %s -> %s", time.Now().Format(time.RFC3339), message)

			parser.ParseTwitchMessage(c, string(message))
			log.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return
		}
	}
}
