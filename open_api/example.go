package open_api

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"net/http"
)

type Person struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Hobby string `json:"hobby"`
}

type PersonRequest struct {
	Authorization
	Name         string `json:"name" required:"true"`
	Hobby        string `json:"hobby"`
	CustomHeader string `header:"custom-header"`
}

type PersonRequestUpdate struct {
	ID    string `path:"id" json:"id"  example:"XXX-XXXXX"`
	Name  string `json:"name"`
	Hobby string `json:"hobby"`
}

type PersonRequestId struct {
	ID string `path:"id" example:"XXX-XXXXX"`
}

func ExampleOpenApi() {

	Handlers = []Handler{
		{
			Method:       "PUT",
			Resource:     "/person",
			Summary:      "Insert Person",
			FunctionName: "HandlerCreatePerson",
			Request:      new(PersonRequest),
			Response: []OpenApiResponse{
				{
					Code:   http.StatusOK,
					Output: new(Person),
				},
			},
		},
		{
			Method:       "GET",
			Resource:     "/person",
			FunctionName: "HandlerGetPerson",
			Summary:      "Get Person",
			Request:      "",
			Response: []OpenApiResponse{
				{
					Code:   http.StatusOK,
					Output: new([]Person),
				},
			},
			Authorization: aws.String("test"),
		},
		{
			Method:       "GET",
			Resource:     "/person/{id}",
			Summary:      "Get by id person",
			FunctionName: "HandlerGetByIdPerson",
			Request:      new(PersonRequestId),
			Response: []OpenApiResponse{
				{
					Code:   http.StatusOK,
					Output: new(Person),
				},
			},
			Authorization: aws.String("test"),
		},
		{
			Method:       "PATCH",
			Resource:     "/person/{id}",
			FunctionName: "HandlerUpdatePerson",
			Summary:      "Update Person",
			Request:      new(PersonRequestUpdate),
			Response: []OpenApiResponse{
				{
					Code:   http.StatusOK,
					Output: new(Person),
				},
			},
			Authorization: aws.String("test"),
		},
		{
			Method:       "DELETE",
			Resource:     "/person/{id}",
			FunctionName: "HandlerDeletePerson",
			Summary:      "Delete Person",
			Request:      new(PersonRequestId),
			Response: []OpenApiResponse{
				{
					Code:   http.StatusOK,
					Output: new(Person),
				},
			},
			Authorization: aws.String("test"),
		},
	}

	err := GenerateOpenApi("Testing", "For testing openapi generator", "openapi")
	if err != nil {
		fmt.Println(err.Error())
	}
}
