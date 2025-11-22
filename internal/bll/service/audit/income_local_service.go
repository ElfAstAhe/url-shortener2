package audit

import (
	"context"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/dal/storage"
	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
)

type IncomeLocalService struct {
	storage *storage.IncomeAuditStorageWriter
}

func NewIncomeLocalService(appConfig config.Config) (*IncomeLocalService, error) {
	auditStorage, err := storage.NewIncomeAuditStorageWriter(appConfig.AuditFile)
	if err != nil {
		return nil, err
	}

	return &IncomeLocalService{
		storage: auditStorage,
	}, nil
}

// Closer interface

func (i *IncomeLocalService) Close() error {
	return i.storage.Close()
}

// IncomeObserver interface

func (i *IncomeLocalService) Observe(ctx context.Context, dto *dto.IncomeAuditDto) error {
	return i.storage.SaveData(ctx, dto)
}

func (i *IncomeLocalService) GetID() string {
	return "LOCAL_INCOME_AUDIT"
}
