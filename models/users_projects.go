package models

import "time"

const (
	_ = iota
	Creator
	Collaborator
)

//UsersProjects ManyToMany rel
type UsersProjects struct {
	ID        int     `gorm:"primary_key" json:"id"`
	Role      int     `json:"role"`
	ProjectID int     `json:"project_id"`
	UserID    int     `json:"user_id"`
	Project   Project `json:"project"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//TableName for UsersProjects
func (u UsersProjects) TableName() string {
	return "users_projects"
}
