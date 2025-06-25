package model

import "time"

type Job struct {
	ID        int64
	Type 	string
	Attempts     int64
	Payload   map[string]any
	Status    string
	CreatedAt time.Time
	CompletedAt time.Time
}

type JobRequestPayload struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}
