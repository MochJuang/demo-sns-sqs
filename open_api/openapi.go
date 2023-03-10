package open_api

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/swaggest/openapi-go/openapi3"
	"io/ioutil"
	"net/http"
)

var reflector openapi3.Reflector

var Handlers []Handler

type Authorization struct {
	Authorization string `header:"Authorization" required:"true"'`
}

type Handler struct {
	Method        string `json:"method"`
	Resource      string `json:"resource"`
	FunctionName  string `json:"functionName"`
	Summary       string
	Request       interface{}
	Response      []OpenApiResponse
	Authorization *string `json:"authorization"`
}

type ResponseData struct {
	Error            string            `json:"error"`
	Message          string            `json:"message,omitempty"`
	Data             interface{}       `json:"data"`
	Count            *uint64           `json:"count,omitempty"`
	Start            *uint64           `json:"start,omitempty"`
	Total            *uint64           `json:"total,omitempty"`
	LastEvaluatedKey interface{}       `json:"lastEvaluatedKey,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
}

type OpenApiResponse struct {
	Code   int
	Output interface{}
}

func selectMethod(req string) (method string, err error) {

	switch req {
	case "GET":
		method = http.MethodGet
	case "PUT":
		method = http.MethodPut
	case "DELETE":
		method = http.MethodDelete
	case "PATCH":
		method = http.MethodPatch
	case "POST":
		method = http.MethodPost
	default:
		err = errors.New("Method not found")
	}

	if err != nil {
		return
	}

	return
}

func GenerateOpenApi(title string, description string, filename string) (err error) {
	reflector = openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	reflector.Spec.Info.
		WithTitle(title).
		WithVersion("1.2.3").
		WithDescription(description)

	for _, request := range Handlers {
		var method string

		api := openapi3.Operation{
			Summary: aws.String(fmt.Sprintf("[%s][%s] %s", request.Method, request.Resource, request.Summary)),
		}

		method, err = selectMethod(request.Method)
		if nil != err {
			return
		}

		err = reflector.SetRequest(&api, request.Request, method)
		if err != nil {
			return
		}

		if len(request.Response) < 1 {
			err = errors.New("Response cannot be null!")
		}

		for _, item := range request.Response {
			item.Output = ResponseData{Data: item.Output}
			err = reflector.SetJSONResponse(&api, item.Output, item.Code)
			if err != nil {
				return
			}
		}

		err = reflector.Spec.AddOperation(method, request.Resource, api)

		if err != nil {
			return
		}
	}

	var schema []byte
	schema, err = reflector.Spec.MarshalYAML()
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filename+".yml", schema, 0644)
	if err != nil {
		errors.New("Unable to write data into the file")
	}

	fmt.Println(string(schema))

	fmt.Println("Success exported !")
	return
}
