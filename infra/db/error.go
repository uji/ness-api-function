package db

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
)

func GetAWSErrCode(err error) string {
	if aerr, ok := err.(awserr.Error); ok {
		return aerr.Code()
	}
	return "error is not awserr"
}

func EqualAWSErrCode(err error, code string) bool {
	if aerr, ok := err.(awserr.Error); ok {
		return aerr.Code() == code
	}
	return false
}
