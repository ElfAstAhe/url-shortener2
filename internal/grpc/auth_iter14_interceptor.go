package grpc

import (
	"context"
	"slices"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthIter14Interceptor struct {
	methods []string
	log     logger.Logger
}

func NewAuthIter14Interceptor(methods []string, log logger.Logger) *AuthIter14Interceptor {
	return &AuthIter14Interceptor{
		methods: methods,
		log:     log.GetLogger("gRPC auth-iter14"),
	}
}

func (aii *AuthIter14Interceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	aii.log.Info("AuthIter14Interceptor start")
	defer aii.log.Info("AuthIter14Interceptor finish")

	aii.log.Infof("AuthIter14Interceptor full method [%s]", info.FullMethod)

	// только в зарегистрированных методах gRPC
	if len(aii.methods) > 0 && slices.Contains(aii.methods, info.FullMethod) {
		// работаем только при наличии метаданных
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get(MetaDataAuthorization)
			// работаем только при отсутствии токена
			if len(values) == 0 {
				token, err := auth.NewJWTStringFromUserInfo(auth.BuildRandomUserInfo())
				// добавляем токен (UserInfo)
				if err == nil {
					md.Append(MetaDataAuthorization, token)
				} else {
					aii.log.Errorf("error add jwt token into gRPC request [%s] meta data info with error [%v]", info.FullMethod, err)
				}
			}
		}
	}

	return handler(ctx, req)
}
