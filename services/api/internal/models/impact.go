package models

type Consumer struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Endpoint string   `json:"endpoint"`
	Fields   []string `json:"fields"`
}

type ImpactedConsumer struct {
	ConsumerID   string   `json:"consumerId"`
	ConsumerName string   `json:"consumerName"`
	Endpoint     string   `json:"endpoint"`
	Reasons      []string `json:"reasons"`
}

type ImpactResponse struct {
	TotalConsumers    int                `json:"totalConsumers"`
	ImpactedCount     int                `json:"impactedCount"`
	ImpactedConsumers []ImpactedConsumer `json:"impactedConsumers"`
}

type ImpactRequest struct {
	Changes   []Change   `json:"changes"`
	Consumers []Consumer `json:"consumers"`
}
