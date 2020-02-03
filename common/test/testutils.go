package test_utils

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"gotest.tools/assert"
)

func MapDecode(data interface{}, s interface{}) {
	mapstructure.Decode(data, s)
}

func Assert(t *testing.T, condition bool) {
	assert.Assert(t, condition)
}

func Equal(t *testing.T, a interface{}, b interface{}) {
	assert.Equal(t, a, b)
}
