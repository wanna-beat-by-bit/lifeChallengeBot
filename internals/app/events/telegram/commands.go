package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	tgClient "github.com/wanna-beat-by-bit/lifeChallengeBot/internals/app/clients/telegram"
	ce "github.com/wanna-beat-by-bit/lifeChallengeBot/internals/pkg/customError"
	store "github.com/wanna-beat-by-bit/lifeChallengeBot/internals/pkg/storage"

	"github.com/google/uuid"
)

// commands
const (
	StartCmd  = "/start"
	TestCmd   = "проверка"
	ActiveCmd = "активные задания"
)

var (
	ErrNoEventString = errors.New("event string not found in message")
	ErrNoTimeString  = errors.New("time string not found in message")
)

type Sender func(string) error

func NewSender(chatId int, client *tgClient.Client) Sender {
	return func(text string) error {
		return client.SendMessage(text, chatId)
	}
}

func (p *Processor) doCmd(chatId int, username string, text string) (err error) {
	const op = "events.telegram.doCmd"
	defer func() { err = ce.WrapIfError(op, err) }()

	log.Printf("[INF]: Received CHATID:'%d', NICK:'%s', MESSAGE'%s'", chatId, username, text)
	sendMessage := NewSender(chatId, p.tg)

	text = strings.TrimSpace(text)

	if ok := isMission(text); ok {
		if err := p.createMission(text, sendMessage); err != nil {
			return err
		}
	}

	switch text {
	case StartCmd:
		sendMessage(StartMessage)
		//log.Printf("[INF]: send '%d', '%s', '%s'", chatId, username, "message")
	case TestCmd:
		sendMessage(TestMessage)
		//log.Printf("[INF]: send '%d', '%s', '%s'", chatId, username, message)
	case ActiveCmd:
		if err := p.activeCmd(sendMessage); err != nil {
			return err
		}
		//log.Printf("[INF]: send '%d', '%s', '%s'", chatId, username, NoTaskAvailable)
	}

	return nil
}

func (p *Processor) createMission(text string, sendMessage Sender) error {
	const op = "events.telegram.commands.createMission"

	mission, err := parseEventAndTime(text)
	if err != nil {
		return ce.Wrap(op, err)
	}
	if mission.Deadline.Before(time.Now()) {
		sendMessage(IncorrectDeadline)
		return nil
	}
	mission.Id = uuid.New().String()

	//if err = p.storage.CreateMission(context.Background(), &mission); err != nil {
	//	return ce.Wrap(op, err)
	//}

	if err := sendMessage(
		fmt.Sprintf("Создано задание № %s 👌: '%s' [%s]",
			mission.Id,
			mission.Text,
			mission.Deadline.Format("2006-01-02 15:04:05")),
	); err != nil {
		return ce.Wrap(op, err)
	}

	return nil
}

func (p *Processor) activeCmd(sendMessage Sender) error {
	const op = "events.telegram.commands.activeCmd"

	activeMissions, err := p.storage.ReadAll(context.Background())
	if err != nil {
		return ce.Wrap(op, err)
	}

	if len(activeMissions) == 0 {
		sendMessage(NoTaskAvailable)
		return nil
	}
	var result string
	for _, mission := range activeMissions {
		result += fmt.Sprintf("%s: %s [%s]\n", mission.Id, mission.Text, mission.Deadline.Format("2006-01-02 15:04:05"))
	}
	sendMessage(result)

	return nil
}

func parseEventAndTime(input string) (store.Mission, error) {

	eventStart := strings.Index(input, "event=")
	if eventStart == -1 {
		return store.Mission{}, ErrNoEventString
	}
	eventStart += len("event=")

	timeStart := strings.Index(input, "time=")
	if timeStart == -1 {
		return store.Mission{}, ErrNoTimeString
	}
	timeStart += len("time=")

	eventEnd := timeStart - len("time=")
	timeStr := input[timeStart:]

	// Extract event text
	event := strings.TrimSpace(input[eventStart:eventEnd])

	// Parse time
	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return store.Mission{}, errors.New("can't parse time from message's time string")
	}

	return store.Mission{Text: event, Deadline: parsedTime}, nil
}

func isMission(input string) bool {
	if !strings.Contains(input, "event=") {
		return false
	}
	if !strings.Contains(input, "time=") {
		return false
	}
	return true

}
