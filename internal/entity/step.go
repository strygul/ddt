package entity

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Step struct {
	Url          string
	Method       string
	Headers      map[string]string
	Body         string
	Placeholders map[string]string
	JsonPath     string
	Next         *Step
}

func (s Step) SetNext(step *Step) {
	s.Next = step
}

func (s Step) parseJsonPath() []string {
	return strings.Split(s.JsonPath, ".")
}

//finish different types
func AccessResponseBodyByJsonPath(responseBody io.ReadCloser, path []string) (string, error) {
	bytes, err := ioutil.ReadAll(responseBody)
	get, dataType, _, err := jsonparser.Get(bytes, path...)
	if err != nil {
		return "", err
	}
	switch dataType {
	case jsonparser.String:
		return string(get), nil
	case jsonparser.Object:
		return string(get), nil
	default:
		return "", errors.New("Could not parse data: unknown data type")
	}
}

func (s Step) GetRequest() (*http.Request, error) {
	body := bytes.NewReader([]byte(s.resolvePlaceholders(s.Body)))
	url := s.resolvePlaceholders(s.Url)
	req, err := http.NewRequest(s.Method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range s.Headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

func (s Step) ExecuteRequest() (io.ReadCloser, error) {
	request, err := s.GetRequest()
	if err != nil {
		log.Error("Could not execute step request. Message: " + err.Error())
		return nil, err
	}
	timeout := 5 * time.Second
	client := http.Client{Timeout: timeout}
	do, err := client.Do(request)
	defer do.Body.Close()
	return do.Body, nil
}

func (s Step) resolvePlaceholders(str string) string {
	for k, v := range s.Placeholders {
		str = strings.ReplaceAll(str, fmt.Sprintf("{%s}", k), v)
	}
	return str
}
