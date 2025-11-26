package facade

import (
	"context"
	"encoding/json"
	"io"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/mapper"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type UserFacadeImpl struct {
	shortenService service.Shorter
	logger         logger.Logger
}

func NewUserFacadeImpl(shorten service.Shorter, logger logger.Logger) *UserFacadeImpl {
	return &UserFacadeImpl{
		shortenService: shorten,
		logger:         logger.GetLogger("user facade"),
	}
}

func (uf *UserFacadeImpl) BatchDelete(ctx context.Context, rawData io.Reader) error {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return err
	}
	if userInfo == nil {
		return errs.NewAuthUnauthorizedError("user info absent", err)
	}
	var batchData = make(dto.ShortenBatchDeleteRequest, 0)
	dec := json.NewDecoder(rawData)
	err = dec.Decode(&batchData)
	if err != nil {
		return apperrs.NewEpJSONTransformError("dto.ShortenBatchDeleteRequest", err)
	}

	incomeData, err := mapper.UserBatchDeletesFromDto(batchData)

	return uf.shortenService.BatchDelete(ctx, userInfo.UserID, incomeData)
}
