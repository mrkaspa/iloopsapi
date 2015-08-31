package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/kiloops/api/gitadmin"
	"bitbucket.org/kiloops/api/models"
	"github.com/jinzhu/gorm"
	"github.com/mrkaspa/go-helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Projects", func() {

	BeforeEach(func() {
		user = saveUser()
		ssh = addSSH(user)
	})

	AfterEach(func() {
		gitadmin.DeleteSSH(user.Email, ssh.ID)
	})

	Describe("POST /projects", func() {

		It("create a project", func() {
			project := models.Project{Name: "Demo Project"}
			projectJSON, _ := json.Marshal(project)
			resp, _ := client.CallRequestWithHeaders("POST", "/projects", bytes.NewReader(projectJSON), authHeaders(user))
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			defer resp.Body.Close()
			getBodyJSON(resp, &project)
			Expect(project.Slug).ToNot(BeNil())
			Expect(project.Slug).ToNot(BeEmpty())
			projectsOwned := user.OwnedProjects()
			Expect(len(projectsOwned)).To(Equal(1))
			path := gitadmin.ProjectPath(project.Slug)
			Expect(helpers.FileExists(path)).To(BeTrue())
			gitadmin.DeleteProject(project.Slug)
		})

	})

	Context("After adding a project", func() {

		BeforeEach(func() {
			project = addProject(user)
		})

		AfterEach(func() {
			gitadmin.DeleteProject(project.Slug)
		})

		Describe("GET /projects", func() {

			It("lists all the projects", func() {
				var projects []models.Project
				resp, _ := client.CallRequestNoBodytWithHeaders("GET", "/projects", authHeaders(user))
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				defer resp.Body.Close()
				getBodyJSON(resp, &projects)
				Expect(len(projects)).To(Equal(1))
			})

		})

		Describe("GET /projects/:slug", func() {

			It("gets a project", func() {
				var projectResp models.Project
				resp, _ := client.CallRequestNoBodytWithHeaders("GET", fmt.Sprintf("/projects/%s", project.Slug), authHeaders(user))
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				defer resp.Body.Close()
				getBodyJSON(resp, &projectResp)
				Expect(projectResp.Name).To(Equal(project.Name))
			})

		})

		Describe("DELETE /projects/:slug", func() {

			It("deletes a project", func() {
				resp, _ := client.CallRequestNoBodytWithHeaders("DELETE", fmt.Sprintf("/projects/%s", project.Slug), authHeaders(user))
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				projectsOwned := user.OwnedProjects()
				Expect(len(projectsOwned)).To(Equal(0))
			})

		})

		Describe("PUT /projects/:slug/leave", func() {

			It("an admin tries to leave a project", func() {
				resp, _ := client.CallRequestNoBodytWithHeaders("PUT", fmt.Sprintf("/projects/%s/leave", project.Slug), authHeaders(user))
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

		})

		Context("Adding another user to the project", func() {

			BeforeEach(func() {
				otherUser = saveOtherUser()
				anotherSSH = addAnotherSSH(otherUser)
			})

			AfterEach(func() {
				gitadmin.DeleteSSH(otherUser.Email, anotherSSH.ID)
			})

			Describe("PUT /projects/:slug/add/:email", func() {

				It("adds another user to the project", func() {
					resp, _ := client.CallRequestNoBodytWithHeaders("PUT", fmt.Sprintf("/projects/%s/add/%s", project.Slug, otherUser.Email), authHeaders(user))
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					projectsCollab := otherUser.CollaboratorProjects()
					Expect(len(projectsCollab)).To(Equal(1))
				})

			})

			Describe("PUT /projects/:slug/delegate/:email", func() {

				It("delegates admin role to another user", func() {
					models.InTx(func(txn *gorm.DB) bool {
						project.AddUser(txn, &otherUser)
						return true
					})
					resp, _ := client.CallRequestNoBodytWithHeaders("PUT", fmt.Sprintf("/projects/%s/delegate/%s", project.Slug, otherUser.Email), authHeaders(user))
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					projectsCollab := user.CollaboratorProjects()
					Expect(len(projectsCollab)).To(Equal(1))
					projectsOwned := otherUser.OwnedProjects()
					Expect(len(projectsOwned)).To(Equal(1))
				})

			})

			Describe("PUT /projects/:slug/leave", func() {

				It("an user leaves a project", func() {
					models.InTx(func(txn *gorm.DB) bool {
						project.AddUser(txn, &otherUser)
						return true
					})
					resp, _ := client.CallRequestNoBodytWithHeaders("PUT", fmt.Sprintf("/projects/%s/leave", project.Slug), authHeaders(otherUser))
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})

			})

			Describe("PUT /projects/:slug/remove/:email", func() {

				It("remove an user from a project", func() {
					models.InTx(func(txn *gorm.DB) bool {
						project.AddUser(txn, &otherUser)
						return true
					})
					resp, _ := client.CallRequestNoBodytWithHeaders("DELETE", fmt.Sprintf("/projects/%s/remove/%s", project.Slug, otherUser.Email), authHeaders(user))
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					projectsCollab := user.CollaboratorProjects()
					Expect(len(projectsCollab)).To(Equal(0))
				})

			})

		})

	})

	Context("Scheduling a project", func() {

		Describe("POST /projects/:slug/schedule", func() {

			BeforeEach(func() {
				project = addProject(user)
			})

			AfterEach(func() {
				gitadmin.DeleteProject(project.Slug)
			})

			It("Should get ok", func() {
				projectConfig := models.ProjectConfig{
					Name:  project.Name,
					AppID: project.Slug,
					Loops: models.Loops{CronFormat: "@every 1m"},
				}
				projectConfigJSON, _ := json.Marshal(projectConfig)
				resp, _ := client.CallRequest("POST", fmt.Sprintf("/projects/%s/schedule", project.Slug), bytes.NewReader(projectConfigJSON))
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

		})

	})

})
