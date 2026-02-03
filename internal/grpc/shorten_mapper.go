package grpc

import (
	pb "github.com/ElfAstAhe/url-shortener2/api/proto"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
)

func UserURLToURLData(shortURL, originalURL string) *pb.URLData {
	return pb.URLData_builder{
		OriginalUrl: originalURL,
		ShortUrl:    shortURL,
	}.Build()
}

func UserURLsToURLsData(shorts model.UserShorts) []*pb.URLData {
	res := make([]*pb.URLData, 0, len(shorts))

	for original, short := range shorts {
		res = append(res, UserURLToURLData(short, original))
	}

	return res
}

func KeyToURL(baseURL, key string) string {
	return utils.BuildNewURI(baseURL, key)
}
