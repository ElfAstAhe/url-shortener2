package audit

import (
	"context"
	"time"

	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit"
	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
)

type IncomeRemoteService struct {
	client  audit.IncomeClient
	observe bool
}

func NewIncomeRemoteService(baseURL string, observe bool) *IncomeRemoteService {
	return &IncomeRemoteService{
		client:  audit.NewSimpleClient(baseURL, 3*time.Second),
		observe: observe,
	}
}

// Closer interface

func (i *IncomeRemoteService) Close() error {
	return nil
}

// IncomeObserver interface

func (i *IncomeRemoteService) Observe(ctx context.Context, dto *dto.IncomeAuditDto) error {
	if !i.observe {
		return nil
	}

	return i.client.AuditIncome(ctx, dto)
}

func (i *IncomeRemoteService) GetID() string {
	return "REMOTE_INCOME_AUDIT"
}
