package entity

import (
	"github.com/stretchr/testify/assert"
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
