package models

import "time"

//Project on the system
type Project struct {
	ID     int    `gorm:"primary_key" json:"id"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//DeleteRels UsersProjects
func (p *Project) DeleteRels(txn *KDB) {
	txn.Where("project_id = ?", p.ID).Delete(UsersProjects{})
}
