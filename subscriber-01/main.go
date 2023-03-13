package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"gitlab.com/myhelpers/aws_helpers"
	"log"
)

func main() {

	var confSns aws.Config
	var confSqs aws.Config
	var err error
	var snsClient *sns.Client
	var sqsClient *sqs.Client
	var messageHelper aws_helpers.MessageHelper

	confSqs, err = aws_helpers.GetConfigSqs()
	if err != nil {
		fmt.Println(err.Error())
	}

	confSns, err = aws_helpers.GetConfigSns()
	if err != nil {
		fmt.Println(err.Error())
	}

	snsClient = sns.NewFromConfig(confSns)
	sqsClient = sqs.NewFromConfig(confSqs)
	messageHelper = aws_helpers.MessageInit(snsClient, sqsClient)
	var forever = make(chan struct{})

	go messageHelper.Consumer("change-configuration-ticket-limit-queue", func(jsonData string) (_err error) {
		log.Println(jsonData)
		return
	})
	<-forever
}
