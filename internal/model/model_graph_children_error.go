/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// GraphChildrenError - Error Message for GraphChild list
type GraphChildrenError struct {

	// Error message for overall of children
	Message string `json:"message,omitempty"`

	Items []GraphChildError `json:"items,omitempty"`
}
