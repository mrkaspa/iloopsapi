package ierrors

const (
	_ = iota
	ErrCodeGeneral
	ErrCodeValidation
	ErrCodeBadInput
	ErrCodeNotAccess
	ErrCodeNotFound
	ErrCodeCredential
	ErrCodeUserInactive
	ErrCodeCreate
	ErrCodeDelete
	ErrCodeAdminLeaveProject
	ErrCodeUserLeaveProject
	ErrCodeExceedMaxProjects
	ErrCodeProjectStart
	ErrCodeProjectStop
)

var (
	ErrProjectCreate         = AppError{Code: ErrCodeCreate, ErrorS: "Could not create the project"}
	ErrProjectDelete         = AppError{Code: ErrCodeDelete, ErrorS: "Could not delete the project"}
	ErrAdminCantLeaveProject = AppError{Code: ErrCodeAdminLeaveProject, ErrorS: "An admin user can't leave a project"}
	ErrUserLeaveProject      = AppError{Code: ErrCodeUserLeaveProject, ErrorS: "Could not leave the project"}
	ErrProjectAddUser        = AppError{Code: ErrCodeGeneral, ErrorS: "Could not add the user"}
	ErrProjectRemoveUser     = AppError{Code: ErrCodeGeneral, ErrorS: "Could not remove the user"}
	ErrProjectDelegateUser   = AppError{Code: ErrCodeGeneral, ErrorS: "Could not delegate the project to the user"}
	ErrSSHCreate             = AppError{Code: ErrCodeCreate, ErrorS: "Could not create the SSH"}
	ErrSSHDelete             = AppError{Code: ErrCodeDelete, ErrorS: "Could not delete the SSH"}
	ErrUserCreate            = AppError{Code: ErrCodeCreate, ErrorS: "Could not create the user"}
	ErrUserLogin             = AppError{Code: ErrCodeCredential, ErrorS: "Could not authenticate the User"}
	ErrUserInactive          = AppError{Code: ErrCodeUserInactive, ErrorS: "The user is inactive"}

	ErrTaskNotScheduled      = AppError{Code: ErrCodeProjectStart, ErrorS: "The task can't be scheduled"}
	ErrTaskNotStopped        = AppError{Code: ErrCodeProjectStop, ErrorS: "The task can't be stopped"}
	ErrUserExceedMaxProjects = AppError{Code: ErrCodeExceedMaxProjects, ErrorS: "The user can't have more projects associated"}
	ErrUserProjectNotSaved   = AppError{Code: ErrCodeCreate, ErrorS: "User Project can't be saved"}
	ErrProjectNotSaved       = AppError{Code: ErrCodeCreate, ErrorS: "Project can't be saved"}
	ErrUserNotFound          = AppError{Code: ErrCodeNotFound, ErrorS: "User not found"}
	ErrCreatorNotRemoved     = AppError{Code: ErrCodeDelete, ErrorS: "You can't remove a Creator from a project"}
	ErrUserIsNotCollaborator = AppError{Code: ErrCodeNotAccess, ErrorS: "The user doesn't have collaborator access to the project"}
	ErrProjectNotFound       = AppError{Code: ErrCodeNotFound, ErrorS: "Project not found"}
)
