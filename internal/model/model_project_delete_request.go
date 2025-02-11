/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// ProjectDeleteRequest - Request Body for Project Delete API
type ProjectDeleteRequest struct {
	User UserOnlyId `json:"user"`

	Project ProjectOnlyId `json:"project"`
}
