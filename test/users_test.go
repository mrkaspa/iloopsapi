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

	Describe("POST /users", func() {

		It("create an user", func() {
			userLogin := defaultUser()
			userJSON, _ := json.Marshal(userLogin)
			client.CallRequest("POST", "/users", bytes.NewReader(userJSON)).WithResponse(func(resp *http.Response) error {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				return nil
			})
		})

	})

	Describe("POST /users/login", func() {

		BeforeEach(func() {
			fmt.Println("***saveUser()***")
			user = saveUser()
		})

		It("login successfully", func() {
			userLogin := defaultUser()
			userJSON, _ := json.Marshal(userLogin)
			client.CallRequest("POST", "/users/login", bytes.NewReader(userJSON)).WithResponse(func(resp *http.Response) error {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				return nil
			})
		})

		It("fails login incorrect password", func() {
			userLogin := models.UserLogin{Email: "michel.ingesoft@gmail.com", Password: "h1"}
			userJSON, _ := json.Marshal(userLogin)
			client.CallRequest("POST", "/users/login", bytes.NewReader(userJSON)).WithResponse(func(resp *http.Response) error {
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				return nil
			})
		})

	})

})
