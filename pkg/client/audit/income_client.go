package audit

import (
	"context"

	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
)

type IncomeClient interface {
	AuditIncome(ctx context.Context, dto *dto.IncomeAuditDto) error
}
