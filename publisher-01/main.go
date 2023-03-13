package main

import (
	"fmt"
	"gitlab.com/myhelpers/aws_helpers"
	"log"
	"time"
)

func main() {

	var err error
	var sqsHelper aws_helpers.SqsHelper

	err = sqsHelper.Connect()
	if err != nil {
		log.Println(err.Error())
	}
	var counter = 1

	for {
		mess := fmt.Sprintf("{\"CounterNumber\" : %v}", counter)
		err = sqsHelper.Publish("change-configuration-ticket-limit-queue", mess)
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Println("Send:", mess)
		counter++
		time.Sleep(time.Second * 2)
	}
}
