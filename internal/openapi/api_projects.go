/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"github.com/gin-gonic/gin"
)

type ProjectsAPI interface {

	// ProjectsCreate Post /api/projects/create
	// Create new project
	ProjectsCreate(c *gin.Context)

	// ProjectsDelete Post /api/projects/delete
	// Delete project
	ProjectsDelete(c *gin.Context)

	// ProjectsFind Get /api/projects/find
	// Find project
	ProjectsFind(c *gin.Context)

	// ProjectsList Get /api/projects/list
	// Get list of projects
	ProjectsList(c *gin.Context)

	// ProjectsUpdate Post /api/projects/update
	// Update project
	ProjectsUpdate(c *gin.Context)
}
