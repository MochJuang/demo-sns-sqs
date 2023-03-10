package helpers

import (
	"encoding/json"
	"fmt"
)

func PrintJson(models interface{}) {
	var jsonData []byte
	var err error
	jsonData, err = json.MarshalIndent(models, "", "\t")
	//jsonData, err = json.Marshal(models)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(jsonData))
}
