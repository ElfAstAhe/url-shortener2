package audit

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type IncomeEventService struct {
	observers map[string]IncomeObserver
	log       logger.Logger
}

func NewIncomeEventService(log logger.Logger) *IncomeEventService {
	return &IncomeEventService{
		observers: make(map[string]IncomeObserver),
		log:       log,
	}
}

// Closer interface

func (i *IncomeEventService) Close() error {
	var eg errgroup.Group
	for _, observer := range i.observers {
		observer := observer
		eg.Go(func() error {
			return observer.Close()
		})
	}

	return eg.Wait()
}

// IncomePublisher interface

func (i *IncomeEventService) Register(observer IncomeObserver) {
	i.observers[observer.GetID()] = observer
}

func (i *IncomeEventService) Deregister(observer IncomeObserver) {
	delete(i.observers, observer.GetID())
}

func (i *IncomeEventService) NotifyAsync(dto *dto.IncomeAuditDto) {
	ctx := context.Background()
	go func() {
		select {
		case <-ctx.Done():
			i.log.Warnf("audit income data context canceled [%v]", zap.Error(ctx.Err()))
			return
		default:
			err := i.notify(ctx, dto)
			if err != nil {
				i.log.Errorf("audit income data got error [%v]", zap.Error(err))
			}
		}
	}()
}

func (i *IncomeEventService) notify(ctx context.Context, dto *dto.IncomeAuditDto) error {
	var eg errgroup.Group
	for _, observer := range i.observers {
		observer := observer
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return observer.Observe(ctx, dto)
			}
		})
	}

	return eg.Wait()
}
