/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// ChapterCreateErrorResponse - Error Response Body for Chapter Create API
type ChapterCreateErrorResponse struct {

	// Error message when request body format is invalid
	Message string `json:"message"`

	User UserOnlyIdError `json:"user,omitempty"`

	Project ProjectOnlyIdError `json:"project,omitempty"`

	Chapter ChapterWithoutAutofieldError `json:"chapter,omitempty"`
}
