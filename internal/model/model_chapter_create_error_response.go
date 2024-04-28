/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// ChapterCreateErrorResponse - Error Response Body for Chapter Create API
type ChapterCreateErrorResponse struct {
	User UserOnlyIdError `json:"user"`

	Project ProjectOnlyIdError `json:"project"`

	Chapter ChapterWithoutAutofieldError `json:"chapter"`
}