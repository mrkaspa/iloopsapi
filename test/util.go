package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"

	"bitbucket.org/kiloops/api/models"
)

var emptyJSON, _ = json.Marshal(struct{}{})

func authHeaders(user models.User) map[string]string {
	return map[string]string{
		"AUTH_ID":    fmt.Sprintf("%d", user.ID),
		"AUTH_TOKEN": user.Token,
	}
}

func defaultUser() models.User {
	return models.User{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
}

func saveUser() models.User {
	user := defaultUser()
	models.InTx(func(txn *gorm.DB) bool {
		if txn.Create(&user).Error != nil {
			panic("error creating the user")
		}
		return true
	})
	return user
}

func saveOtherUser() models.User {
	user := models.User{Email: "angelbotto@gmail.com", Password: "h1h1h1h1h1h1"}
	models.InTx(func(txn *gorm.DB) bool {
		if txn.Create(&user).Error != nil {
			panic("error creating the user")
		}
		return true
	})
	return user
}

func addSSH(user models.User) models.SSH {
	sshKey := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDlCc96zWY05/vFIcP5NLhi8bIVkcUdSyet1Dw7+rQqbeEJaQ0Ifz/x17AGkQAnC0ZjdHI7sCFjVGuvk6agw6MJzKU8a+iWisAVu4hvv22DXBPKYak28GMEW3e0Ba/8mUiCdLCW5lfQ85QmDABqdWb6BGy2VSJ/k4NfWW728RwbQf1MZSwS2+kqvR3XjpkpMETLz5DmRty6Dqp3al73JbE7raWhidoYeS0wiJKsWiaucfewz+feubNkEnO5/p1v1zpAlaPYEVvZEeG5ABchNZ4Co+SGvVd4+FuxVgLkPOqpV5y3JFFrmSJE4HMsin96u/3OHcgVwew6nyE9YyoKZ/rL michel.ing@hotmail.com`
	ssh := models.SSH{Name: "demo", PublicKey: sshKey, UserID: user.ID}
	models.InTx(func(txn *gorm.DB) bool {
		if txn.Create(&ssh).Error != nil {
			panic("error creating the ssh")
		}
		return true
	})
	return ssh
}

func addAnotherSSH(user models.User) models.SSH {
	sshKey := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDBZ5qb4XIWo4yQbRJWqjSdJs/9TYs8QranpnUOoHQDd/6Teik4NptSJ8QaMpF+9p6euscG7ceIRvBUFvc0Cecy/E7uVoQVm/BosypvbrywvTcnwNkHeNJftu05wk4iUj4AYol57zXxYXH0hqq/UD7ijSlsG/d24n9NR9r3Ocng1BWpjd+ZdJFqzYLwp/1vVMOohVdOKtJymnkjFWBoEqaZ3g4p8glb/NrjC0154r6vLq3FBLglEqCdcXXX6dy0QFEtGPHrqAGVj10vnJYUsPrGjZzOHYHxQwvzdtmtI8lMSoPEHZ39ODrABzfv7b07bT4it9YVFodUUyfNn+bZYp/F mrkaspa@github.com`
	ssh := models.SSH{Name: "demo", PublicKey: sshKey, UserID: user.ID}
	models.InTx(func(txn *gorm.DB) bool {
		if txn.Create(&ssh).Error != nil {
			panic("error creating the ssh")
		}
		return true
	})
	return ssh
}

func addProject(user models.User) models.Project {
	project := models.Project{Name: "Demo Project"}
	models.InTx(func(txn *gorm.DB) bool {
		if err := user.CreateProject(txn, &project); err != nil {
			panic("error creating the project")
		}
		return true
	})
	return project
}

func getBodyJSON(resp *http.Response, i interface{}) {
	if jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body); err == nil {
		if err := json.Unmarshal([]byte(jsonDataFromHTTP), &i); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
}

func debugResponse(resp *http.Response) {
	contents, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("*****************")
	fmt.Println(string(contents))
	fmt.Println("*****************")
}
