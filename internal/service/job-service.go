package service

import (
	"Job-Queue/internal/model"
	"fmt"
)

type JobService struct {
	Queue *model.RedisQueue
}

func NewJobService(queue *model.RedisQueue) *JobService {
	return &JobService{
		Queue: queue,
	}
}

func (j *JobService) GetAllJobs(start, end int64, status, jobType string) ([]*model.Job, error) {
	cmdMap, err := j.Queue.GetJobs(start, end)
	if err != nil {
		return nil, err
	}
	jobs := []*model.Job{}
	for id, cmd := range cmdMap {
		data := cmd.Val()
		if len(data) == 0 {
			continue
		}

		job := model.HydrateJob(data, id)

		if status != "" && job.Status != status {
			continue
		}
		if jobType != "" && job.Type != jobType {
			continue
		}

		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (j *JobService) GetJobByID(ID string) (*model.Job, error) {
	job, err := j.Queue.GetJobByID(ID)

	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j *JobService) PushJob(job *model.Job) error {
	if err := j.Queue.Enqueue(job); err != nil {
		return fmt.Errorf("error in enqueueing: %v", err)
	}
	return nil
}
