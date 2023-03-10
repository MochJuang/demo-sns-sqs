package main

import (
	"context"
	"fmt"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
	"test-sns/aws_helper"
	"test-sns/helpers"
	"test-sns/open_api"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var svcSqs *sqs.Client
var svcSns *sns.Client

var messageHelper aws_helper.MessageHelper

func getconfig(url string) (cfg aws.Config) {
	var err error
	cfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: url}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "4jr2s", SecretAccessKey: "nr37w8", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)

	if err != nil {
		fmt.Println("Connection to " + url + " failded")
	}

	return
}

type ProductType struct {
	PK               *string    `json:",omitempty" dynamodbav:",omitempty"`
	SK               *string    `json:",omitempty" dynamodbav:",omitempty"`
	GSI1PK           *string    `json:",omitempty" dynamodbav:",omitempty"`
	GSI1SK           *string    `json:",omitempty" dynamodbav:",omitempty"`
	GSI2PK           *string    `json:",omitempty" dynamodbav:",omitempty"`
	GSI2SK           *string    `json:",omitempty" dynamodbav:",omitempty"`
	Id               *string    `json:",omitempty" dynamodbav:",omitempty"`
	Name             *string    `json:",omitempty" dynamodbav:",omitempty"`
	Type             *string    `json:",omitempty" dynamodbav:",omitempty"`
	Stock            *int       `json:",omitempty" dynamodbav:",omitempty"`
	Price            *int       `json:",omitempty" dynamodbav:",omitempty"`
	Product          []Product  `json:",omitempty" dynamodbav:",omitempty"`
	CreatedTimestamp *time.Time `json:",omitempty" dynamodbav:",omitempty"`
	UpdatedTimestamp *time.Time `json:",omitempty" dynamodbav:",omitempty"`
}

type Product struct {
	Id     *string        `json:",omitempty" dynamodbav:",omitempty"`
	Name   *string        `json:",omitempty" dynamodbav:",omitempty"`
	Imei   *string        `json:",omitempty" dynamodbav:",omitempty"`
	Detail *ProductDetail `json:",omitempty" dynamodbav:",omitempty"`
}

type ProductDetail struct {
	Description *string `json:",omitempty" dynamodbav:",omitempty"`
	Date        *string `json:",omitempty" dynamodbav:",omitempty"`
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func createOpenApi() {
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	reflector.Spec.Info.
		WithTitle("Things API").
		WithVersion("1.2.3").
		WithDescription("Put something here")

	type req struct {
		ID     string `path:"id" example:"XXX-XXXXX"`
		Locale string `query:"locale" pattern:"^[a-z]{2}-[A-Z]{2}$"`
		Title  string `json:"string"`
		Amount uint   `json:"amount"`
		Items  []struct {
			Count uint   `json:"count"`
			Name  string `json:"name"`
		} `json:"items"`
	}

	type resp struct {
		ID     string `json:"id" example:"XXX-XXXXX"`
		Amount uint   `json:"amount"`
		Items  []struct {
			Count uint   `json:"count"`
			Name  string `json:"name"`
		} `json:"items"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	openApi := helpers.NewOpenApi("Testing", "Testing openapi")

	openApi.SetOpenApi(helpers.OpenApiOption{
		Path:    "/things",
		Method:  http.MethodGet,
		Request: new(req),
		Response: []helpers.OpenApiResponse{
			{
				Code:   http.StatusOK,
				Output: new(resp),
			},
		},
	})

	openApi.SetOpenApi(helpers.OpenApiOption{
		Path:    "/things/{id}",
		Method:  http.MethodPut,
		Request: new(req),
		Response: []helpers.OpenApiResponse{
			{
				Code:   http.StatusOK,
				Output: new(resp),
			},
			{
				Code:   http.StatusConflict,
				Output: new([]resp),
			},
		},
	})

	err := openApi.GenerateOpenApi("openapi", "")
	if err != nil {
		fmt.Println(err)
	}

}

func main() {

	//var err error
	//var _model = []ProductType{}
	//err = helpers.CreateDummyData(helpers.FakerOption{
	//	Count: 2,
	//	Fields: []helpers.FakerOptions{
	//		{
	//			FieldName: "PK",
	//			Prefix:    aws.String("PRODUCTTYPE#"),
	//			Type:      "uuid",
	//		},
	//		{
	//			FieldName: "SK",
	//			Prefix:    aws.String("PRODUCTTYPE#"),
	//			Type:      "uuid",
	//		},
	//		{
	//			FieldName: "Price",
	//			Type:      "number",
	//		},
	//		{
	//			FieldName: "Product",
	//			Children: &helpers.FakerOption{
	//				Count: 3,
	//				Fields: []helpers.FakerOptions{
	//					{
	//						FieldName: "Id",
	//						Type:      "uuid",
	//						Prefix:    aws.String(""),
	//					},
	//					{
	//						FieldName: "Imei",
	//						Type:      "uuid",
	//						Prefix:    aws.String("IMEI#"),
	//					},
	//					{
	//						FieldName: "Detail",
	//						Child: &helpers.FakerOption{
	//							Fields: []helpers.FakerOptions{
	//								{
	//									FieldName: "Description",
	//									Type:      "fullname",
	//								},
	//								{
	//									FieldName: "Date",
	//									Type:      "date",
	//								},
	//							},
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}, &_model)
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//helpers.PrintJson(_model)

	//createOpenApi()
	open_api.ExampleOpenApi()

}
