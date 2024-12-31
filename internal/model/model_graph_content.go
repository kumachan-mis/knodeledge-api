/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// GraphContent - Graph object with only content fields
type GraphContent struct {

	// Auto-generated section ID
	Id string `json:"id"`

	// Graph paragraph
	Paragraph string `json:"paragraph"`

	Children []GraphChild `json:"children"`
}
