package dto

type ShortenBatchCreateItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
