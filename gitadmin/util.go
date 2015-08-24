package gitadmin

import "github.com/codeskyblue/go-sh"

// import "github.com/codeskyblue/go-sh"

func CommitChange(path string) error {
	session := sh.NewSession()
	return session.SetDir(path).Command("git", "pull", "origin", "master").Command("git", "add", "-A").Command("git", "commit", "-m", "update repo").Command("git", "push", "origin", "master").Run()
}

func RevertChange(path string) error {
	session := sh.NewSession()
	return session.SetDir(path).Command("git", "reset", "--hard", "HEAD").Run()
}
