package test

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mrkaspa/iloopsapi/ierrors"
	"github.com/mrkaspa/iloopsapi/models"
	"github.com/mrkaspa/iloopsapi/utils"
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
			var appError ierrors.AppError
			client.CallRequest("POST", "/users/login", bytes.NewReader(userJSON)).Solve(utils.MapExec{
				http.StatusConflict: utils.InfoExec{
					&appError,
					func(resp *http.Response) error {
						Expect(appError.Code).To(Equal(ierrors.ErrCodeCredential))
						return appError
					},
				},
			})
		})

	})

	Describe("POST /users/forgot", func() {

		BeforeEach(func() {
			user = saveUser()
		})

		It("requests a password change", func() {
			email := models.Email{Value: user.Email}
			emailJSON, _ := json.Marshal(email)
			client.CallRequest("POST", "/users/forgot", bytes.NewReader(emailJSON)).WithResponse(func(resp *http.Response) error {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				return nil
			})
		})

	})

	Describe("POST /users/change_password", func() {

		var token string

		BeforeEach(func() {
			user = saveUser()
			var passwordRequest models.PasswordRequest
			email := models.Email{Value: user.Email}
			emailJSON, _ := json.Marshal(email)
			client.CallRequest("POST", "/users/forgot", bytes.NewReader(emailJSON)).Solve(utils.MapExec{
				http.StatusOK: utils.InfoExec{
					&passwordRequest,
					func(resp *http.Response) error {
						token = passwordRequest.Token
						return nil
					},
				},
			})
		})

		It("requests a password change", func() {
			newPassword := "jokalive123"
			changePassword := models.ChangePassword{
				Token:    token,
				Password: newPassword,
			}
			changePasswordJSON, _ := json.Marshal(changePassword)
			client.CallRequest("POST", "/users/change_password", bytes.NewReader(changePasswordJSON)).WithResponse(func(resp *http.Response) error {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				return nil
			})
		})

	})

})
