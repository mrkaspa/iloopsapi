package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"bitbucket.org/kiloops/api/gitadmin"
	"bitbucket.org/kiloops/api/models"
	"github.com/mrkaspa/go-helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SSH", func() {

	BeforeEach(func() {
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

		BeforeEach(func() {
			ssh = addSSH(user)
		})

		AfterEach(func() {
			gitadmin.DeleteSSH(user.Email, ssh.ID)
		})

		It("create an ssh", func() {
			resp, _ := client.CallRequestNoBodytWithHeaders("DELETE", fmt.Sprintf("/ssh/%d", ssh.ID), authHeaders(user))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("throws error when delete an unknown ssh", func() {
			resp, _ := client.CallRequestNoBodytWithHeaders("DELETE", "/ssh/-1", authHeaders(user))
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

	})

	Describe("After a project is created", func() {

		BeforeEach(func() {
			ssh = addSSH(user)
			project = addProject(user)
			anotherSSH = addAnotherSSH(user)
		})

		AfterEach(func() {
			gitadmin.DeleteSSH(user.Email, ssh.ID)
		})

		It("should containg the two keys in the default project file", func() {
			projectPath := gitadmin.ProjectPath(project.Slug)
			data, err := ioutil.ReadFile(projectPath)
			if err != nil {
				panic("error reading file " + projectPath)
			}
			eq := "@users_" + project.Slug + " = " + user.Email + "-" + strconv.Itoa(ssh.ID) + " " + user.Email + "-" + strconv.Itoa(anotherSSH.ID) + "\nrepo " + project.Slug + "\n  RW+ = @users_" + project.Slug
			Expect(string(data)).To(BeEquivalentTo(eq))
		})

	})

})
