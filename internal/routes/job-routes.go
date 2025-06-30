package routes

import (
	"Job-Queue/internal/model"
	"Job-Queue/internal/handlers"
	"Job-Queue/internal/service"

	"github.com/gorilla/mux"
)

var JobService *service.JobService

func RegisterJobRoutes(router *mux.Router, queue *model.RedisQueue) {
	JobService = service.NewJobService(queue)
	jobHandler := handlers.NewJobHandler(JobService)
	router.HandleFunc("/submit-job", jobHandler.SubmitJob).Methods("POST")
	router.HandleFunc("/jobs", jobHandler.GetAllJobs).Methods("GET")
	router.HandleFunc("/job/{id}", jobHandler.GetJobByID).Methods("GET")
}
