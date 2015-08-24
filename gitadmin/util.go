package gitadmin

import (
	"bitbucket.org/kiloops/api/utils"
	"github.com/codeskyblue/go-sh"
)

func CommitChange(path string) error {
	if utils.IsTest() {
		return nil
	}
	session := sh.NewSession()
	return session.SetDir(path).Command("git", "pull", "origin", "master").Command("git", "add", "-A").Command("git", "commit", "-m", "update repo").Command("git", "push", "origin", "master").Run()
}

func RevertChange(path string) error {
	if utils.IsTest() {
		return nil
	}
	session := sh.NewSession()
	return session.SetDir(path).Command("git", "reset", "--hard", "HEAD").Run()
}
