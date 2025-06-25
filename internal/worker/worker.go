package worker

import (
	"Job-Queue/internal/config"
	"Job-Queue/internal/model"
	"Job-Queue/pkg"

	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

func StartWorkerPool(queue *model.Queue, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go worker(i, queue)
	}
}

func worker(id int, queue *model.Queue) {
	for {
		job, err := queue.Dequeue()
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		job.Status = "processing"

		pkg.Log.WithFields(logrus.Fields{
			"worker_id": id,
			"job_id":    job.ID,
			"status":    "processing",
			"job_type":  job.Type,
		}).Info("Worker started processing task")

		// Simulate job processing
		startTime := time.Now()

		err = processJob(job, id)

		if err != nil {
			pkg.Log.WithFields(logrus.Fields{
				"worker_id": id,
				"job_id":    job.ID,
				"duration":  time.Since(startTime).Milliseconds(),
				"status":    "failed",
			}).WithError(err).Error("Worker failed to process task")
			return
		}

		pkg.Log.WithFields(logrus.Fields{
			"worker_id": id,
			"job_id":    job.ID,
			"duration":  time.Since(startTime).Milliseconds(),
			"status":    "success",
			"job_type":  job.Type,
		}).Info("Worker completed processing task successfully")

		job.Status = "completed"
	}
}

func processJob(job *model.Job, worker_id int) error {
	delay := 1 * time.Second
	maxAttempts, valid := config.MaxAttempts[job.Type]
	if !valid {
		maxAttempts = 2
	}
	for attempts := 1; attempts <= maxAttempts; attempts++ {

		// Call their respective Job function handlers...
		job.Attempts++
		err := taskFunc(job)

		if err == nil {
			return nil
		}

		if attempts == maxAttempts {
			job.Status = "failed"
			job.CompletedAt = time.Now()
			return err
		}
		pkg.Log.WithFields(logrus.Fields{
			"worker_id": worker_id,
			"job_id":    job.ID,
			"attempts":  attempts,
			"error":     err.Error(),
			"status":    "retrying",
		}).Warn("Retrying Job...")

		jitter := time.Duration(rand.Int63n(int64(delay))) // Randomize the backoff
		time.Sleep(delay + jitter)
		delay *= 2 // Exponential backoff
	}
	return nil
}

func taskFunc(job *model.Job) error {
	num := rand.Float32()
	fmt.Println(num)
	if num < 0.5 {
		return fmt.Errorf("task failed")
	}
	job.CompletedAt = time.Now()
	return nil
}
