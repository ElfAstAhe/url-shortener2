package grpc

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/ElfAstAhe/url-shortener2/api/proto"
	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
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
	resp, err := as.facade.CreateURL(ctx, req)
	if err != nil && !errors.As(err, &apperrs.BllConflictErr) {
		// unauthorized
		if errors.As(err, &errs.AuthUnauthorizedErr) || errors.As(err, &errs.AuthInfoInvalidErr) {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("unauthorized with error [%v]", err))
		}

		// internal
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create URLs with error [%v]", err))
	}

	// conflict
	if err != nil && !errors.As(err, &apperrs.BllConflictErr) {
		err = status.Error(codes.AlreadyExists, fmt.Sprintf("URL already exists, error [%v]", err))
	}

	return resp, err
}

func (as *ShortenGRPCService) ExpandURL(ctx context.Context, req *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	resp, err := as.facade.GetURL(ctx, req)
	if err != nil {
		// gone
		if errors.As(err, &apperrs.DalSoftRemovedErr) {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("data gone, error [%v]", err))
		}

		// internal
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get URLs with error [%v]", err))
	}

	// not found
	if resp == nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return resp, nil
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
