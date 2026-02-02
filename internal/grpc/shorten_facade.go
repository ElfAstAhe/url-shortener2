package grpc

import (
	"context"
	"errors"

	pb "github.com/ElfAstAhe/url-shortener2/api/proto"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type ShortenGRPCFacade struct {
	service service.Shorter
	log     logger.Logger
	baseURL string
}

func NewShortenGRPCFacade(service service.Shorter, baseURL string, log logger.Logger) *ShortenGRPCFacade {
	return &ShortenGRPCFacade{
		service: service,
		baseURL: baseURL,
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

func (sf *ShortenGRPCFacade) CreateURL(ctx context.Context, req *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	modelData, err := sf.service.Store(ctx, userInfo.UserID, req.GetUrl())
	if err != nil && !errors.As(err, &apperrs.BllConflictErr) {
		return nil, err
	}

	return pb.URLShortenResponse_builder{
		Result: KeyToURL(sf.baseURL, modelData),
	}.Build(), err
}

func (sf *ShortenGRPCFacade) GetURL(ctx context.Context, req *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	key, err := sf.service.GetURL(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if key == "" {
		return nil, nil
	}

	resp := pb.URLExpandResponse_builder{
		Result: KeyToURL(sf.baseURL, key),
	}.Build()

	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return resp, nil
	}

	// get user data
	_, errUser := sf.service.GetURLUser(ctx, userInfo.UserID, key)
	if errUser != nil && !errors.As(errUser, &apperrs.DalSoftRemovedErr) {
		return nil, errUser
	}

	return resp, errUser
}
