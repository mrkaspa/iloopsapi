package gitadmin

import (
	"fmt"

	"bitbucket.org/kiloops/api/utils"
	"github.com/codeskyblue/go-sh"
)

var (
	ChanCommit   chan ChanReq
	ChanRollback chan ChanReq
)

type ChanReq struct {
	Path     string
	ChanResp *chan error
}

func InitGitAdmin() {

	ChanCommit = make(chan ChanReq)
	ChanRollback = make(chan ChanReq)

	fmt.Println("***INIT GIT ADMIN***")

	go func() {
		for {
			select {
			case req, ok := <-ChanCommit:
				if !ok {
					return
				}
				*req.ChanResp <- CommitChange(req.Path)
			case req, ok := <-ChanRollback:
				if !ok {
					return
				}
				*req.ChanResp <- RollbackChange(req.Path)
			}
		}
	}()
}

func GetCloseChanResponse(chanResp *chan error) error {
	err := <-*chanResp
	close(*chanResp)
	return err
}

func FinishGitAdmin() {
	fmt.Println("***Closing channels***")
	close(ChanCommit)
	close(ChanRollback)
}

func CommitChange(path string) error {
	if utils.IsTest() {
		return nil
	}
	session := sh.NewSession()
	return session.SetDir(path).Command("git", "pull", "origin", "master").Command("git", "add", "-A").Command("git", "commit", "-m", "update repo").Command("git", "push", "origin", "master").Run()
}

func RollbackChange(path string) error {
	if utils.IsTest() {
		return nil
	}
	session := sh.NewSession()
	return session.SetDir(path).Command("git", "reset", "--hard", "HEAD").Run()
}
