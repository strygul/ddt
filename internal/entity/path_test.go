package entity

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestAddStep(t *testing.T) {
	p := Path{}
	s := Step{
		Url:                   "http://localhost:80",
		Method:                Get,
		Headers:               make(map[string]string),
		Body:                  "",
		Placeholders:          make(map[string]string),
		PlaceholderNameToPath: make(map[PlaceholderName]JsonPath),
	}
	s2 := Step{
		Url:                   "http://localhost:80",
		Method:                Get,
		Headers:               make(map[string]string),
		Body:                  "",
		Placeholders:          make(map[string]string),
		PlaceholderNameToPath: make(map[PlaceholderName]JsonPath),
	}
	p.AddStep(&s)
	assert.Equal(t, 1, len(p.steps))

	p.AddStep(&s2)
	assert.Equal(t, 2, len(p.steps))
	assert.Equal(t, &s2, p.steps[0].next)
}

func TestPath(t *testing.T) {
	server := httptest.NewServer(&MyHandler{})
	defer server.Close()

	p := Path{}
	s := Step{
		Url:                   server.URL + "/step1",
		Method:                Get,
		Headers:               make(map[string]string),
		Body:                  "",
		Placeholders:          make(map[string]string),
		PlaceholderNameToPath: map[PlaceholderName]JsonPath{"userId": "userId"},
	}
	s2 := Step{
		Url:                   server.URL + "/step2/{{userId}}",
		Method:                Get,
		Headers:               make(map[string]string),
		Body:                  "",
		Placeholders:          make(map[PlaceholderName]string),
		PlaceholderNameToPath: make(map[PlaceholderName]JsonPath),
	}
	s3 := Step{
		Url:                   server.URL + "/step3/{{userId}}",
		Method:                Get,
		Headers:               make(map[string]string),
		Body:                  "",
		Placeholders:          make(map[PlaceholderName]string),
		PlaceholderNameToPath: make(map[PlaceholderName]JsonPath),
	}
	p.AddStep(&s)
	assert.Equal(t, 1, len(p.steps))

	p.AddStep(&s2)
	assert.Equal(t, 2, len(p.steps))
	assert.Equal(t, &s2, p.steps[0].next)

	p.AddStep(&s3)
	assert.Equal(t, 3, len(p.steps))
	assert.Equal(t, &s3, p.steps[1].next)

	err := p.Execute()
	assert.NoError(t, err, "No error when executing the path")
}

type MyHandler struct {
	sync.Mutex
	count int
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data []byte

	switch r.URL.Path {
	case "/step1":
		data = []byte(`{"userId": "myUserId"}`)
	case "/step2/myUserId":
		data = []byte(`{"name": "Vasyl", "surname": "Lomachenko", "alias": "Hi-Tech", "stance": "southpaw"}`)
	case "/step3/myUserId":
		data = []byte(`Success`)
	default:
		log.Fatal(fmt.Sprintf("Error: unexpected call: %s", r.URL.Path))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
