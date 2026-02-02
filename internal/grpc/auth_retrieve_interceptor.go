package grpc

import (
	"context"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthRetrieveInterceptor struct {
	log logger.Logger
}

func NewAuthRetrieveInterceptor(log logger.Logger) *AuthRetrieveInterceptor {
	return &AuthRetrieveInterceptor{
		log: log.GetLogger("gRPC auth-retrieve"),
	}
}

func (ari *AuthRetrieveInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ari.log.Info("UnaryInterceptor start")
	defer ari.log.Info("UnaryInterceptor finish")

	ctxUserInfo := ctx
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("authorization")
		var userInfo *auth.UserInfo
		var err error
		if len(values) > 0 {
			userInfo, err = auth.UserInfoFromJWTString(values[0])
			if err != nil {
				ari.log.Errorf("error retrieve user info from authorization header: [%v]", err)
			}

			ctxUserInfo = context.WithValue(ctx, auth.ContextUserInfo, userInfo)
		}
	}

	return handler(ctxUserInfo, req)
}
