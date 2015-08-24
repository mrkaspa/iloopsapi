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

var TemplateProjectConf = "@users_%s = %s\nrepo %s\n  RW+ = @users_%s"

func CreateProject(slug string) error {
	path := ProjectPath(slug)
	if helpers.FileExists(path) {
		return ErrProjectFileExists
	}
	if _, err := os.Create(path); err == nil {
		return saveProjectFile(path, slug, &[]string{}, true)
	} else {
		return err
	}
}

func AddSSHToProject(email string, sshID int, slug string) error {
	if !helpers.FileExists(KeyPath(email, sshID)) {
		return ErrSSHFileNotFound
	}
	path := ProjectPath(slug)
	users := currentUsers(path)
	*users = append(*users, UserKeyValue(email, sshID))
	return saveProjectFile(path, slug, users, false)
}

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

func DeleteProject(slug string) error {
	path := ProjectPath(slug)
	if err := os.Remove(path); err == nil {
		return CommitChange(GITOLITEPATH)
	} else {
		return err
	}
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
	if err := ioutil.WriteFile(path, []byte(content), os.ModePerm); err == nil {
		if commit {
			return CommitChange(GITOLITEPATH)
		} else {
			return nil
		}
	} else {
		return err
	}
}

func ProjectPath(slug string) string {
	return GITOLITEPATH + "conf" + "/" + "repos" + "/" + slug + ".conf"
}
