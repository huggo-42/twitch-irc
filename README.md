# Twitch chatbot implementation

Translations:
- [Português](README_ptBR.md)

# Index
- [Building and running the chatbot](#building-and-running-the-chatbot)
- [Provide necessary information to connect](#provide-necessary-information-to-connect)
- [Commands structure](#commands-structure)
- [How to add a new commands](#to-add-a-new-command)

# Building and running the chatbot
```console
$ go build -o ./bin/chatbot cmd/twitchbot/main.go
$ ./bin/chatbot
```

# Provide necessary information to connect
- **token** (will be hidden)
- **App name**
- **Twich channel**
```console
$ ./bin/chatbot
connecting to ws://irc-ws.chat.twitch.tv:80/
Log in credentials
Token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
App name: appName
Twitch channel: twitchChannel
```

# Commands structure
### at `./internal/parser/parse.go` you will find
- **declaration** of the available commands
```go
// available commands
const (
	HELP  = "help"
	HELLO = "hello"
	BYE   = "bye"
)
```
- a **switch** for those commands
    - it appends the message you wish to return to the twitch chat to `messageToWrite`
```go
switch command {
case HELLO:
    messageToWrite += handleHelloCommand(args)
...
default:
    messageToWrite += fmt.Sprintf("command not recognized, are you trying to break me?")
}
```
- and its **handlerFunction**
```go
func handleHelloCommand(args string) string {
	messageToWrite := ""
	if args != "" {
		messageToWrite = fmt.Sprintf("Hello, %s!", args)
	} else {
		messageToWrite = fmt.Sprintf("Hello, friend!")
	}
	return messageToWrite
}
```
# To add a new command
- Declare the command name at `// available commands`
- Add a case to the switch statement
- Create it's handler function
    - **MUST** return a string (that will be the **MESSAGE SENT** in the chat)
    - can receive **(args string)** as arguments
        - text that came after the command
        - `!hello John Doe` ➜ `args == "John Doe"`
    - signature example:
        - `func handleHelloCommand(args string) string {`
