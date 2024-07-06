/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// PaperError - Error Message for Paper object
type PaperError struct {

	// Error message for paper ID
	Id string `json:"id,omitempty"`

	// Error message for paper content
	Content string `json:"content,omitempty"`
}
