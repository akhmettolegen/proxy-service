package entity

type Task struct {
	Id             string
	Status         string
	HttpStatusCode int
	Headers        map[string][]string
	Length         int
}
