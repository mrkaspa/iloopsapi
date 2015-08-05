package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"

	"bitbucket.org/kiloops/api/models"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Users", func() {

	Describe("POST /users", func() {

		It("create an user", func() {
			// Around(func(c *Client) {
			userLogin := models.UserLogin{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
			userJSON, _ := json.Marshal(userLogin)
			resp, _ := client.CallRequest("POST", "/users", bytes.NewReader(userJSON))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			// })
		})

	})

	Describe("POST /users/login", func() {
		var user models.User

		BeforeEach(func() {
			user = saveUser()
		})

		It("login successfully", func() {
			// Around(func(c *Client) {
			userLogin := models.UserLogin{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
			userJSON, _ := json.Marshal(userLogin)
			resp, _ := client.CallRequest("POST", "/users/login", bytes.NewReader(userJSON))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			// })
		})

		It("fails login incorrect password", func() {
			// Around(func(c *Client) {
			userLogin := models.UserLogin{Email: "michel.ingesoft@gmail.com", Password: "h1"}
			userJSON, _ := json.Marshal(userLogin)
			resp, _ := client.CallRequest("POST", "/users/login", bytes.NewReader(userJSON))
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			// })
		})

	})

})
