package gitadmin

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mrkaspa/go-helpers"
)

//TemplateProjectConf for user access
var TemplateProjectConf = "@users_%s = %s\nrepo %s\n  RW+ = @users_%s"

//CreateProject git file
func CreateProject(slug string) error {
	path := ProjectPath(slug)
	if helpers.FileExists(path) {
		return ErrProjectFileExists
	}
	if _, err := os.Create(path); err != nil {
		return err
	}
	return saveProjectFile(path, slug, &[]string{}, true)
}

//AddSSHToProject file key to project
func AddSSHToProject(email string, sshID int, slug string) error {
	if !helpers.FileExists(KeyPath(email, sshID)) {
		return ErrSSHFileNotFound
	}
	path := ProjectPath(slug)
	users := currentUsers(path)
	*users = append(*users, UserKeyValue(email, sshID))
	return saveProjectFile(path, slug, users, false)
}

//RemoveSSHFromProject file key to project
func RemoveSSHFromProject(email string, sshID int, slug string) error {
	if !helpers.FileExists(KeyPath(email, sshID)) {
		return ErrSSHFileNotFound
	}
	path := ProjectPath(slug)
	users := currentUsers(path)
	key := UserKeyValue(email, sshID)
	usersFiltered := []string{}
	for _, v := range *users {
		if v != key {
			usersFiltered = append(usersFiltered, v)
		}
	}
	return saveProjectFile(path, slug, &usersFiltered, false)
}

//DeleteProject git file
func DeleteProject(slug string) error {
	path := ProjectPath(slug)
	if err := os.Remove(path); err != nil {
		return err
	}
	chanResp := make(chan error)
	ChanCommit <- ChanReq{GITOLITEPATH, &chanResp}
	return GetCloseChanResponse(&chanResp)
}

func currentUsers(path string) *[]string {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		s1 := strings.Split(line, "=")
		if len(s1) == 2 {
			usersCad := strings.TrimSpace(s1[1])
			s2 := strings.Split(usersCad, " ")
			return &s2
		}
		break
	}
	return &[]string{}
}

func saveProjectFile(path string, slug string, users *[]string, commit bool) error {
	var usersBuff bytes.Buffer
	for _, user := range *users {
		usersBuff.WriteString(user + " ")
	}
	content := fmt.Sprintf(TemplateProjectConf, slug, strings.TrimSpace(usersBuff.String()), slug, slug)
	if err := ioutil.WriteFile(path, []byte(content), os.ModePerm); err != nil {
		return err
	}
	if !commit {
		return nil
	}
	chanResp := make(chan error)
	ChanCommit <- ChanReq{GITOLITEPATH, &chanResp}
	return GetCloseChanResponse(&chanResp)
}

//ProjectPath generator
func ProjectPath(slug string) string {
	return GITOLITEPATH + "conf" + "/" + "repos" + "/" + slug + ".conf"
}
