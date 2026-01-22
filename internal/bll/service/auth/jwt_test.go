package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJWTStringFromUserInfo_AllCases(t *testing.T) {
	// act
	t.Run("positive info exists", func(t *testing.T) {
		// act
		actual, err := NewJWTStringFromUserInfo(BuildRandomUserInfo())
		assert.NoError(t, err)
		assert.NotEqual(t, "", actual)
	})
	// act
	t.Run("positive info nil", func(t *testing.T) {
		actual, err := NewJWTStringFromUserInfo(nil)
		assert.Error(t, err)
		assert.Equal(t, "", actual)
	})
}

func BenchNewJWTStringFromUserInfo(b *testing.B) {
	b.StopTimer()

	type result struct {
		jwtStr string
		err    error
	}

	var data = make([]result, 0, b.N)

	b.StartTimer()
	// act
	b.Run("generating a new JWT string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resJWT, err := NewJWTStringFromUserInfo(BuildRandomUserInfo())
			data[i].jwtStr = resJWT
			data[i].err = err
		}
	})
}
