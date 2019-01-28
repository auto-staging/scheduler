
prepare:
	dep ensure -v

build: prepare
	go build -o ./bin/auto-staging-scheduler

tests:
	go test ./... -v -cover

run:
	go run main.go