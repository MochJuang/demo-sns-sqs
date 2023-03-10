package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"test-sns/aws_helper"
	"time"
)

type Company struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Location string `json:"location"`
}

type Request struct{
	MessageId string
	Message string
}

type Queue struct{
	message aws_helper.MessageHelper
}

func NewQueue(	mess aws_helper.MessageHelper) Queue {
	return Queue{
		message : mess,
	}
}

var companies = []Company{}
var isTimeout bool = true

func (that Queue) QueueInsertCompany() (err error) {
	var company = &Company{}
	var message = &Request{}
	that.message.Consumer("create-company-queue", func(jsondata string) (_err error) {
		_err = json.Unmarshal([]byte(jsondata), message)
		if _err != nil {
			_err = aws_helper.CreateError("02", _err.Error(), message.Message, "create-company")
			return
		}

		_err = json.Unmarshal([]byte(message.Message), company)
		
		if _err != nil {
			_err = aws_helper.CreateError("02", _err.Error(), message.Message, "create-company")
			return
		}
		
		if !strings.Contains(company.Name, "-company") {
			_err = aws_helper.CreateError("02", "Invalid company name", message.Message, "create-company")
			return
		}

		if isTimeout {
			time.Sleep(time.Second * 2)
			isTimeout = false
			_err = aws_helper.CreateError("01", "Create company is timeout", message.Message, "create-company")
			return

		}

		companies = append(companies, *company)

		json2, _ := json.Marshal(companies)
		fmt.Println(string(json2))
		log.Println("---------  Success Added Company  ---------")
		return
	})

	return
}
