package entity

const (
	TaskStatusNew       = "new"
	TaskStatusInProcess = "in_process"
	TaskStatusDone      = "done"
	TaskStatusError     = "error"
)

type TaskRequest struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    any               `json:"body"`
}

type Task struct {
	Id             string              `json:"-"`
	Status         string              `json:"status"`
	HttpStatusCode int                 `json:"http_status_code"`
	Headers        map[string][]string `json:"headers"`
	Length         int64               `json:"length"`
}
