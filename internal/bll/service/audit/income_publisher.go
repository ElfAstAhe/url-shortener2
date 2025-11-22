package audit

import (
	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
)

type IncomePublisher interface {
	Register(IncomeObserver)
	Deregister(IncomeObserver)
	NotifyAsync(dto *dto.IncomeAuditDto)
	Close() error
}
