package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/kiloops/api/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Users", func() {

	BeforeEach(func() {
		cleanDB()
	})

	Describe("POST /users", func() {

		It("create an user", func() {
			userLogin := models.UserLogin{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
			userJSON, _ := json.Marshal(userLogin)
			resp, _ := client.CallRequest("POST", "/users", bytes.NewReader(userJSON))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

	})

	Describe("POST /users/login", func() {

		BeforeEach(func() {
			fmt.Println("***saveUser()***")
			user = saveUser()
		})

		It("login successfully", func() {
			userLogin := models.UserLogin{Email: "michel.ingesoft@gmail.com", Password: "h1h1h1h1h1h1"}
			userJSON, _ := json.Marshal(userLogin)
			resp, _ := client.CallRequest("POST", "/users/login", bytes.NewReader(userJSON))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("fails login incorrect password", func() {
			userLogin := models.UserLogin{Email: "michel.ingesoft@gmail.com", Password: "h1"}
			userJSON, _ := json.Marshal(userLogin)
			resp, _ := client.CallRequest("POST", "/users/login", bytes.NewReader(userJSON))
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

	})

})
