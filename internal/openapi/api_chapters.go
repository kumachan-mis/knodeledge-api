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

type ChaptersAPI interface {

	// ChaptersCreate Post /api/chapters/create
	// Create new Chapter
	ChaptersCreate(c *gin.Context)

	// ChaptersDelete Post /api/chapters/delete
	// Delete chapter
	ChaptersDelete(c *gin.Context)

	// ChaptersList Get /api/chapters/list
	// Get list of chapters for a project
	ChaptersList(c *gin.Context)

	// ChaptersUpdate Post /api/chapters/update
	// Update chapter
	ChaptersUpdate(c *gin.Context)
}
