package main

import (
	"flag"
	"log"
	tgClient "tgBot/clients/telegram"
	consumer "tgBot/consumer"
	tgProcessor "tgBot/events/telegram"
)

const (
	telegramHost = "api.telegram.org"
	batchSize    = 10
)

func main() {
	token := mustToken()

	client := tgClient.NewClient(token, telegramHost)
	processor := tgProcessor.NewProcessor(client)
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
