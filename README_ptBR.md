# Implementação de um chatbot para Twitch

# Índice
- [Configurar .env](#configurar-env)
- [Construir e rodar o chatbot](#construir-e-rodar-o-chatbot)
- [Estrutura dos comandos](#estrutura-dos-comandos)
- [Como adicionar novos comandos](#como-adicionar-novos-comandos)

# Configurar .env
- Copie `env.example` para `.env`
```console
$ cp env.example .env
```
- Modifique `.env` de acordo com seu token, nome do bot e canal
```
TOKEN=token
NICK=botname
CHANNEL=channel
```

# Construir e rodar o chatbot
```console
$ go build -o ./bin/chatbot cmd/twitchbot/main.go && ./bin/chatbot
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
