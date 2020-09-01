package entity

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpMethod string

const (
	get     = "GET"
	Get     = HttpMethod(get)
	head    = "HEAD"
	Head    = HttpMethod(head)
	post    = "POST"
	Post    = HttpMethod(post)
	put     = "PUT"
	Put     = HttpMethod(put)
	delete  = "DELETE"
	Delete  = HttpMethod(delete)
	connect = "CONNECT"
	Connect = HttpMethod(connect)
	options = "OPTIONS"
	Options = HttpMethod(options)
	trace   = "TRACE"
	Trace   = HttpMethod(trace)
	patch   = "PATCH"
	Patch   = HttpMethod(patch)
)

func (hm HttpMethod) String() string {
	switch hm {
	case Get:
		return get
	case Head:
		return head
	case Post:
		return post
	case Put:
		return put
	case Delete:
		return delete
	case Connect:
		return connect
	case Options:
		return options
	case Trace:
		return trace
	case Patch:
		return patch
	default:
		return "unknown"
	}
}

type Step struct {
	Url          string
	Method       HttpMethod
	Headers      map[string]string
	Body         string
	Placeholders map[string]string
	JsonPath     string
	next         *Step
	client       *http.Client
}

func (s Step) SetClient(c *http.Client) {
	s.client = c
}

func (s Step) SetNext(step *Step) {
	s.next = step
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

func (s Step) ConstructRequest() (*http.Request, error) {
	body := s.serializeResolvedBody()
	url := s.resolvePlaceholders(s.Url)
	req, err := http.NewRequest(s.Method.String(), url, body)
	if err != nil {
		return nil, err
	}
	s.addHeaders(req)
	return req, nil
}

func (s Step) addHeaders(req *http.Request) {
	for k, v := range s.Headers {
		req.Header.Set(k, v)
	}
}

func (s Step) serializeResolvedBody() *bytes.Reader {
	return bytes.NewReader([]byte(s.resolvePlaceholders(s.Body)))
}

func (s Step) ExecuteRequest() (io.ReadCloser, error) {
	request, err := s.ConstructRequest()
	if err != nil {
		return nil, err
	}
	if s.client == nil {
		s.client = s.defaultClient()
	}
	do, err := s.client.Do(request)
	defer do.Body.Close()
	return do.Body, nil
}

func (s Step) defaultClient() *http.Client {
	timeout := 5 * time.Second
	return &http.Client{Timeout: timeout}

}

func (s Step) resolvePlaceholders(str string) string {
	for k, v := range s.Placeholders {
		str = strings.ReplaceAll(str, fmt.Sprintf("{%s}", k), v)
	}
	return str
}
