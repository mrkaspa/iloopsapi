package gitadmin

import "errors"

var (
	ErrSSHFileExists     = errors.New("The SSH file already exists")
	ErrSSHFileNotFound   = errors.New("The SSH file does not exist")
	ErrProjectFileExists = errors.New("The Project file already exists")
)
