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
	StartCmd  = "/realStart"
	TestCmd   = "проверка"
	ActiveCmd = "активные задания"
)

func NewSender(chatId int, client *tgClient.Client) func(string) error {
	return func(text string) error {
		return client.SendMessage(text, chatId)
	}
}

func (p *Processor) doCmd(chatId int, username string, text string) (err error) {
	defer func() { err = ce.WrapIfError("Error while processing commands", err) }()

	text = strings.TrimSpace(text)

	log.Printf("[INF]: Received '%d', '%s', '%s'", chatId, username, text)

	sendMessage := NewSender(chatId, p.tg)

	if ok := isMission(text); ok {
		mission, err := parseEventAndTime(text)
		if err != nil {
			return err
		}
		mission.Id = uuid.New().String()

		p.storage.CreateMission(context.Background(), &mission)

		sendMessage(
			fmt.Sprintf("Создано задание № %s 👌: '%s' [%s]",
				mission.Id,
				mission.Text,
				mission.Deadline.Format("2006-01-02 15:04:05")),
		)
	}

	switch text {
	case StartCmd:
		sendMessage("Hello! Here i am testing bot!")
		log.Printf("[INF]: send '%d', '%s', '%s'", chatId, username, "message")
	case TestCmd:
		message := "Тестирование бота"
		sendMessage(message)
		log.Printf("[INF]: send '%d', '%s', '%s'", chatId, username, message)
	case ActiveCmd:
		activeMissions, err := p.storage.ReadAll(context.Background())
		if err != nil {
			return err
		}
		if len(activeMissions) == 0 {
			message := "Пока нет заданий"
			sendMessage(message)
			log.Printf("[INF]: send '%d', '%s', '%s'", chatId, username, message)
			return nil
		}
		var result string
		for _, mission := range activeMissions {
			result += fmt.Sprintf("%s: %s [%s]\n", mission.Id, mission.Text, mission.Deadline.Format("2006-01-02 15:04:05"))
		}
		sendMessage(result)
		log.Printf("[INF]: send '%d', '%s', '%s'", chatId, username, "missions")

	}

	return nil
}

func parseEventAndTime(input string) (store.Mission, error) {

	eventStart := strings.Index(input, "event=")
	if eventStart == -1 {
		return store.Mission{}, errors.New("event string not found in message")
	}
	eventStart += len("event=")

	timeStart := strings.Index(input, "time=")
	if timeStart == -1 {
		return store.Mission{}, errors.New("time string not found in message")
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
