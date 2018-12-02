package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(eventJson json.RawMessage) error {

	return nil
}

func main() {
	lambda.Start(Handler)
}
