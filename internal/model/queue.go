package model

import (
	"Job-Queue/internal/config"
	"Job-Queue/pkg"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{
		client: client,
	}
}

func (q *RedisQueue) Enqueue(job *Job) error {
	job.ID = time.Now().UnixMicro()
	job.Status = "queued"
	job.Attempts = 0
	job.CreatedAt = time.Now()
	payloadBytes, err := json.Marshal(job.Payload)
	if err != nil {
		return err
	}
	jobKey := fmt.Sprintf("job:%v", job.ID)
	_, err = q.client.TxPipelined(config.Ctx, func(pipe redis.Pipeliner) error {
		jobData := map[string]any{
			"status":     job.Status,
			"attempts":   job.Attempts,
			"created_at": job.CreatedAt.Unix(),
			"type":       job.Type,
			"payload":    string(payloadBytes),
		}

		pipe.HSet(config.Ctx, jobKey, jobData)
		pipe.ZAdd(config.Ctx, "job_index", redis.Z{
			Score:  float64(job.CreatedAt.Unix()),
			Member: job.ID,
		})
		pipe.RPush(config.Ctx, "worker_queue", job.ID)
		return nil
	})
	return err
}

func (q *RedisQueue) RedisFetcher(id int, jobChan chan *Job) {
	for {
		jobID, err := q.client.BRPopLPush(config.Ctx, "worker_queue", "in_progress_queue", 5*time.Second).Result()
		if err != nil {
			if err != redis.Nil {
				pkg.Log.WithFields(logrus.Fields{
					"fetcher_id": id,
				}).WithError(err).Error("Fetcher failed to fetch")
			}
			continue
		}
		jobData, err := q.client.HGetAll(config.Ctx, "job:"+jobID).Result()
		if err != nil || len(jobData) == 0 {
			pkg.Log.WithFields(logrus.Fields{
				"fetcher_id": id,
				"job_id":     jobID,
			}).WithError(err).Error("failed to fetch job details")
			continue
		}

		job := HydrateJob(jobData, jobID)
		job.Status = "processing"
		if err := q.SaveJob(job); err != nil {
			pkg.Log.WithFields(logrus.Fields{
				"fetcher_id": id,
			}).WithError(err).Error("error in saving job")
			continue
		}

		select {
		case jobChan <- job:
			pkg.Log.WithFields(logrus.Fields{
				"fetcher_id": id,
				"job_id":     job.ID,
				"job_type":   job.Type,
			}).Info("Fetcher enqueued job")
		case <-config.Ctx.Done():
			return
		}
	}
}

func (q *RedisQueue) GetJobByID(jobID string) (*Job, error) {
	jobKey := fmt.Sprintf("job:%v", jobID)
	data, err := q.client.HGetAll(config.Ctx, jobKey).Result()
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}
	return HydrateJob(data, jobID), nil
}

func (q *RedisQueue) SaveJob(job *Job) error {
	jobKey := fmt.Sprintf("job:%v", job.ID)
	_, err := q.client.HSet(config.Ctx, jobKey, map[string]any{
		"status":       job.Status,
		"attempts":     job.Attempts,
		"completed_at": job.CompletedAt.Unix(),
	}).Result()
	return err
}

func (q *RedisQueue) AckJob(jobID int64) error {
	_, err := q.client.LRem(config.Ctx, "in_progress_queue", 0, fmt.Sprintf("%v", jobID)).Result()
	return err
}

func (q *RedisQueue) FailJob(jobID int64) error {
	_, err := q.client.LRem(config.Ctx, "in_progress_queue", 0, fmt.Sprintf("%v", jobID)).Result()
	if err != nil {
		return err
	}
	_, err = q.client.RPush(config.Ctx, "dead_letter_queue", jobID).Result()
	return err
}

func (q *RedisQueue) GetJobs(start int64, end int64) (map[string]*redis.MapStringStringCmd, error) {
	ids, err := q.client.ZRevRange(config.Ctx, "job_index", int64(start), int64(end)).Result()
	if err != nil {
		return nil, err
	}
	pipe := q.client.Pipeline()
	cmdMap := make(map[string]*redis.MapStringStringCmd)

	for _, id := range ids {
		key := fmt.Sprintf("job:%s", id)
		cmdMap[id] = pipe.HGetAll(config.Ctx, key)
	}

	_, err = pipe.Exec(config.Ctx)
	if err != nil {
		return nil, err
	}
	return cmdMap, nil
}

// err = q.client.HSet(config.Ctx, jobKey, map[string]any{
// 	"status":     job.Status,
// 	"attempts":   job.Attempts,
// 	"created_at": job.CreatedAt.Unix(),
// 	"type":       job.Type,
// 	"payload":    string(payloadBytes),
// }).Err()

// if err != nil {
// 	return err
// }
// _, err = q.client.ZAdd(config.Ctx, "job_index", redis.Z{
// 	Score:  float64(job.ID),
// 	Member: job.ID,
// }).Result()

// if err != nil {
// 	return err
// }
// _, err = q.client.RPush(config.Ctx, "worker_queue", job.ID).Result()
