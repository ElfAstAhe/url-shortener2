package grpc

import (
	"context"

	shortener "github.com/ElfAstAhe/url-shortener2/api/proto"
    "github.com/ElfAstAhe/url-shortener2/internal/app/config"
    "github.com/ElfAstAhe/url-shortener2/internal/bll/service"
    "github.com/ElfAstAhe/url-shortener2/pkg/logger"
    "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AppGRPCService struct {
    shortener.UnimplementedShortenerServiceServer
    service service.Shorter
    log     logger.Logger
    conf    *config.Config
}

func NewAppGRPCService(
    service service.Shorter,
    conf *config.Config,
    logger logger.Logger,
) *AppGRPCService {
    return &AppGRPCService{
        service: service,
        log:     logger,
        conf:    conf,
    }
}

func (as *AppGRPCService) ShortenURL(context.Context, *shortener.URLShortenRequest) (*shortener.URLShortenResponse, error) {
	// ToDo: implement

	return nil, status.Error(codes.Unimplemented, "method ShortenURL not implemented")
}

func (as *AppGRPCService) ExpandURL(context.Context, *shortener.URLExpandRequest) (*shortener.URLExpandResponse, error) {
	// ToDo: implement

	return nil, status.Error(codes.Unimplemented, "method ExpandURL not implemented")
}

func (as *AppGRPCService) ListUserURLs(context.Context, *emptypb.Empty) (*shortener.UserURLsResponse, error) { {
    // ToDo: implement

    return nil, status.Error(codes.Unimplemented, "method ListUserURLs not implemented")
}
