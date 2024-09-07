/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// PaperUpdateErrorResponse - Error Response Body for Paper Update API
type PaperUpdateErrorResponse struct {

	// Error message when request body format is invalid
	Message string `json:"message,omitempty"`

	User UserOnlyIdError `json:"user,omitempty"`

	Project ProjectOnlyIdError `json:"project,omitempty"`

	Paper PaperError `json:"paper,omitempty"`
}
