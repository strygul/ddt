package main

import (
	"fmt"
	"github.com/strygul/ddt/internal/entity"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	fmt.Println("hello world")
	emptyMap := make(map[string]string)
	step := entity.Step{
		Url:          "https://webhook.site/1495182d-0096-47b9-b369-8e8536c0cde1",
		Method:       "GET",
		Headers:      emptyMap,
		Body:         "testing",
		Placeholders: emptyMap,
	}
	step.ExecuteRequest()
	request, err := step.GetRequest()
	if err != nil {
		println(err.Error())
	}
	timeout := 5 * time.Second
	client := http.Client{Timeout: timeout}
	do, err := client.Do(request)
	defer do.Body.Close()
	if err != nil {
		println(err.Error())
	} else {
		body, err := ioutil.ReadAll(do.Body)
		if err != nil {
			println(err.Error())
		} else {
			println(string(body))
		}
	}
}
