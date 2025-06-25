package service

import (
	"Job-Queue/internal/model"
	"fmt"
)

type JobService struct {
	Queue *model.Queue
}

func NewJobService() *JobService {
	NewQueue := model.NewQueue(10)
	return &JobService{
		Queue: NewQueue,
	}
}

func (j *JobService) GetAllJobs() []model.Job {
	return j.Queue.GetJobs()
}

func (j *JobService) GetJobByID(ID int64) (*model.Job, error) {
	job, err := j.Queue.GetJobByID(ID)

	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j *JobService) PushJob(job model.Job) error {
	if err := j.Queue.Enqueue(job); err != nil {
		return fmt.Errorf("error in enqueueing: %v", err)
	}
	return nil
}
