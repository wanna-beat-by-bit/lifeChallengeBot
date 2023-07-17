APP_PATH := ./cmd/lifeChallengeBot.go
CONFIG_PATH := ./config/local.yaml

build: 
	go build $(APP_PATH)

run:
	./lifeChallengeBot -tgBotToken=$(TOKEN) -config=$(CONFIG_PATH)

# Combine build and run targets
.PHONY: build run
