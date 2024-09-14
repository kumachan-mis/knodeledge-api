/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// GraphUpdateRequest - Request Body for Graph Update API
type GraphUpdateRequest struct {
	User UserOnlyId `json:"user"`

	Project ProjectOnlyId `json:"project"`

	Chapter ChapterOnlyId `json:"chapter"`

	Graph GraphContent `json:"graph"`
}
