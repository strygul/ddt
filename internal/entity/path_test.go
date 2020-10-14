package entity

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/strygul/ddt/internal/testutility"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddStep(t *testing.T) {
	p := Path{}
	s := Step{
		Url:          "http://localhost:80",
		Method:       Get,
		Headers:      make(map[string]string),
		Body:         "",
		Placeholders: make(map[string]string),
		JsonPath:     "",
	}
	s2 := Step{
		Url:          "http://localhost:80",
		Method:       Get,
		Headers:      make(map[string]string),
		Body:         "",
		Placeholders: make(map[string]string),
		JsonPath:     "",
	}
	p.AddStep(&s)
	assert.Equal(t, 1, len(p.steps))

	p.AddStep(&s2)
	assert.Equal(t, 2, len(p.steps))
	assert.Equal(t, &s2, p.steps[0].next)
}

func TestPath(t *testing.T) {

	server := httptest.NewServer(&testutility.MyHandler{})
	defer server.Close()

	body, err := http.Get(server.URL)
	if err != nil {
		logrus.Error("Could not execute a call to the test server")
	}

	bytes, err := ioutil.ReadAll(body.Body)
	println(string(bytes))

	p := Path{}
	s := Step{
		Url:          server.URL,
		Method:       Get,
		Headers:      make(map[string]string),
		Body:         "",
		Placeholders: make(map[string]string),
		JsonPath:     "",
	}
	s2 := Step{
		Url:          server.URL,
		Method:       Get,
		Headers:      make(map[string]string),
		Body:         "",
		Placeholders: make(map[string]string),
		JsonPath:     "",
	}
	p.AddStep(&s)
	assert.Equal(t, 1, len(p.steps))

	p.AddStep(&s2)
	assert.Equal(t, 2, len(p.steps))
	assert.Equal(t, &s2, p.steps[0].next)
}
