package entity

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
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

type JsonPath string

func (p JsonPath) Split() []string {
	return strings.Split(string(p), ".")
}

type PlaceholderName = string

type Step struct {
	Url                   string
	Method                HttpMethod
	Headers               map[string]string
	Body                  string
	Placeholders          map[PlaceholderName]string
	PlaceholderNameToPath map[PlaceholderName]JsonPath
	Description           string
	next                  *Step
	client                Doer // e.g. a net/*http.Client to use for requests
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

func (s *Step) SetClient(c Doer) {
	s.client = c
}

func (s *Step) SetNext(step *Step) {
	s.next = step
}

//TODO finish different types
func AccessJsonByPath(jsonBytes []byte, path []string) (string, error) {
	get, dataType, _, err := jsonparser.Get(jsonBytes, path...)
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

func (s Step) ExecuteRequest() ([]byte, error) {
	request, err := s.ConstructRequest()
	if err != nil {
		return nil, err
	}
	if s.client == nil {
		s.client = s.defaultClient()
	}
	do, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	body := do.Body
	defer body.Close()
	return ioutil.ReadAll(body)
}

func (s Step) defaultClient() *http.Client {
	timeout := 5 * time.Second
	return &http.Client{Timeout: timeout}

}

func (s Step) resolvePlaceholders(str string) string {
	for k, v := range s.Placeholders {
		strToReplace := fmt.Sprintf("{{%s}}", k)
		str = strings.ReplaceAll(str, strToReplace, v)
	}
	return str
}
