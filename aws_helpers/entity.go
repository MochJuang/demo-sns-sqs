package aws_helpers

import (
	"encoding/json"
	"errors"
)

type ResponseErr struct {
	Code      string
	Message   string
	Err       string
	QueueName string
}

func CreateError(code string, err string, message string, queueName string) (_err error) {
	responseErr, _ := json.Marshal(ResponseErr{
		Code: code,
		// Payload
		Message:   string(message),
		Err:       err,
		QueueName: queueName,
	})
	_err = errors.New(string(responseErr))
	return
}
