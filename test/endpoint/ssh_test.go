package endpoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/kiloops/api/models"
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
			sshKey := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDCadJM0DdJotRnSWW7coFcCxMW1cCZIJqkyW3wMQoUOU2VHuLExh44tDpSAiz2EEeFlqJ5hI67ZI+3bSx7puKSr44l78H/Kb8UDLidAUao7JZoo0thq7bAVesGr+8aligmULvxH3sQqstI9yNcifJ56jHUVTB14PslBmhA56pmGOva0ojmdt9l2aBy4LxQBDc5Js+AcPlfC2zXE7rtaiafB/M3992V+7CEisbAv7CpsI3SPdpW2p4mfR1zMVpf4Jt6lQJW6Sr53/bzAP4/Tif3fgbZhoSL8qnnLi3556gWi90FwFhCoqqDR/lN3sxJQx5NxxCF8mbNgpmS5qDptFyF michel.ingesoft@gmail.com`
			ssh := models.SSH{PublicKey: sshKey}
			sshJSON, _ := json.Marshal(ssh)
			resp, _ := client.CallRequestWithHeaders("POST", "/ssh", bytes.NewReader(sshJSON), authHeaders(user))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("crashes when creates an ssh empty", func() {
			ssh := models.SSH{PublicKey: ""}
			sshJSON, _ := json.Marshal(ssh)
			resp, _ := client.CallRequestWithHeaders("POST", "/ssh", bytes.NewReader(sshJSON), authHeaders(user))
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
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
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

	})

})
