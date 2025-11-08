include .env
.PHONY: test clean build

build:
	go mod download
	go build -o setup main.go

test:
	go get github.com/newm4n/goornogo
	go test ./... -v -covermode=count -coverprofile=coverage.out
	goornogo -c 20 -i coverage.out

clean: 
	go clean

run:
	go mod download
	go run main.go

migrateup:
	migrate -path migrations -database "${DB_DSN}" -verbose up

migratedown:
	migrate -path migrations -database "${DB_DSN}" -verbose down