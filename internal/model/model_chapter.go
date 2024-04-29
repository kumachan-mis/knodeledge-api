/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// Chapter - Chapter object
type Chapter struct {

	// Auto-generated chapter ID
	Id string `json:"id"`

	// Chapter name
	Name string `json:"name"`

	Sections []Section `json:"sections"`

	// next chapter ID
	NextId string `json:"nextId"`
}
