/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// Graph - Graph object
type Graph struct {

	// Auto-generated section ID
	Id string `json:"id"`

	// Graph name
	Name string `json:"name"`

	// Graph paragraph
	Paragraph string `json:"paragraph"`

	Children []GraphChild `json:"children"`
}
