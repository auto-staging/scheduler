
prepare:
	dep ensure -v

build: prepare
	GOOS=linux go build -o ./bin/auto-staging-scheduler -v -ldflags "-X main.commitHash=`git rev-parse HEAD` -X main.buildTime=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X main.branch=`git rev-parse --abbrev-ref HEAD` -X main.version=`git describe --abbrev=0 --tags` -d -s -w" -tags netgo -installsuffix netgo

tests:
	go test ./... -v -cover

run:
	go run main.go