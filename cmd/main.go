package main

import (
	"net/http"
)

func HelloServer(w http.ResponseWriter, r *http.Request) {
	// create response binary data
	data := []byte(`{"foo": "bar"}`) // slice of bytes
	// write `data` to response
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	//fmt.Println("hello world")
	//emptyMap := make(map[string]string)
	//step := entity.Step{
	//	Url:          "https://webhook.site/1b127957-0d09-4447-a754-2c3c56ca351e",
	//	Method:       entity.Get,
	//	Headers:      emptyMap,
	//	Body:         "testing",
	//	Placeholders: emptyMap,
	//	Description:  "some description of what the step does",
	//	JsonPath:     "",
	//}
	//step.ExecuteRequest()

	m := http.NewServeMux()
	s := http.Server{Addr: ":8000", Handler: m}
	m.HandleFunc("/", HelloServer)
	s.ListenAndServe()

	//
	//request, err := step.ConstructRequest()
	//if err != nil {
	//	println(err.Error())
	//}
	//timeout := 5 * time.Second
	//client := http.Client{Timeout: timeout}
	//do, err := client.Do(request)
	//defer do.Body.Close()
	//if err != nil {
	//	println(err.Error())
	//} else {
	//	body, err := ioutil.ReadAll(do.Body)
	//	if err != nil {
	//		println(err.Error())
	//	} else {
	//		println(string(body))
	//	}
	//}
}
