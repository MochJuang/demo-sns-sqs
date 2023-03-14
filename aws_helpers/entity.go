package aws_helpers

import (
	"encoding/json"
	"errors"
	"gitlab.com/ptami_lib/util"
)

type ResponseErr struct {
	Code      string
	Message   string
	Err       string
	QueueName string
}

type MessageBody struct {
	ReceiverName string
	Body         string
	MessageId    *string
}

func (m MessageBody) ToJson() (message string) {
	message = util.StructToString(m)
	return
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
