package utils

import (
	"Job-Queue/internal/model"
	"Job-Queue/pkg"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

var WebhookSecret = "my-secret-key"

func SendWebhookNotification(job *model.Job, q *model.RedisQueue) {
	body, _ := json.Marshal(map[string]any{
		"job_id":   job.ID,
		"status":   job.Status,
		"attempts": job.Attempts,
		"type":     job.Type,
		"payload":  job.Payload,
	})
	signature := computeHMAC(body, WebhookSecret)

	req, _ := http.NewRequest("POST", job.WebhookURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature", "sha256="+signature)

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode >= 400 {
		pkg.Log.WithError(err).WithFields(logrus.Fields{
			"job_id": job.ID,
			"url":    job.WebhookURL,
		}).Error("Failed to send webhook")
		if job.Status == "failed" {
			err := q.FailJob(job.ID)
			if err == nil {
				pkg.Log.WithFields(logrus.Fields{
					"job_id": job.ID,
				}).Warn("Sedning Job to DLQ")
			}

		}
	}
}

func computeHMAC(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}
