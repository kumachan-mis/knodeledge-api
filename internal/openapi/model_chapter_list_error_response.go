/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// ChapterListErrorResponse - Error Response Body for Chapter List API
type ChapterListErrorResponse struct {

	// Error message when request body format is invalid
	Message string `json:"message"`

	// Error message for user ID
	UserId string `json:"userId,omitempty"`

	// Error message for project ID
	ProjectId string `json:"projectId,omitempty"`
}
