package models

import "time"

type HealthResponse struct {
	Status    string    `json:"status"`
	Operator  string    `json:"operator"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address,omitempty"`
	Version   string    `json:"version,omitempty"`
}
