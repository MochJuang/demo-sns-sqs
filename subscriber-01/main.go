package main

import (
	"gitlab.com/myhelpers/aws_helpers"
	"log"
)

func main() {

	var err error
	var sqsHelper aws_helpers.SqsHelper

	err = sqsHelper.Connect("counter")
	if err != nil {
		log.Println(err.Error())
	}

	var forever = make(chan struct{})

	go sqsHelper.Consumer("change-configuration-ticket-limit-queue-02", func(jsonData aws_helpers.MessageBody) (_err error) {
		log.Println(jsonData)
		return
	})
	<-forever
}
