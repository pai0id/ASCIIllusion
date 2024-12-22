NAME ?= main

run:
	CGO_ENABLED=0 go run ./cmd/$(NAME)

build:
	CGO_ENABLED=0 go build -o $(NAME).out ./cmd/$(NAME)

tidy:
	go mod tidy

clean:
	rm *.out tmp

default:
	make build
	./main.out -font-config=fonts/IBM_config_cool.json -slice=fonts/slice.json -font-file=fonts/IBM.ttf