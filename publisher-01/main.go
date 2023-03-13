package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"gitlab.com/myhelpers/aws_helpers"
	"log"
	"os"
	"time"
)

func main() {

	var confSqs aws.Config
	var err error
	var sqsClient *sqs.Client
	//var messageHelper aws_helpers.MessageHelper

	confSqs, err = aws_helpers.GetConfigSqs()
	if err != nil {
		fmt.Println(err.Error())
	}

	sqsClient = sqs.NewFromConfig(confSqs)
	var counter = 1

	//_, err = snsClient.CreateTopic(context.TODO(), &sns.CreateTopicInput{
	//	Name: aws.String("change-configuration-ticket-limit"),
	//})
	//if err != nil {
	//	fmt.Println(err.Error())
	//	os.Exit(1)
	//}

	result, err := sqsClient.ListQueues(context.TODO(), &sqs.ListQueuesInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list queues, %v\n", err)
		os.Exit(1)
	}

	for _, queueURL := range result.QueueUrls {
		fmt.Println(queueURL)
	}

	for {
		mess := fmt.Sprintf("{\"CounterNumber\" : %v}", counter)
		sqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
			MessageBody: aws.String(mess),
			QueueUrl:    aws.String("http://sqs:9324/queue/change-configuration-ticket-limit-queue"),
		})
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Println("Send:", mess)
		counter++
		time.Sleep(time.Second * 2)
	}
}
