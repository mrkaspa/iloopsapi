package gitadmin

import "os"

var (
	GITOLITEPATH, GITURLROOT string
)

func InitVars() {
	GITOLITEPATH = os.Getenv("GITOLITE_PATH")
	GITURLROOT = os.Getenv("GIT_URL_ROOT")
}
