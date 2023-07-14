package main

import (
	"context"
	"flag"
	"log"

	tgClient "github.com/wanna-beat-by-bit/lifeChallengeBot/internals/app/clients/telegram"
	consumer "github.com/wanna-beat-by-bit/lifeChallengeBot/internals/app/consumer"
	tgProcessor "github.com/wanna-beat-by-bit/lifeChallengeBot/internals/app/events/telegram"
	storeSqlite "github.com/wanna-beat-by-bit/lifeChallengeBot/internals/pkg/storage/sqlite"
)

const (
	telegramHost = "api.telegram.org"
	batchSize    = 10
)

// token: 5949295798:AAGuGBtl8_nHtyEWRZ_FWwB0375DJfg2LPs
func main() {
	token := mustToken()

	db, err := storeSqlite.NewStorage("test.db")
	if err != nil {
		log.Fatalf("Error while creating database: %s", err.Error())
	}

	db.Init(context.Background())

	client := tgClient.NewClient(token, telegramHost)
	processor := tgProcessor.NewProcessor(client, db)
	consumer := consumer.NewConsumer(processor, processor, batchSize)

	log.Println("Server started")
	consumer.Consume()
}

func mustToken() string {
	token := flag.String(
		"tgBotToken",
		"",
		"Specify a telegram bot token",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("Must specify a bot token!")
	}

	return *token
}
