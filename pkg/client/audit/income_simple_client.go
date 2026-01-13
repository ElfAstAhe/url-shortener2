package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	utl "github.com/ElfAstAhe/url-shortener2/internal/utils"
	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
	"github.com/ElfAstAhe/url-shortener2/pkg/utils"
)

type SimpleClient struct {
	client  *http.Client
	baseURL string
}

func NewSimpleClient(baseURL string, timeOut time.Duration) *SimpleClient {
	return &SimpleClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: timeOut},
	}
}

func (client *SimpleClient) AuditIncome(ctx context.Context, data *dto.IncomeAuditDto) error {
	if data == nil {
		return nil
	}

	auditURL, err := url.Parse(client.baseURL)
	if err != nil {
		return NewClientError("Error parsing base URL", err)
	}
	buffer, err := json.Marshal(data)
	if err != nil {
		return NewClientError("Error marshalling audit data", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, auditURL.String(), bytes.NewBuffer(buffer))
	if err != nil {
		return NewClientError("Error creating request", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.client.Do(req)
	if err != nil {
		return NewClientError("Error sending request", err)
	}
	defer utl.CloseOnly(resp.Body)

	if !utils.IsSuccess(resp.StatusCode) {
		return NewClientError(fmt.Sprintf("Error in response with status code [%v]", resp.StatusCode), nil)
	}

	return nil
}
