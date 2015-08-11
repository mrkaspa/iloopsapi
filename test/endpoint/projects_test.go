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

var _ = Describe("Projects", func() {

	var user models.User

	BeforeEach(func() {
		cleanDB()
		user = saveUser()
	})

	Describe("POST /projects", func() {

		It("create a project", func() {
			project := models.Project{Name: "Demo Project"}
			projectJSON, _ := json.Marshal(project)
			resp, _ := client.CallRequestWithHeaders("POST", "/projects", bytes.NewReader(projectJSON), authHeaders(user))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			getBodyJSON(resp, &project)
			Expect(project.Slug).ToNot(BeNil())
			projectsOwned := user.OwnedProjects()
			Expect(len(projectsOwned)).To(Equal(1))
		})

	})

	Context("After adding a project", func() {

		var project models.Project

		BeforeEach(func() {
			project = addProject(user)
		})

		Describe("GET /projects", func() {

			It("lists all the projects", func() {
				var projects []models.Project
				resp, _ := client.CallRequestWithHeaders("GET", "/projects", bytes.NewReader(emptyJSON), authHeaders(user))
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				getBodyJSON(resp, &projects)
				Expect(len(projects)).To(Equal(1))
			})

		})

		Describe("GET /projects/:id", func() {

			It("gets a project", func() {
				var projectResp models.Project
				resp, _ := client.CallRequestWithHeaders("GET", fmt.Sprintf("/projects/%d", project.ID), bytes.NewReader(emptyJSON), authHeaders(user))
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				getBodyJSON(resp, &projectResp)
				Expect(projectResp.Name).To(Equal(project.Name))
			})

		})

		Describe("DELETE /projects/:id", func() {

			It("deletes a project", func() {
				resp, _ := client.CallRequestWithHeaders("DELETE", fmt.Sprintf("/projects/%d", project.ID), bytes.NewReader(emptyJSON), authHeaders(user))
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				projectsOwned := user.OwnedProjects()
				Expect(len(projectsOwned)).To(Equal(0))
			})

		})

		Describe("PUT /projects/:id/leave", func() {

			FIt("deletes a project", func() {
				resp, _ := client.CallRequestWithHeaders("PUT", fmt.Sprintf("/projects/%d/leave", project.ID), bytes.NewReader(emptyJSON), authHeaders(user))
				Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
			})

		})

		Context("Adding another user to the project", func() {

		})

	})

})
