package models

import (
	"encoding/json"
	"time"
)

type ImpactReport struct {
	ID          int64           `json:"id"`
	ServiceName string          `json:"serviceName"`
	OldVersion  string          `json:"oldVersion"`
	NewVersion  string          `json:"newVersion"`
	Result      json.RawMessage `json:"result"`
	CreatedAt   time.Time       `json:"createdAt"`
}
