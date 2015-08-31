package gitadmin

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/mrkaspa/go-helpers"
)

//AddSSH to gitolite
func AddSSH(email string, sshID int, content string) error {
	path := KeyPath(email, sshID)
	if helpers.FileExists(path) {
		return ErrSSHFileExists
	}
	if _, err := os.Create(path); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, []byte(content), os.ModePerm); err != nil {
		return err
	}
	chanResp := make(chan error)
	ChanCommit <- ChanReq{GITOLITEPATH, &chanResp}
	return GetCloseChanResponse(&chanResp)
}

//DeleteSSH from gitolite
func DeleteSSH(email string, sshID int) error {
	path := KeyPath(email, sshID)
	if err := os.Remove(path); err != nil {
		return err
	}
	chanResp := make(chan error)
	ChanCommit <- ChanReq{GITOLITEPATH, &chanResp}
	return GetCloseChanResponse(&chanResp)
}

//KeyPath generator
func KeyPath(email string, sshID int) string {
	return GITOLITEPATH + "keydir" + "/" + UserKeyValue(email, sshID) + ".pub"
}

//UserKeyValue generator
func UserKeyValue(email string, sshID int) string {
	return email + "-" + strconv.Itoa(sshID)
}
