package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jinzhu/gorm"

	"bitbucket.org/kiloops/api/models"
	"bitbucket.org/kiloops/api/test"
)

func TestDemo(t *testing.T) {
	test.Around(func(c *test.Client) {
		userLogin := models.UserLogin{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
		userJSON, _ := json.Marshal(userLogin)
		resp, _ := c.CallRequest("POST", "/users", bytes.NewReader(userJSON))
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status >> %s", resp.Status)
		}
	})
}

func TestLogin(t *testing.T) {
	saveUser()
	test.Around(func(c *test.Client) {
		userLogin := models.UserLogin{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
		userJSON, _ := json.Marshal(userLogin)
		resp, _ := c.CallRequest("POST", "/users/login", bytes.NewReader(userJSON))
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status >> %s", resp.Status)
		}
	})
}

func TestLoginWrongParams(t *testing.T) {
	saveUser()
	test.Around(func(c *test.Client) {
		userLogin := models.UserLogin{Email: "michel.ingesoft@gmail.com", Password: "h1"}
		userJSON, _ := json.Marshal(userLogin)
		resp, _ := c.CallRequest("POST", "/users/login", bytes.NewReader(userJSON))
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status >> %s", resp.Status)
		}
	})
}

func saveUser() models.User {
	user := models.User{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
	models.Gdb.InTx(func(txn *gorm.DB) {
		txn.Save(&user)
	})
	return user
}
