package grpc

import (
	shortener "github.com/ElfAstAhe/url-shortener2/api/proto"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
)

func UserURLToURLData(shortURL, originalURL string) *shortener.URLData {
	return shortener.URLData_builder{
		OriginalUrl: originalURL,
		ShortUrl:    shortURL,
	}.Build()
}

func UserURLsToURLsData(shorts model.UserShorts) []*shortener.URLData {
	res := make([]*shortener.URLData, 0, len(shorts))

	for original, short := range shorts {
		res = append(res, UserURLToURLData(short, original))
	}

	return res
}
