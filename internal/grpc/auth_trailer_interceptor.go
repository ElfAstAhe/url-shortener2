package grpc

import (
	"context"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthTrailerInterceptor struct {
	log logger.Logger
}

func NewAuthTrailerInterceptor(log logger.Logger) *AuthTrailerInterceptor {
	return &AuthTrailerInterceptor{
		log: log.GetLogger("gRPC auth-trailer"),
	}
}

func (ati *AuthTrailerInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ati.log.Info("AuthTrailerInterceptor start")
	defer ati.log.Info("AuthTrailerInterceptor finish")

	ati.log.Infof("AuthTrailerInterceptor full method [%s]", info.FullMethod)

	// passthru pipeline
	resp, respErr := handler(ctx, req)

	// processing
	var userInfo *auth.UserInfo
	var err error
	var tokenString string
	userInfo, err = auth.UserInfoFromContext(ctx)
	if err == nil && userInfo != nil {
		tokenString, err = auth.NewJWTStringFromUserInfo(userInfo)
		if err != nil {
			ati.log.Warnf("has trouble build jwt token string with error [%v]", err)

			return resp, respErr
		}
		metaData := metadata.Pairs(MetaDataAuthorization, tokenString)
		// set authorization info into answer
		err = grpc.SetTrailer(ctx, metaData)
		if err != nil {
			ati.log.Warnf("failed to set trailer with error [%v]", err)
		}
	}

	return resp, respErr
}
