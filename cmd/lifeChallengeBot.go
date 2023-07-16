package main

import (
	"flag"
	"fmt"
	"log"
)

// token: 5949295798:AAGuGBtl8_nHtyEWRZ_FWwB0375DJfg2LPs
func main() {
	token := mustToken()
	fmt.Println("Provided token:", token)
	//app := app.New(token)
	//if err := app.Init(); err != nil {
	//	log.Fatalf("Error whilce creating application: %s", err.Error())
	//}

	//app.Run()
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
