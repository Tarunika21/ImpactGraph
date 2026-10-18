package models

import (
	"encoding/json"
	"time"
)

type CreateConsumerRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Consumer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateConsumerDependencyRequest struct {
	ConsumerID  string          `json:"consumerId"`
	ServiceName string          `json:"serviceName"`
	Endpoint    string          `json:"endpoint"`
	Fields      json.RawMessage `json:"fields"`
}

type ConsumerDependency struct {
	ID          int64           `json:"id"`
	ConsumerID  string          `json:"consumerId"`
	ServiceName string          `json:"serviceName"`
	Endpoint    string          `json:"endpoint"`
	Fields      json.RawMessage `json:"fields"`
	CreatedAt   time.Time       `json:"createdAt"`
}
