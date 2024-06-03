package controllers

import (
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id
		task.Status = domain.NewTakStatus
		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Created(w, tDto)
	}
}

func (c TaskController) FindByUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		tasks, err := c.taskService.FindByUserId(user.Id)
		if err != nil {
			log.Printf("TaskController -> FindBuUserId: %s", err)
			InternalServerError(w, err)
			return
		}

		var tsDto resources.TasksDto
		tsDto = tsDto.DomainToDtoCollection(tasks)
		Success(w, tsDto)
	}
}

/*
func GetTaskIdFromRequest(r *http.Request) (uint64, error) {
	// Extract the path from the request
	urlPath := r.URL.Path

	// Split the path by "/" to get segments
	pathSegments := strings.Split(urlPath, "/")

	// Check if there are enough segments (base_url + /tasks/ + ID)
	if len(pathSegments) < 3 {
		return 0, errors.New("Invalid URL format for task ID")
	}

	// The third segment should be the ID
	taskIdStr := pathSegments[2]

	// Convert the ID string to uint64
	taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
	if err != nil {
		return 0, errors.New("Invalid task ID format")
	}

	return taskId, nil
}*/
/*
func (c TaskController) FindByTaskId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract task ID from the URL
		taskIdStr := chi.URLParam(r, "taskId") // Access the "taskId" path parameter

		// Convert the ID string to uint64
		taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
		// Call task service to find task
		task, err := c.taskService.FindByTaskId(taskId)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) { // Handle specific task not found error
				NotFound(w, "Task not found")
				return
			}
			log.Printf("TaskController -> FindById: %s", err)
			InternalServerError(w, err)
			return
		}

		// Convert task to DTO and respond (assuming you want the full task data)
		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Success(w, tDto)
	}
}*/

func (c TaskController) FindByTaskId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId, err := requests.ParseTaskId(r)
		if err != nil {
			log.Printf("TaskController -> FindByTaskId: %s", err)
			BadRequest(w, err)
			return
		}

		task, err := c.taskService.FindByTaskId(taskId)
		if err != nil {
			log.Printf("TaskController -> FindByTaskId: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Success(w, tDto)
	}
}

func (c TaskController) DeleteByTaskId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId, err := requests.ParseTaskId(r)
		if err != nil {
			log.Printf("TaskController -> DeleteByTaskId: %s", err)
			BadRequest(w, err)
			return
		}

		err = c.taskService.DeleteByTaskId(taskId)
		if err != nil {
			log.Printf("TaskController -> DeleteByTaskId: %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}

/*func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func SuccessUpdate(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}*/

func (c TaskController) UpdateByTaskId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId, err := requests.ParseTaskId(r)
		if err != nil {
			log.Printf("TaskController -> UpdateByTaskId: %s", err)
			BadRequest(w, err)
			return
		}

		taskReq, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> UpdateByTaskId: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		taskReq.UserId = user.Id
		taskReq.Id = taskId

		updatedTask, err := c.taskService.UpdateByTaskId(taskReq)
		if err != nil {
			log.Printf("TaskController -> UpdateByTaskId: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(updatedTask)
		Success(w, tDto)
	}
}
