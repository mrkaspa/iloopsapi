package models

import "errors"

var (
	ErrUserProjectNotSaved   = errors.New("User Project can't be saved")
	ErrProjectNotSaved       = errors.New("Project can't be saved")
	ErrUserNotFound          = errors.New("User not found")
	ErrCreatorNotRemoved     = errors.New("You can't remove a Creator from a project")
	ErrUserIsNotCollaborator = errors.New("The user doesn't have collaborator access to the project")
	ErrProjectNotFound       = errors.New("Project not found")
	ErrTaskNotScheduled      = errors.New("The task can't be scheduled")
)
