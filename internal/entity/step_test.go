package entity

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParsingOfPlaceholders(t *testing.T) {
	s := Step{
		"",
		"",
		map[string]string{},
		"",
		map[string]string{"piggy's name": "Naf Naf", "what does piggy say": "oink"},
		"",
		nil,
		nil,
	}

	result := s.resolvePlaceholders("And then {{piggy's name}} said, `{{what does piggy say}}`.")
	assert.Equal(t, "And then Naf Naf said, `oink`.", result)
}

func TestAccessResponseBodyByJsonPath(t *testing.T) {
	json := "{\"foo\":  {\"bar\":  [{\"baz\":  \"targetString\"}]}}"
	body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	obj, err := AccessResponseBodyByJsonPath(body, strings.Split("foo.bar.[0].baz", "."))
	assert.Nil(t, err, "Should be no errors")
	assert.Equal(t, "targetString", obj, "parsed string should be equal to the target")
}

type HttpClientMock struct {
}

func (c *HttpClientMock) Do(r *http.Request) (*http.Response, error) {
	defer r.Body.Close()
	all, err := ioutil.ReadAll(r.Body)
	recorder := httptest.NewRecorder()
	recorder.Write(all)
	return recorder.Result(), err
}

func TestStepExecution(t *testing.T) {
	placeholders := make(map[string]string)
	placeholders["foo"] = "bar"
	body := "testing"
	step := Step{
		Url:          "https://webhook.site/1b127957-0d09-4447-a754-2c3c56ca351e",
		Method:       Get,
		Headers:      make(map[string]string),
		Body:         body,
		Placeholders: placeholders,
		JsonPath:     "",
	}
	step.SetClient(&HttpClientMock{})
	r, err := step.ExecuteRequest()
	assert.NoError(t, err, "Should be no error")
	assert.Equal(t, body, string(r))
}
