package main

import (
	"flag"
	"log"

	"github.com/wanna-beat-by-bit/lifeChallengeBot/internals/app"
)

// token: 5949295798:AAGuGBtl8_nHtyEWRZ_FWwB0375DJfg2LPs
func main() {
	token, config := mustArgs()
	app := app.New(token, config)
	if err := app.Init(); err != nil {
		log.Fatalf("Error whilce creating application: %s", err.Error())
	}

	app.Run()
}

func mustArgs() (string, string) {
	var token string
	var config string

	flag.StringVar(
		&token,
		"tgBotToken",
		"",
		"Specify a telegram bot token",
	)

	flag.StringVar(
		&config,
		"config",
		"",
		"Specify a config path",
	)

	flag.Parse()

	if token == "" {
		log.Fatal("Must specify a bot token!")
	}
	if config == "" {
		log.Fatal("Must specify a config!")
	}

	return token, config
}
