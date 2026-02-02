package grpc

import (
	"context"

	pb "github.com/ElfAstAhe/url-shortener2/api/proto"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type ShortenGRPCFacade struct {
	service service.Shorter
	log     logger.Logger
}

func NewShortenGRPCFacade(service service.Shorter, log logger.Logger) *ShortenGRPCFacade {
	return &ShortenGRPCFacade{
		service: service,
		log:     log.GetLogger("shorten-grpc-facade"),
	}
}

func (sf *ShortenGRPCFacade) ListAllUserURLs(ctx context.Context) (*pb.UserURLsResponse, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	modelData, err := sf.service.GetAllUserShorts(ctx, userInfo.UserID)
	if err != nil {
		return nil, err
	}

	return pb.UserURLsResponse_builder{
		Url: UserURLsToURLsData(modelData),
	}.Build(), nil
}
