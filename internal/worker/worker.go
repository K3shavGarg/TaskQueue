package worker

import (
	"Job-Queue/internal/config"
	"Job-Queue/internal/model"
	"Job-Queue/metrics"
	"Job-Queue/pkg"

	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

var JobChan = make(chan *model.Job, 100)

func StartQueueProcessor(queue *model.RedisQueue, fetchers int, workers int) {
	for i := 0; i < fetchers; i++ {
		go queue.RedisFetcher(i, JobChan)
	}

	for i := 0; i < workers; i++ {
		go worker(i, JobChan, queue)
	}
}

func worker(id int, jobChan chan *model.Job, queue *model.RedisQueue) {
	for job := range jobChan {
		func() {
			startTime := time.Now()
			defer func() {
				if r := recover(); r != nil {
					pkg.Log.WithFields(logrus.Fields{
						"worker_id":   id,
						"job_id":      job.ID,
						"panic":       r,
						"duration_ms": time.Since(startTime).Milliseconds(),
					}).Error("Worker panicked while processing job")

					job.Status = "failed"
					job.CompletedAt = time.Now()
					_ = queue.SaveJob(job)
					_ = queue.FailJob(job.ID)
					pkg.Log.WithField("job_id", job.ID).Warn("Moved job to DLQ due to panic")
				}
			}()

			pkg.Log.WithFields(logrus.Fields{
				"worker_id":  id,
				"job_id":     job.ID,
				"job_type":   job.Type,
				"start_time": startTime,
				"attempts":   0,
			}).Info("Worker started processing task")

			err := processJob(job, id, queue)

			if err != nil {
				metrics.JobsFailed.Inc()
				pkg.Log.WithFields(logrus.Fields{
					"worker_id":   id,
					"job_id":      job.ID,
					"duration_ms": time.Since(startTime).Milliseconds(),
					"status":      "failed",
					"job_type":    job.Type,
					"attempts":    job.Attempts,
				}).WithError(err).Error("Job ultimately failed after max retries, moving to DLQ")
			} else {
				metrics.JobsProcessed.Inc()
				metrics.JobDuration.Observe(time.Since(startTime).Seconds())
				pkg.Log.WithFields(logrus.Fields{
					"worker_id":   id,
					"job_id":      job.ID,
					"duration_ms": time.Since(startTime).Milliseconds(),
					"status":      "success",
					"job_type":    job.Type,
					"attempts":    job.Attempts,
				}).Info("Worker completed processing task successfully")
			}
		}()
	}
}

func processJob(job *model.Job, worker_id int, queue *model.RedisQueue) error {
	delay := 1 * time.Second
	maxAttempts, valid := config.MaxAttempts[job.Type]
	if !valid {
		maxAttempts = 2
	}
	for attempts := 1; attempts <= maxAttempts; attempts++ {

		if job.Type == "panic" {
			// For panic testing
			panic("Gopher Paniked !!!")
		}
		// Call their respective Job function handlers...
		err := taskFunc(job)
		job.Attempts++

		if err == nil {
			job.Status = "completed"
			job.CompletedAt = time.Now()
			queue.SaveJob(job)
			queue.AckJob(job.ID) // removes from in_progress_queue
			return nil
		}

		if attempts == maxAttempts {
			job.Status = "failed"
			job.CompletedAt = time.Now()
			queue.SaveJob(job)
			queue.FailJob(job.ID) // move to DLQ
			return err
		}
		queue.SaveJob(job)

		pkg.Log.WithFields(logrus.Fields{
			"worker_id": worker_id,
			"job_id":    job.ID,
			"attempts":  attempts,
			"error":     err.Error(),
			"status":    "failure",
		}).Warn("Retrying Job...")

		jitter := time.Duration(rand.Int63n(int64(delay))) // Randomize the backoff
		time.Sleep(delay + jitter)
		delay *= 2 // Exponential backoff
	}
	return nil
}

func taskFunc(job *model.Job) error {
	num := rand.Float32()
	if num < 0.5 {
		return fmt.Errorf("task failed")
	}
	time.Sleep(1 * time.Second)
	return nil
}
