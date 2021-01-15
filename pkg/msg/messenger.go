package msg

import "fmt"

type Messanger struct{}

func (m *Messanger) SentMessage(msg string) {
	//TODO sent to messenger
	fmt.Println("Service do not sent 'info' in Telegram or WebHook if you read this message.")
	fmt.Printf("info: %s\n", msg)
}

func NewMessanger() *Messanger {
	return &Messanger{}
}
