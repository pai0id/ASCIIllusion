run:
	@go run  cmd/main/main.go

build:
	@go build -o app.exe cmd/main/main.go

testrun:
	@go run  cmd/testing/main.go

testbuild:
	@go build -o app.exe cmd/testing/main.go

tidy:
	@go mod tidy