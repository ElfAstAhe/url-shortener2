package grpc

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/ElfAstAhe/url-shortener2/api/proto"
	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ShortenGRPCService struct {
	pb.UnimplementedShortenerServiceServer
	facade *ShortenGRPCFacade
	log    logger.Logger
	conf   *config.Config
}

func NewAppGRPCService(
	facade *ShortenGRPCFacade,
	conf *config.Config,
	logger logger.Logger,
) *ShortenGRPCService {
	return &ShortenGRPCService{
		facade: facade,
		log:    logger.GetLogger("shorten-grpc-service"),
		conf:   conf,
	}
}

func (as *ShortenGRPCService) ShortenURL(ctx context.Context, req *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	// ToDo: implement

	return nil, status.Error(codes.Unimplemented, "method ShortenURL not implemented")
}

func (as *ShortenGRPCService) ExpandURL(ctx context.Context, req *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	// ToDo: implement

	return nil, status.Error(codes.Unimplemented, "method ExpandURL not implemented")
}

func (as *ShortenGRPCService) ListUserURLs(ctx context.Context, req *emptypb.Empty) (*pb.UserURLsResponse, error) {
	resp, err := as.facade.ListAllUserURLs(ctx)
	if err != nil {
		// no content
		if errors.As(err, &errs.AuthInfoAbsentErr) {
			// no content
			return pb.UserURLsResponse_builder{
				Url: make([]*pb.URLData, 0),
			}.Build(), nil
		}

		// unauthorized
		if errors.As(err, &errs.AuthInfoInvalidErr) || errors.As(err, &errs.AuthUnauthorizedErr) {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("unauthorized with error [%v]", err))
		}

		// other errors internal server error
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list user URLs with error [%v]", err))
	}

	return resp, nil
}
