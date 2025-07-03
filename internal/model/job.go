package model

import (
	"encoding/json"
	"strconv"
	"time"
)

type Job struct {
	ID          int64
	Type        string
	Attempts    int64
	Payload     map[string]any
	Status      string
	CreatedAt   time.Time
	CompletedAt time.Time
	Delay_ms    int    `json:"delay_ms,omitempty"`
	WebhookURL  string `json:"webhook_url,omitempty"`
}

type JobRequestPayload struct {
	Type       string         `json:"type"`
	Payload    map[string]any `json:"payload"`
	Delay_ms   int            `json:"delay_ms"`
	WebhookURL string         `json:"webhook_url,omitempty"`
}

func HydrateJob(data map[string]string, id string) *Job {
	attempts, _ := strconv.Atoi(data["attempts"])
	jobId, _ := strconv.Atoi(id)
	createdAtUnix, _ := strconv.ParseInt(data["created_at"], 10, 64)
	var payload map[string]any
	json.Unmarshal([]byte(data["payload"]), &payload)

	return &Job{
		ID:        int64(jobId),
		Status:    data["status"],
		Type:      data["type"],
		Attempts:  int64(attempts),
		CreatedAt: time.Unix(createdAtUnix, 0),
		Payload:   payload,
		WebhookURL: data["webhook_url"],
	}
}
