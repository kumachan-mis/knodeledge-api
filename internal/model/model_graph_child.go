/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// GraphChild - GraphChild object
type GraphChild struct {

	// Child node name of the graph
	Name string `json:"name"`

	// Graph relation
	Relation string `json:"relation"`

	// Graph description
	Description string `json:"description"`

	Children []GraphChild `json:"children"`
}