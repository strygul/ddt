package entity

import (
	"bytes"
	"fmt"
	"io"
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
	JsonPath	string
	Next 		*Step
}

func (s Step) SetNext(step *Step) {
	s.Next = step
}

func (s Step) parseJsonPath() []string {
	return strings.Split(s.JsonPath, ".")
}

func (s Step) GetRequest()(*http.Request, error) {
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