NAME ?= main

run:
	go run ./cmd/$(NAME)

build:
	go build -o $(NAME).out ./cmd/$(NAME)

tidy:
	go mod tidy

clean:
	rm *.out tmp