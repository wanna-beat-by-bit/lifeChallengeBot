package telegram

import (
	"errors"
	tgClient "tgBot/clients/telegram"
	ev "tgBot/events"
	ce "tgBot/pkg/customError"
)

type Processor struct {
	tg     *tgClient.Client
	offset int
}

type Meta struct {
	chatId   int
	username string
}

var (
	ErrorUnknownMeta  = errors.New("Error while validating meta data")
	ErrorUnknownEvent = errors.New("Error while validatig event type")
)

func NewProcessor(tg *tgClient.Client) *Processor {
	return &Processor{
		tg:     tg,
		offset: 0,
	}
}

func (p *Processor) Fetch(limit int) (evs []ev.Event, err error) {
	defer func() { err = ce.WrapIfError("Error occured while fetching events", err) }()

	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, err
	}

	if len(updates) == 0 {
		return nil, nil
	}

	events := make([]ev.Event, len(updates))

	for _, update := range updates {
		events = append(events, updateToEvent(update))
	}

	p.offset = updates[len(updates)-1].Id + 1

	return events, nil
}

func (p *Processor) Process(event ev.Event) error {
	switch event.Type {
	case ev.Message:
		p.processMessage(event)
	default:
		return ErrorUnknownEvent
	}

	return nil
}

func (p *Processor) processMessage(event ev.Event) (err error) {
	defer func() { err = ce.WrapIfError("Error while processing message", err) }()

	meta, err := takeMeta(event)
	if err != nil {
		return err
	}

	if err := p.doCmd(meta.chatId, meta.username, event.Text); err != nil {
		return err
	}

	return nil
}

func takeMeta(event ev.Event) (Meta, error) {
	meta, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrorUnknownMeta
	}

	return meta, nil
}

func updateToEvent(update tgClient.Update) ev.Event {
	updateType := checkType(update)

	event := ev.Event{
		Type: updateType,
		Text: fetchMessageText(update),
	}

	if updateType == ev.Message {
		event.Meta = Meta{
			chatId:   update.Message.Chat.Id,
			username: update.Message.From.Username,
		}
	}

	return event
}

func checkType(update tgClient.Update) ev.Type {
	if update.Message != nil {
		return ev.Message
	}

	return ev.Unknown
}

func fetchMessageText(update tgClient.Update) string {
	if update.Message != nil {
		return update.Message.Text
	}

	return ""
}
