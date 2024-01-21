# Implementação de um chatbot para Twitch

# Índice
- [Construir e rodar o chatbot](#construir-e-rodar-o-chatbot)
- [Forneça as informações necessárias para o bot se conectar](#forneça-as-informações-necessárias-para-o-bot-se-conectar)
- [Estrutura dos comandos](#estrutura-dos-comandos)
- [Como adicionar novos comandos](#como-adicionar-novos-comandos)

# Construir e rodar o chatbot
```console
$ go build -o ./bin/chatbot cmd/twitchbot/main.go
$ ./bin/chatbot
```

# Forneça as informações necessárias para o bot se conectar
- **token** (estará escondido)
- **Nome do aplicativo** (App name)
- **Canal da twich** (Twitch channel)
```console
$ ./bin/chatbot
connecting to ws://irc-ws.chat.twitch.tv:80/
Log in credentials
Token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
App name: appName
Twitch channel: twitchChannel
```

# Estrutura dos comandos
### em `./internal/parser/parse.go` você acha
- **declarações** dos comandos disponíveis
```go
// available commands
const (
	HELP  = "help"
	HELLO = "hello"
	BYE   = "bye"
)
```
- um **switch** para esse comandos
    - que acrescenta a mensagem que você deseja retornar para o chat da twitch para `messageToWrite`
```go
switch command {
case HELLO:
    messageToWrite += handleHelloCommand(args)
...
default:
    messageToWrite += fmt.Sprintf("command not recognized, are you trying to break me?")
}
```
- e sua **função manipuladora (handlerFunction)**
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
# Como adicionar novos comandos
- Declare o nome do comando em `// available commands`
- Adicione um caso no switch
- Crie sua função manipuladora (handlerFunction)
    - **DEVE** retornar uma string (que será a **MENSAGEM ENVIADA** no chat)
    - pode receber **(args string)** como parâmetro
        - texto que veio depois do comando
        - `!hello John Doe` ➜ `args == "John Doe"`
    - exemplo da assinatura:
        - `func handleHelloCommand(args string) string {`
