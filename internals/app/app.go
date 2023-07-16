package app

import (
	"context"
	"fmt"
	"log"

	tgClient "github.com/wanna-beat-by-bit/lifeChallengeBot/internals/app/clients/telegram"
	"github.com/wanna-beat-by-bit/lifeChallengeBot/internals/app/consumer"
	tgProcessor "github.com/wanna-beat-by-bit/lifeChallengeBot/internals/app/events/telegram"
	"github.com/wanna-beat-by-bit/lifeChallengeBot/internals/pkg/storage"
	"github.com/wanna-beat-by-bit/lifeChallengeBot/internals/pkg/storage/sqlite"
)

const (
	telegramHost = "api.telegram.org"
	batchSize    = 10
)

type App struct {
	token     string
	client    *tgClient.Client
	consumer  *consumer.Consumer
	processor *tgProcessor.Processor
	storage   storage.Storage
}

func New(token string) *App {
	return &App{
		token: token,
	}
}

func (a *App) Init() error {
	db, err := sqlite.NewStorage("test.db")
	if err != nil {
		return fmt.Errorf("can't create a database: %s", err.Error())
	}
	a.storage = db
	a.storage.Init(context.Background()) //db.Init(context.Background())
	a.client = tgClient.NewClient(a.token, telegramHost)

	a.processor = tgProcessor.NewProcessor(a.client, a.storage)
	a.consumer = consumer.NewConsumer(a.processor, a.processor, batchSize)

	log.Println("Service initialized")
	return nil
}

func (a *App) Run() {
	log.Println("Service started")
	a.consumer.Consume()
}
