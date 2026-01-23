package mapper

import (
	"fmt"
	"testing"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
)

func BenchmarkShortenBatchResponseFromKeys(b *testing.B) {
	// prepare
	b.StopTimer()

	maxSize := 100_000
	data := make(map[string]string, maxSize)

	for i := 0; i < maxSize; i++ {
		data[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}
	var cnt int64 = 0

	// act
	b.StartTimer()
	b.Run("", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res, err := ShortenBatchResponseFromKeys(data)
			if err != nil {
				fmt.Printf("got error [%v]", err)
			} else {
				cnt += int64(len(res))
			}
		}

		b.ReportMetric(float64(cnt), "total/cnt")
	})
}

func BenchmarkShortenBatchResponseFromEntity(b *testing.B) {
	// prepare
	b.StopTimer()

	maxSize := 100_000
	baseURL := "http://localhost:8080"
	data := make(map[string]*model.ShortURI, maxSize)

	for i := 0; i < maxSize; i++ {
		data[fmt.Sprintf("key%d", i)] = buildShortURI(fmt.Sprintf("id%d", i), fmt.Sprintf("%s/value%d", baseURL, i), fmt.Sprintf("key%d", i))
	}
	var cnt int64 = 0

	// act
	b.StartTimer()
	b.Run("", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := ShortenBatchResponseFromEntity(baseURL, data)
			if err != nil {
				fmt.Printf("got error [%v]\n", err)
			} else {
				cnt += int64(len(data))
			}
		}

		b.ReportMetric(float64(cnt), "total/cnt")
	})
}

func BenchmarkUserShortensFromModel(b *testing.B) {
	// prepare
	b.StopTimer()

	maxSize := 100_000
	data := make(map[string]string, maxSize)

	for i := 0; i < maxSize; i++ {
		data[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}

	var cnt int64 = 0

	b.StartTimer()
	b.Run("", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res, err := UserShortensFromModel(data)
			if err != nil {
				fmt.Printf("got error [%v]\n", err)
			} else {
				cnt += int64(len(res))
			}
		}

		b.ReportMetric(float64(cnt), "total/cnt")
	})
}

func buildShortURI(id string, url string, key string) *model.ShortURI {
	res, err := model.NewShortURIFull(id, url, key)
	if err != nil {
		// omit error
	}
	return res
}
