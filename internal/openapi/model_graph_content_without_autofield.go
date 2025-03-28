/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// GraphContentWithoutAutofield - Graph object with only content fields without auto-generated fields
type GraphContentWithoutAutofield struct {

	// Graph paragraph
	Paragraph string `json:"paragraph"`

	Children []GraphChild `json:"children"`
}
