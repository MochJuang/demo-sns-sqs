package aws_helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type MessageHelper struct {
	snsClient *sns.Client
	sqsClient *sqs.Client
}

func MessageInit(snsClient *sns.Client, sqsClient *sqs.Client) MessageHelper {
	return MessageHelper{
		snsClient : snsClient,
		sqsClient : sqsClient,
	}
}

func (that MessageHelper) GetUrlQueue(queue string) (out *sqs.GetQueueUrlOutput, err error) {

	out, err = that.sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return
	}

	return
}

func (that MessageHelper) GetMessage(url string) (out *sqs.ReceiveMessageOutput, err error) {

	out, err = that.sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl: &url,
		MaxNumberOfMessages: 1,
		VisibilityTimeout: 60,
	})

	if err != nil {
		return
	}

	return
}

func (that MessageHelper) Consumer(queueName string, callback func(json string) (_err error) ) (_err error) {
	var queueUrl *sqs.GetQueueUrlOutput
	var message *sqs.ReceiveMessageOutput

	queueUrl, _err = that.GetUrlQueue(queueName)
	
	if _err != nil {
		fmt.Println("Error get url queue")
		fmt.Println(_err.Error())
		return
	}
	// fmt.Println(*queueUrl.QueueUrl)
	message, _err = that.GetMessage(*queueUrl.QueueUrl)
	if _err != nil {
		fmt.Println("Error message queue")
		fmt.Println(_err.Error())
		return
	}
	
	if len(message.Messages) > 0 {
		fmt.Println("==================== CONSUMER "+ strings.ToUpper(queueName) +" ================================")

		for _, msg := range message.Messages {
			var body = *msg.Body
			
			_err = callback(string(body))
			
			if _err != nil {
				var errQueue = &ResponseErr{}
				json.Unmarshal([]byte(_err.Error()), errQueue)
				log.Println(errQueue.Err)
		
				if errQueue.Code == "01" {
					that.Sender(errQueue.Topic, errQueue.Message)
				}
			}

			dMInput := &sqs.DeleteMessageInput{
				QueueUrl:      queueUrl.QueueUrl,
				ReceiptHandle: msg.ReceiptHandle,
			}
	
			// var output *sqs.DeleteMessageOutput
			_, _err = that.sqsClient.DeleteMessage(context.TODO(), dMInput)
			if _err != nil {
				fmt.Println("Error deleted message")
				return
			}
		}
		fmt.Println("")		
	}
	
	return
}