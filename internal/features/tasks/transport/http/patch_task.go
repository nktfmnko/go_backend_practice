package tasks_transport_http

import (
	"fmt"
	"net/http"
	"practice/internal/core/domain"
	core_logger "practice/internal/core/logger"
	core_http_request "practice/internal/core/transport/http/request"
	core_http_response "practice/internal/core/transport/http/response"
	core_http_types "practice/internal/core/transport/http/types"
)

type PatchTaskResponse TaskDTOResponse

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("title cant be null")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("title must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("description must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("completed cant be null")
		}
	}
	return nil
}

func (h *TasksHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	const (
		taskIDPath = "id"
	)
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	taskID, err := core_http_request.GetIntPathValue(r, taskIDPath)
	if err != nil {
		responseHandler.ErrorResponse(err, "get id from path")
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "validate request")
		return
	}

	taskPatch := taskPatchFromRequest(request)
	taskPatchDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)

	if err != nil {
		responseHandler.ErrorResponse(err, "patch task")
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskPatchDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
