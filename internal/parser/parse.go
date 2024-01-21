package parser

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

// available commands
const (
	HELP  = "help"
	HELLO = "hello"
	BYE   = "bye"
)

// twitch message types
const (
	JOIN    = "JOIN"
	NICK    = "NICK"
	NOTICE  = "NOTICE"
	PART    = "PART"
	PASS    = "PASS"
	PING    = "PING"
	PONG    = "PONG"
	PRIVMSG = "PRIVMSG"
)

type Message struct {
	Identity string
	Type     string
	Channel  string
	Message  string
}

func isCommand(str string) bool {
	return str[0] == '!'
}

func isPing(str string) bool {
	return strings.Split(str, " ")[0] == PING
}

func handleHelpCommand() string {
	return fmt.Sprintf("Available commands: !%s, !%s, !%s", HELP, HELLO, BYE)
}

func handleHelloCommand(args string) string {
	messageToWrite := ""
	if args != "" {
		messageToWrite = fmt.Sprintf("Hello, %s!", args)
	} else {
		messageToWrite = fmt.Sprintf("Hello, friend!")
	}
	return messageToWrite
}

func handleByeCommand(args string) string {
	messageToWrite := ""
	if args != "" {
		messageToWrite = fmt.Sprintf("Cya, %s!", args)
	} else {
		messageToWrite = fmt.Sprintf("Cya, friend!")
	}
	return messageToWrite
}

func ParseTwitchMessage(c *websocket.Conn, str string) {
	if isPing(str) {
		_, mess, _ := strings.Cut(str, ":")
		fmt.Println(PING + " " + PONG)
		c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s %s", PONG, mess)))
		return
	}
	_, identityTypeChannel, _ := strings.Cut(str, ":")
	strSplitted := strings.Split(identityTypeChannel, " ")
	id := strSplitted[0]
	messType := strSplitted[1]
	channel := strSplitted[2]
	if messType != PRIVMSG {
		switch messType {
		// For future implementations
		default:
			fmt.Println("MessageType not implemented: " + messType)
			return
		}
	}
	_, mess, _ := strings.Cut(identityTypeChannel, ":")
	msg := Message{
		Identity: id,
		Type:     messType,
		Channel:  channel,
		Message:  mess,
	}
	if isCommand(msg.Message) == false {
		return
	}
	command := strings.TrimSpace(strings.Split(msg.Message, " ")[0][1:])
	_, args, _ := strings.Cut(msg.Message, " ")
	args = strings.TrimSpace(args)
	messageToWrite := fmt.Sprintf("%s %s :", PRIVMSG, channel)
	switch command {
	case HELP:
		messageToWrite += handleHelpCommand()
	case HELLO:
		messageToWrite += handleHelloCommand(args)
	case BYE:
		messageToWrite += handleByeCommand(args)
	default:
		messageToWrite += fmt.Sprintf("command not recognized, are you trying to break me?")
	}
	fmt.Println("Command: " + command + " -> Writting to chat: " + messageToWrite)
	c.WriteMessage(websocket.TextMessage, []byte(messageToWrite))
}
