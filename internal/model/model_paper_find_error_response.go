/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// PaperFindErrorResponse - Error Response Body for Paper Find API
type PaperFindErrorResponse struct {

	// Error message when request body format is invalid
	Message string `json:"message,omitempty"`

	User UserOnlyIdError `json:"user,omitempty"`

	Project ProjectOnlyIdError `json:"project,omitempty"`

	Chapter ChapterOnlyIdError `json:"chapter,omitempty"`
}
