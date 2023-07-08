package telegram

import (
	"log"
	"strings"
)

// commands
const (
	StartCmd = "/start"
	TestCmd  = "/test"
)

func (p *Processor) doCmd(chatId int, username string, text string) error {
	text = strings.TrimSpace(text)

	log.Printf("[INF]: '%d', '%s', '%s'", chatId, username, text)
	switch text {
	case StartCmd:
		p.tg.SendMessage("Hello! Here i am testing bot!", chatId)
		log.Println("[INF]: Sent message")
	case TestCmd:
		p.tg.SendMessage("You are testing me!", chatId)
	}

	return nil
}
