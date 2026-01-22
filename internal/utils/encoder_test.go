package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeURI_AllCases(t *testing.T) {
	// prepare
	bytes := []byte("hello world")
	// act
	t.Run("positive", func(t *testing.T) {
		// act
		actual := EncodeURI(bytes)
		// assert
		assert.True(t, len(actual) > 0)
	})
	t.Run("negative", func(t *testing.T) {
		// act
		actual := EncodeURI(nil)
		// assert
		assert.Nil(t, actual)
	})
}

func TestEncodeURIStr_AllCases(t *testing.T) {
	// prepare
	src := "hello world"
	expected := ""
	// act
	t.Run("positive", func(t *testing.T) {
		// act
		actual := EncodeURIStr(src)
		// assert
		assert.NotEqual(t, expected, actual)
	})
	t.Run("negative", func(t *testing.T) {
		// act
		actual := EncodeURIStr(expected)
		// assert
		assert.Equal(t, expected, actual)
	})
}
