package endpoint

import (
	"bytes"
	"errors"
	"fmt"
)

type JSONError struct {
	Code      int
	ErrorCad  string
	MapErrors map[string]string
}

func (j JSONError) Error() string {
	if j.ErrorCad != "" {
		return j.ErrorCad
	}
	var buffer bytes.Buffer
	buffer.WriteString("Validation errors:\n")
	for key, value := range j.MapErrors {
		buffer.WriteString(fmt.Sprintf("Field validation for '%s' failed on the field '%s'", key, value))
	}
	return buffer.String()
}

type AppError struct {
	Code  int
	Error error
}

const (
	_ = iota
	ErrCodeGeneral
	ErrCodeValidation
	ErrCodeBadInput
	ErrCodeCredential
	ErrCodeActive
	ErrCodeCreate
	ErrCodeDelete
	ErrCodeAdminLeaveProject
	ErrCodeUserLeaveProject
)

var (
	ErrProjectCreate         = AppError{ErrCodeCreate, errors.New("Could not create the project")}
	ErrProjectDelete         = AppError{ErrCodeDelete, errors.New("Could not delete the project")}
	ErrAdminCantLeaveProject = AppError{ErrCodeAdminLeaveProject, errors.New("An admin user can't leave a project")}
	ErrUserLeaveProject      = AppError{ErrCodeUserLeaveProject, errors.New("Could not leave the project")}
	ErrProjectAddUser        = AppError{ErrCodeGeneral, errors.New("Could not add the user")}
	ErrProjectRemoveUser     = AppError{ErrCodeGeneral, errors.New("Could not remove the user")}
	ErrProjectDelegateUser   = AppError{ErrCodeGeneral, errors.New("Could not delegate the project to the user")}
	ErrSSHCreate             = AppError{ErrCodeCreate, errors.New("Could not create the SSH")}
	ErrSSHDelete             = AppError{ErrCodeDelete, errors.New("Could not delete the SSH")}
	ErrUserCreate            = AppError{ErrCodeCreate, errors.New("Could not create the user")}
	ErrUserLogin             = AppError{ErrCodeCredential, errors.New("Could not authenticate the User")}
	ErrUserActive            = AppError{ErrCodeActive, errors.New("The user is inactive")}
)
