NAME ?= main

run:
	go run ./cmd/$(NAME)

tidy:
	go mod tidy