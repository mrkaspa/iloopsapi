package models

//Task executed recurrently
type Task struct {
	ID          string `json:"id"`
	Periodicity string `json:"periodicity`
	Command     string `json:"command"`
}
