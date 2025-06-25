package model

import (
	"fmt"
	"sync"
)

type Queue struct {
	jobs       []Job
	mu         sync.Mutex
	maxSize    int
	jobHistory map[int64]*Job
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		jobs:       make([]Job, 0),
		maxSize:    capacity,
		jobHistory: make(map[int64]*Job),
		mu:         sync.Mutex{},
	}
}

func (q *Queue) Enqueue(job Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.jobs) >= q.maxSize {
		return fmt.Errorf("task queue reached its capacity")
	}
	q.jobs = append(q.jobs, job)
	q.jobHistory[job.ID] = &q.jobs[len(q.jobs)-1]
	return nil
}

func (q *Queue) Dequeue() (*Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.jobs) == 0 {
		return nil, fmt.Errorf("task queue has no jobs")
	}

	job := &q.jobs[0]
	q.jobs = q.jobs[1:] // remove the job from front
	return job, nil
}

func (q *Queue) GetJobs() []Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	jobs := make([]Job, 0, len(q.jobHistory))
	for _, job := range q.jobHistory {
		jobs = append(jobs, *job)
	}
	return jobs
}

func (q *Queue) GetJobByID(ID int64) (*Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	job, exists := q.jobHistory[ID]
	if !exists {
		return nil, fmt.Errorf("job not found")
	}
	return job, nil
}
