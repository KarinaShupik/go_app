package requests

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	Deadline    *uint64 `json:"deadline"`
}

func (r TaskRequest) ToDomainModel() (interface{}, error) {

	var deadline *time.Time
	if r.Deadline != nil {
		ddl := time.Unix(int64(*r.Deadline), 0)
		deadline = &ddl
	}

	/*descr := ""
	if r.Description != nil {
		descr = *r.Description
	}

	var deadline uint64
	if r.Deadline != nil {
		deadline = *r.Deadline
	}*/

	return domain.Task{
		Title:       r.Title,
		Description: r.Description,
		Deadline:    deadline,
	}, nil
}
