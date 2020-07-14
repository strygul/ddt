package entity

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type Step struct {
	url          string
	method       string
	headers      map[string]string
	body         string
	placeholders map[string]string
}


func (s Step) GetRequest()(*http.Request, error){
	req, err := http.NewRequest(s.method, s.resolvePlaceholders(s.url),  bytes.NewReader([]byte(s.resolvePlaceholders(s.body))))
	if err != nil {
		return nil, err
	}
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}
	return req, nil

}

func (s Step) resolvePlaceholders(str string) string {
	for k, v := range s.placeholders {
		str = strings.ReplaceAll(str, fmt.Sprintf("{%s}", k), v)
	}
	return str
}