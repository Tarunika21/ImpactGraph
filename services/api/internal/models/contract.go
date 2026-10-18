package models

import (
	"encoding/json"
	"time"
)

type CreateContractRequest struct {
	ServiceName   string          `json:"serviceName"`
	Version       string          `json:"version"`
	Specification json.RawMessage `json:"specification"`
}

type Contract struct {
	ID            int64           `json:"id"`
	ServiceName   string          `json:"serviceName"`
	Version       string          `json:"version"`
	Specification json.RawMessage `json:"specification"`
	CreatedAt     time.Time       `json:"createdAt"`
}
