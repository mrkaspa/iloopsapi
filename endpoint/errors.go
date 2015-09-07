package endpoint

import "errors"

type JSONError struct {
	Code      int
	Error     string
	MapErrors map[string]string
}

type AppError struct {
	Code  int
	Error error
}

const (
	_ = iota
	GeneralErrCode
	ValidationErrCode
	BadInputErrCode
	CreateErrCode
	DeleteErrCode
	AdminLeaveProjectErrCode
	UserLeaveProjectErrCode
)

var (
	ProjectCreateErr         = AppError{CreateErrCode, errors.New("Could not create the project")}
	ProjectDeleteErr         = AppError{DeleteErrCode, errors.New("Could not delete the project")}
	AdminCantLeaveProjectErr = AppError{AdminLeaveProjectErrCode, errors.New("An admin user can't leave a project")}
	UserLeaveProjectErr      = AppError{UserLeaveProjectErrCode, errors.New("Could not leave the project")}
	ProjectAddUserErr        = AppError{GeneralErrCode, errors.New("Could not add the user")}
	ProjectRemoveUserErr     = AppError{GeneralErrCode, errors.New("Could not remove the user")}
	ProjectDelegateUserErr   = AppError{GeneralErrCode, errors.New("Could not delegate the project to the user")}
	SSHCreateErr             = AppError{CreateErrCode, errors.New("Could not create the SSH")}
	SSHDeleteErr             = AppError{DeleteErrCode, errors.New("Could not delete the SSH")}
	UserCreateErr            = AppError{CreateErrCode, errors.New("Could not create the user")}
	UserLoginErr             = AppError{GeneralErrCode, errors.New("Could not authenticate the User")}
)
