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

var _ = Describe("Execution", func() {

	var (
		user    models.User
		project models.Project
	)

	BeforeEach(func() {
		cleanDB()
		user = saveUser()
		project = addProject(user)
	})

	Describe("POST /execution/:project_id", func() {

		It("create an execution", func() {
			execution := models.Execution{}
			executionJSON, _ := json.Marshal(execution)
			resp, _ := client.CallRequest("POST", fmt.Sprintf("/executions/%d", project.ID), bytes.NewReader(executionJSON))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

	})

})
