package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNewURI_AllCases(t *testing.T) {
	// prepare
	expected := "http://example.com/api/v1/123"
	// act
	t.Run("positive correct data", func(t *testing.T) {
		// act
		actual := BuildNewURI("http://example.com/api/v1/", "123")
		// assert
		assert.Equal(t, expected, actual)
	})
	// act
	t.Run("positive dirty data", func(t *testing.T) {
		// act
		actual := BuildNewURI(" http://example.com/api/v1/ ", " 123 ")
		// assert
		assert.Equal(t, expected, actual)
	})
	// act
	t.Run("negative empty base url", func(t *testing.T) {
		// act
		actual := BuildNewURI("", "123")
		// assert
		assert.Equal(t, "", actual)
	})
	// act
	t.Run("negative empty key", func(t *testing.T) {
		// act
		actual := BuildNewURI("http://example.com/api/v1", " ")
		// assert
		assert.Equal(t, "", actual)
	})
	// act
	t.Run("positive base url wo trailing slash", func(t *testing.T) {
		actual := BuildNewURI("http://example.com/api/v1", "123")
		assert.Equal(t, "http://example.com/api/v1/123", actual)
	})
	// act
	t.Run("negative incorrect base url", func(t *testing.T) {
		actual := BuildNewURI("http:// example.com/api/v1\\v2", "123")
		assert.Equal(t, "", actual)
	})
}
