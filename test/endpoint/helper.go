package endpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"bitbucket.org/kiloops/api/models"
)

type empty struct{}

var emptyJSON, _ = json.Marshal(empty{})

//Client for http requests
type Client struct {
	*http.Client
	baseURL     string
	contentType string
}

func (c Client) CallRequest(method string, path string, reader io.Reader) (*http.Response, error) {
	return c.CallRequestWithHeaders(method, path, reader, make(map[string]string))
}

func (c Client) CallRequestWithHeaders(method string, path string, reader io.Reader, headers map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest(method, c.baseURL+path, reader)
	req.Header.Set("Content-Type", c.contentType)
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	return c.Do(req)
}

func authHeaders(user models.User) map[string]string {
	return map[string]string{
		"AUTH_ID":    fmt.Sprintf("%d", user.ID),
		"AUTH_TOKEN": user.Token,
	}
}

func saveUser() models.User {
	user := models.User{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
	models.Gdb.InTx(func(txn *models.KDB) {
		if txn.Save(&user).Error != nil {
			panic("error creating the user")
		}
	})
	return user
}

func saveOtherUser() models.User {
	user := models.User{Email: "angelbotto@gmail.com", Password: "h1h1h1h1h1h1"}
	models.Gdb.InTx(func(txn *models.KDB) {
		if txn.Save(&user).Error != nil {
			panic("error creating the user")
		}
	})
	return user
}

func addSSH(user models.User) models.SSH {
	sshKey := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDCadJM0DdJotRnSWW7coFcCxMW1cCZIJqkyW3wMQoUOU2VHuLExh44tDpSAiz2EEeFlqJ5hI67ZI+3bSx7puKSr44l78H/Kb8UDLidAUao7JZoo0thq7bAVesGr+8aligmULvxH3sQqstI9yNcifJ56jHUVTB14PslBmhA56pmGOva0ojmdt9l2aBy4LxQBDc5Js+AcPlfC2zXE7rtaiafB/M3992V+7CEisbAv7CpsI3SPdpW2p4mfR1zMVpf4Jt6lQJW6Sr53/bzAP4/Tif3fgbZhoSL8qnnLi3556gWi90FwFhCoqqDR/lN3sxJQx5NxxCF8mbNgpmS5qDptFyF michel.ingesoft@gmail.com`
	ssh := models.SSH{PublicKey: sshKey, UserID: user.ID}
	models.Gdb.InTx(func(txn *models.KDB) {
		if txn.Save(&ssh).Error != nil {
			panic("error creating the ssh")
		}
	})
	return ssh
}

func addProject(user models.User) models.Project {
	project := models.Project{Name: "Demo Project"}
	models.Gdb.InTx(func(txn *models.KDB) {
		if user.CreateProject(txn, &project) != nil {
			panic("error creating the project")
		}
	})
	return project
}

func getBodyJSON(resp *http.Response, i interface{}) {
	defer resp.Body.Close()
	if jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body); err == nil {
		if err := json.Unmarshal([]byte(jsonDataFromHTTP), &i); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
}

func debugResponse(resp *http.Response) {
	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("*****************")
	fmt.Println(string(contents))
	fmt.Println("*****************")
}
