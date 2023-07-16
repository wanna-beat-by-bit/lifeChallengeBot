APP_PATH := ./cmd/lifeChallengeBot.go

build:
	go build $(APP_PATH)

run:
	./lifeChallengeBot -tgBotToken=$(TOKEN)

# Combine build and run targets
.PHONY: build run
