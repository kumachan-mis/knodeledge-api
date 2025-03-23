/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// ChapterDeleteErrorResponse - Error Response Body for Chapter Delete API
type ChapterDeleteErrorResponse struct {

	// Error message when request body format is invalid
	Message string `json:"message"`

	User UserOnlyIdError `json:"user,omitempty"`

	Project ProjectOnlyIdError `json:"project,omitempty"`

	Chapter ChapterOnlyIdError `json:"chapter,omitempty"`
}
