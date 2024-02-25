package entity

type TaskRequest struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type Task struct {
	Id             string
	Status         string
	HttpStatusCode int
	Headers        map[string][]string
	Length         int
}
