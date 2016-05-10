package gitadmin

import (
	"github.com/mrkaspa/iloopsapi/utils"
	"github.com/codeskyblue/go-sh"
)

var (
	ChanCommit   chan ChanReq
	ChanRollback chan ChanReq
)

//ChanReq represent the requests to the channels
type ChanReq struct {
	Path     string
	ChanResp *chan error
}

//InitGitAdmin go routine
func InitGitAdmin() {

	ChanCommit = make(chan ChanReq)
	ChanRollback = make(chan ChanReq)

	utils.Log.Info("***INIT GIT ADMIN***")

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

//GetCloseChanResponse gets the response and closes the channel
func GetCloseChanResponse(chanResp *chan error) error {
	err := <-*chanResp
	close(*chanResp)
	return err
}

//FinishGitAdmin go routine
func FinishGitAdmin() {
	utils.Log.Info("***Closing Git channels***")
	close(ChanCommit)
	close(ChanRollback)
}

//CommitChange to the master
func CommitChange(path string) error {
	if utils.IsTest() {
		return nil
	}
	session := sh.NewSession()
	session.SetDir(path)
	session.Command("git", "pull", "origin", "master").Run()
	session.Command("git", "add", "-A").Run()
	session.Command("git", "commit", "-m", "update repo").Run()
	return session.Command("git", "push", "origin", "master").Run()
}

//RollbackChange from the gitolite repo
func RollbackChange(path string) error {
	if utils.IsTest() {
		return nil
	}
	session := sh.NewSession()
	session.SetDir(path).Command("git", "reset", "--hard", "HEAD").Run()
	return session.SetDir(path).Command("git", "clean", "-f").Run()
}

//RevertAll changes of the gitolite repo
func RevertAll(path string) error {
	utils.Log.Info("***RevertAll***")
	if utils.IsTest() {
		session := sh.NewSession()
		return session.SetDir(path).Command("git", "clean", "-f").Run()
	}
	return nil
}

//InTx executes a task on the repo inside a git transaction
func InTx(f func() error) error {
	if err := f(); err != nil {
		chanResp := make(chan error)
		ChanRollback <- ChanReq{GITOLITEPATH, &chanResp}
		GetCloseChanResponse(&chanResp)
		return err
	}
	chanResp := make(chan error)
	ChanCommit <- ChanReq{GITOLITEPATH, &chanResp}
	return GetCloseChanResponse(&chanResp)
}
