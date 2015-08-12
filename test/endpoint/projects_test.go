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

		FIt("create a project", func() {
			project := models.Project{Name: "Demo Project"}
			projectJSON, _ := json.Marshal(project)
			resp, _ := client.CallRequestWithHeaders("POST", "/projects", bytes.NewReader(projectJSON), authHeaders(user))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			getBodyJSON(resp, &project)
			Expect(project.Slug).ToNot(BeNil())
			Expect(project.Slug).ToNot(BeEmpty())
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

			It("an admin tries to leave a project", func() {
				resp, _ := client.CallRequestWithHeaders("PUT", fmt.Sprintf("/projects/%d/leave", project.ID), bytes.NewReader(emptyJSON), authHeaders(user))
				Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
			})

		})

		Context("Adding another user to the project", func() {

			var otherUser models.User

			BeforeEach(func() {
				otherUser = saveOtherUser()
			})

			Describe("PUT /projects/:id/add/:user_id", func() {

				It("adds another user to the project", func() {
					resp, _ := client.CallRequestWithHeaders("PUT", fmt.Sprintf("/projects/%d/add/%d", project.ID, otherUser.ID), bytes.NewReader(emptyJSON), authHeaders(user))
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					projectsCollab := otherUser.CollaboratorProjects()
					Expect(len(projectsCollab)).To(Equal(1))
				})

			})

			Describe("PUT /projects/:id/delegate/:user_id", func() {

				It("delegates admin role to another user", func() {
					models.Gdb.InTx(func(txn *models.KDB) {
						project.AddUser(txn, &otherUser)
					})
					resp, _ := client.CallRequestWithHeaders("PUT", fmt.Sprintf("/projects/%d/delegate/%d", project.ID, otherUser.ID), bytes.NewReader(emptyJSON), authHeaders(user))
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					projectsCollab := user.CollaboratorProjects()
					Expect(len(projectsCollab)).To(Equal(1))
					projectsOwned := otherUser.OwnedProjects()
					Expect(len(projectsOwned)).To(Equal(1))
				})

			})

			Describe("PUT /projects/:id/leave", func() {

				It("an user leaves a project", func() {
					models.Gdb.InTx(func(txn *models.KDB) {
						project.AddUser(txn, &otherUser)
					})
					resp, _ := client.CallRequestWithHeaders("PUT", fmt.Sprintf("/projects/%d/leave", project.ID), bytes.NewReader(emptyJSON), authHeaders(otherUser))
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})

			})

		})

	})

})
