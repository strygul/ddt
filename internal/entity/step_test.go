package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsingOfPlaceholders(t *testing.T) {
	s := Step{"",
		"",
		map[string]string{},
		"",
		map[string]string{"piggy's name": "Naf Naf", "what does piggy say": "oink"}}

	result := s.resolvePlaceholders("And then {piggy's name} said, `{what does piggy say}`.")
	assert.Equal(t, "And then Naf Naf said, `oink`.", result)
}
