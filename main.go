package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func handle() (string, error) {
	return "ok", nil
}

func main() {
	lambda.Start(handle)
}
