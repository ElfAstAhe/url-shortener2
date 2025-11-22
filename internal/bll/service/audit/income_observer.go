package audit

import (
	"context"

	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
)

type IncomeObserver interface {
	Observe(ctx context.Context, dto *dto.IncomeAuditDto) error
	GetID() string
	Close() error
}
