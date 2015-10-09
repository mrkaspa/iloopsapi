package models

//Task wrapper json executed recurrently
type Task struct {
	ID          string `json:"id"`
	Periodicity string `json:"periodicity`
	Command     string `json:"command"`
}

type Email struct {
	Value string `json:"email" validate:"required"`
}

type ChangePassword struct {
	Password string `json:"password" validate:"required"`
	Token    string `json:"token" validate:"required"`
}
