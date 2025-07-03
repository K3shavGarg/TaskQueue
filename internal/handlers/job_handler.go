package handlers

import (
	"Job-Queue/internal/model"
	"Job-Queue/internal/service"
	"Job-Queue/metrics"
	"Job-Queue/utils"
	"fmt"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

type JobHandler struct {
	jobService *service.JobService
}

func NewJobHandler(js *service.JobService) *JobHandler {
	return &JobHandler{jobService: js}
}

func (jh *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	status := query.Get("status")
	jobType := query.Get("type")
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	start := (page - 1) * limit
	end := start + limit - 1

	jobs, err := jh.jobService.GetAllJobs(int64(start), int64(end), status, jobType)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong in fetching jobs: %v", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, jobs)
}

func (jh *JobHandler) GetJobByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	job, err := jh.jobService.GetJobByID(params["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusFound, job)
}

func (jh *JobHandler) SubmitJob(w http.ResponseWriter, r *http.Request) {
	var payload model.JobRequestPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	job := model.Job{
		Payload: payload.Payload,
		Type:    payload.Type,
		Delay_ms: payload.Delay_ms,
		WebhookURL: payload.WebhookURL,
	}
	// Call the service
	if err := jh.jobService.PushJob(&job); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	metrics.JobsTotal.Inc()
	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"message": "Job Submitted Successfully",
		"job_id":  job.ID,
	})
}
