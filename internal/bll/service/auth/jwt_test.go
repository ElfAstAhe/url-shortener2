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

func BenchmarkNewJWTStringFromUserInfo(b *testing.B) {
	b.StopTimer()

	type result struct {
		jwtStr string
		//lint:ignore U1000 это грёбаные тесты
		err error
	}
	maxSize := 10_000

	b.StartTimer()
	var cnt = 0
	// act
	b.Run("generating a new JWT string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < maxSize; j++ {
				var data = make([]result, 0, maxSize)
				var res = result{}
				resJWT, _ := NewJWTStringFromUserInfo(BuildRandomUserInfo())
				res.jwtStr = resJWT

				data = append(data, res)

				cnt += len(data)
			}
		}

		b.ReportMetric(float64(cnt), "total/cnt")
	})
}
