NAME ?= main
# make run NAME=testing

.PHONY: test run build tidy clean default

run:
	CGO_ENABLED=0 go run ./cmd/$(NAME)

build:
	CGO_ENABLED=0 go build -o $(NAME).out ./cmd/$(NAME)

tidy:
	go mod tidy

clean:
	rm *.out tmp

test:
	@echo "Running tests in internal/renderer..."
	@cd internal/renderer && go test ./... -cover
	@echo "Running tests in internal/object..."
	@cd internal/object && go test ./... -cover
	@echo "Running tests in internal/transformer..."
	@cd internal/transformer && go test ./... -cover

default:
	make build
	./main.out -font-config=fonts/IBM_config_cool.json -slice=fonts/slice.json -font-file=fonts/IBM.ttf