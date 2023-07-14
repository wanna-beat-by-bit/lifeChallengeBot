DOCKER_CONFIG_PATH := docker/docker-compose.yaml
BUILD_CONTEXT := .
GO_CONFIG := ./config/local.yaml
APP_PATH := ./cmd/lifeChallengeBot.go
TOKEN := asdlkfj23ordsf

build:
	go build $(APP_PATH)

run:
	go run $(APP_PATH) -tgBotToken=$(TOKEN)

# Combine build and run targets
.PHONY: build run
