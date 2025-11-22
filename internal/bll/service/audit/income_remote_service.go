package audit

import (
	"context"
	"time"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit"
	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
)

type IncomeRemoteService struct {
	client audit.IncomeClient
}

func NewIncomeRemoteService(appConfig config.Config) *IncomeRemoteService {
	return &IncomeRemoteService{
		client: audit.NewSimpleClient(appConfig.AuditURL, 3*time.Second),
	}
}

// IncomeObserver interface

func (i *IncomeRemoteService) Observe(ctx context.Context, dto *dto.IncomeAuditDto) error {
	return i.client.AuditIncome(ctx, dto)
}

func (i *IncomeRemoteService) GetID() string {
	return "REMOTE_INCOME_AUDIT"
}

// Closer interface

func (i *IncomeRemoteService) Close() error {
	return nil
}
