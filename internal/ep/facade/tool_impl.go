package facade

import (
	"context"
	"net/http"
	"net/netip"
	"strings"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

type ToolFacadeImpl struct {
	toolRepo          repository.DBConnCheckRepository
	shortURIService   service.Shorter
	trustedSubnetCIDR string
}

func NewToolFacadeImpl(
	connCheckRepo repository.DBConnCheckRepository,
	shortURIService service.Shorter,
	trustedSubnetCIDR string,
) *ToolFacadeImpl {
	return &ToolFacadeImpl{
		toolRepo:          connCheckRepo,
		shortURIService:   shortURIService,
		trustedSubnetCIDR: strings.TrimSpace(trustedSubnetCIDR),
	}
}

func (tf *ToolFacadeImpl) Ping(ctx context.Context) error {
	if err := tf.toolRepo.CheckDBConn(ctx); err != nil {
		return err
	}

	return nil
}

func (tf *ToolFacadeImpl) InternalStats(ctx context.Context, header *http.Header) (*dto.InternalStatsDto, error) {
	// подготовка и валидация
	ipAddr, subNetPrefix, err := tf.prepareAndValidateInternalStatsParams(header)
	if err != nil {
		return nil, err
	}

	// бизнес логика
	if !subNetPrefix.Contains(*ipAddr) {
		return nil, errs.NewAuthUnauthorizedError("untrusted ip", nil)
	}

	// данные
	urlCount, userCount, err := tf.shortURIService.InternalStats(ctx)
	if err != nil {
		return nil, err
	}

	return dto.NewInternalStatsDto(urlCount, userCount), nil
}

func (tf *ToolFacadeImpl) prepareAndValidateInternalStatsParams(header *http.Header) (ip *netip.Addr, subnet *netip.Prefix, err error) {
	realIP := strings.TrimSpace(header.Get("X-Real-IP"))
	if err = tf.validateInternalStatsParams(realIP); err != nil {
		return nil, nil, err
	}

	ipAddr, err := netip.ParseAddr(realIP)
	if err != nil {
		return nil, nil, err
	}
	subNetPrefix, err := netip.ParsePrefix(tf.trustedSubnetCIDR)
	if err != nil {
		return nil, nil, err
	}

	return &ipAddr, &subNetPrefix, nil
}

func (tf *ToolFacadeImpl) validateInternalStatsParams(realIP string) error {
	if tf.trustedSubnetCIDR == "" {
		return errs.NewAuthUnauthorizedError("empty subnet", nil)
	}
	if realIP == "" {
		return errs.NewAuthUnauthorizedError("empty ip", nil)
	}

	return nil
}
