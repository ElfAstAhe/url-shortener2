package facade

import (
	"context"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
)

type ToolFacadeImpl struct {
	repo repository.DBConnCheckRepository
}

func NewToolFacadeImpl(connCheckRepo repository.DBConnCheckRepository) *ToolFacadeImpl {
	return &ToolFacadeImpl{
		repo: connCheckRepo,
	}
}

func (tf *ToolFacadeImpl) Ping(ctx context.Context) error {
	if err := tf.repo.CheckDBConn(ctx); err != nil {
		return err
	}

	return nil
}
