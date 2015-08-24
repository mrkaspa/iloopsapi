package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/kiloops/api/gitadmin"
	"bitbucket.org/kiloops/api/models"
	"github.com/mrkaspa/go-helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SSH", func() {

	var user models.User

	BeforeEach(func() {
		cleanDB()
		user = saveUser()
	})

	Describe("POST /ssh", func() {

		It("create an ssh", func() {
			sshKey := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDlCc96zWY05/vFIcP5NLhi8bIVkcUdSyet1Dw7+rQqbeEJaQ0Ifz/x17AGkQAnC0ZjdHI7sCFjVGuvk6agw6MJzKU8a+iWisAVu4hvv22DXBPKYak28GMEW3e0Ba/8mUiCdLCW5lfQ85QmDABqdWb6BGy2VSJ/k4NfWW728RwbQf1MZSwS2+kqvR3XjpkpMETLz5DmRty6Dqp3al73JbE7raWhidoYeS0wiJKsWiaucfewz+feubNkEnO5/p1v1zpAlaPYEVvZEeG5ABchNZ4Co+SGvVd4+FuxVgLkPOqpV5y3JFFrmSJE4HMsin96u/3OHcgVwew6nyE9YyoKZ/rL michel.ing@hotmail.com`
			ssh := models.SSH{Name: "demo", PublicKey: sshKey}
			sshJSON, _ := json.Marshal(ssh)
			resp, _ := client.CallRequestWithHeaders("POST", "/ssh", bytes.NewReader(sshJSON), authHeaders(user))
			var sshResp models.SSH
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			getBodyJSON(resp, &sshResp)
			path := gitadmin.KeyPath(user.Email, sshResp.ID)
			Expect(helpers.FileExists(path)).To(BeTrue())
			err := gitadmin.DeleteSSH(user.Email, sshResp.ID)
			Expect(err).To(BeNil())
		})

		It("crashes when creates an ssh empty", func() {
			ssh := models.SSH{PublicKey: ""}
			sshJSON, _ := json.Marshal(ssh)
			resp, _ := client.CallRequestWithHeaders("POST", "/ssh", bytes.NewReader(sshJSON), authHeaders(user))
			Expect(resp.StatusCode).To(Equal(http.StatusConflict))
		})

	})

	Describe("DELETE /ssh/:id", func() {

		var ssh models.SSH

		BeforeEach(func() {
			ssh = addSSH(user)
		})

		It("create an ssh", func() {
			resp, _ := client.CallRequestWithHeaders("DELETE", fmt.Sprintf("/ssh/%d", ssh.ID), bytes.NewReader(emptyJSON), authHeaders(user))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("throws error when delete an unknown ssh", func() {
			resp, _ := client.CallRequestWithHeaders("DELETE", "/ssh/-1", bytes.NewReader(emptyJSON), authHeaders(user))
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

	})

})
