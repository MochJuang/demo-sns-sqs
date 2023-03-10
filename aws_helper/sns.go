package aws_helper

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func (that MessageHelper) Sender(topic string, json string) (_err error) {

	input := &sns.PublishInput{
		Message: aws.String(json),
		TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:" + topic),
	}

	// var output *sns.PublishOutput
	_, _err = that.snsClient.Publish(context.TODO(), input)
	log.Println("Sending to topic " + topic)
	return
}
