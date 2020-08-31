package entity

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestParsingOfPlaceholders(t *testing.T) {
	s := Step{"",
		"",
		map[string]string{},
		"",
		map[string]string{"piggy's name": "Naf Naf", "what does piggy say": "oink"},
		"",
		nil}

	result := s.resolvePlaceholders("And then {piggy's name} said, `{what does piggy say}`.")
	assert.Equal(t, "And then Naf Naf said, `oink`.", result)
}

func TestAccessResponseBodyByJsonPath(t *testing.T) {
	json := "{\"foo\":  {\"bar\":  [{\"baz\":  \"targetString\"}]}}"
	body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	obj, err := AccessResponseBodyByJsonPath(body, strings.Split("foo.bar.[0].baz", "."))
	assert.Nil(t, err, "Should be no errors")
	assert.Equal(t, "targetString", obj, "parsed string should be equal to the target")
}
