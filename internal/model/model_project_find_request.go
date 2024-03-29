/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// ProjectFindRequest - Request Body for Project Find API
type ProjectFindRequest struct {
	User UserOnlyId `json:"user"`

	Project ProjectOnlyId `json:"project"`
}
