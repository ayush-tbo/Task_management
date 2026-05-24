package dto

import "github.com/floqast/task-management/backend/internal/model"

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type StatusChartEntry struct {
	Status model.TaskStatus `json:"status"`
	Count  int              `json:"count"`
}

type StatusChart struct {
	ProjectID string             `json:"project_id"`
	Data      []StatusChartEntry `json:"data"`
}

type PriorityChartEntry struct {
	Priority model.Priority `json:"priority"`
	Count    int            `json:"count"`
}

type PriorityChart struct {
	ProjectID string               `json:"project_id"`
	Data      []PriorityChartEntry `json:"data"`
}
