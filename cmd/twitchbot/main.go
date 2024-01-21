package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/huggo-42/twitch-irc/internal/parser"
)

var addr = flag.String("addr", "irc-ws.chat.twitch.tv:80", "http service address")

func main() {
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

	var token string
	var nick string
	var channel string

	log.Println("Log in credentials")
	fmt.Print("Token: ")
	fmt.Print("\033[8m") // Hide input
	fmt.Scan(&token)
	fmt.Print("\033[28m") // Show input

	fmt.Print("App name: ")
	fmt.Scan(&nick)

	fmt.Print("Twitch channel: ")
	fmt.Scan(&channel)

	if token == "" || nick == "" || channel == "" {
		log.Fatal("Error loading .env file. Have you created one as the env.example?")
		return
	}

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
