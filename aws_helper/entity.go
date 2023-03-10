package aws_helper

import (
	"encoding/json"
	"errors"
)

type ResponseErr struct {
	Code    string
	Message string
	Err     string
	Topic     string
}

func CreateError(code string, err string, message string, topic string) (_err error) {
	responseErr, _ := json.Marshal(ResponseErr{
		Code: code,
		// Payload
		Message: string(message),
		Err:     err,
		Topic: topic,
	})
	_err = errors.New(string(responseErr))
	return
}
