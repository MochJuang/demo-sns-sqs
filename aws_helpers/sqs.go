package aws_helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
)

type SqsHelper struct {
	SqsClient *sqs.Client
	Cfg       aws.Config
}

func (that *SqsHelper) Connect() (err error) {

	that.Cfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:9324"}, nil
			})),
	)

	if err != nil {
		return
	}

	that.SqsClient = sqs.NewFromConfig(that.Cfg)

	return
}

func (that *SqsHelper) GetUrlQueue(queue string) (out *sqs.GetQueueUrlOutput, err error) {
	out, err = that.SqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		return
	}

	return
}

func (that *SqsHelper) GetMessage(url string) (out *sqs.ReceiveMessageOutput, err error) {

	out, err = that.SqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            &url,
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   60,
	})

	if err != nil {
		return
	}

	return
}
func (that *SqsHelper) Publish(queueName, message string) (err error) {
	_, err = that.SqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String("http://sqs:9324/queue/" + queueName),
	})
	return
}

func (that *SqsHelper) Consumer(queueName string, callback func(json string) (_err error)) (_err error) {
	var queueUrl *sqs.GetQueueUrlOutput
	var message *sqs.ReceiveMessageOutput

	queueUrl, _err = that.GetUrlQueue(queueName)

	if _err != nil {
		fmt.Println("Error get url queue")
		fmt.Println(_err.Error())
		return
	}
	for {
		message, _err = that.GetMessage(*queueUrl.QueueUrl)
		if _err != nil {
			fmt.Println("Error message queue")
			fmt.Println(_err.Error())
			return
		}

		if len(message.Messages) > 0 {

			for _, msg := range message.Messages {
				var body = *msg.Body

				_err = callback(body)

				if _err != nil {
					var errQueue = &ResponseErr{}
					json.Unmarshal([]byte(_err.Error()), errQueue)
					log.Println(errQueue.Err)

					if errQueue.Code == "01" {
						_err = that.Publish(errQueue.QueueName, errQueue.Message)
					}
				}

				dMInput := &sqs.DeleteMessageInput{
					QueueUrl:      queueUrl.QueueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				}

				// var output *sqs.DeleteMessageOutput
				_, _err = that.SqsClient.DeleteMessage(context.TODO(), dMInput)
				if _err != nil {
					fmt.Println("Error deleted message")
					return
				}
			}
			fmt.Println("")
		}
	}
}
