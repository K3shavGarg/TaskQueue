package handlers

import (
	"Job-Queue/internal/model"
	"Job-Queue/internal/service"
	"Job-Queue/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type JobHandler struct {
	jobService *service.JobService
}

func NewJobHandler(js *service.JobService) *JobHandler {
	return &JobHandler{jobService: js}
}

func (jh *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusAccepted, jh.jobService.GetAllJobs())
}

func (jh *JobHandler) GetJobByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID, err := strconv.Atoi(params["id"])

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id field is neither not present or not numeric"))
		return
	}

	job, err := jh.jobService.GetJobByID(int64(ID))
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
		ID:        time.Now().UnixMilli(), // Simple unique ID
		Attempts:     0,
		CreatedAt: time.Now(),
		Payload:   payload.Payload,
		Status:    "pending",
		Type: 	payload.Type,
	}
	// Call the service
	if err := jh.jobService.PushJob(job); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, map[string]any{
		"message": "Job Submitted Successfully",
		"job_id":  job.ID,
	})
}
