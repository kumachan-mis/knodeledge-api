/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// ChapterWithSections - ChapterWithSections object
type ChapterWithSections struct {

	// Auto-generated chapter ID
	Id string `json:"id"`

	// Chapter name
	Name string `json:"name"`

	// Chapter number
	Number int32 `json:"number"`

	Sections []SectionOfChapter `json:"sections"`
}
