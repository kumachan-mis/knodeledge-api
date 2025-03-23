/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// ChapterListResponse - Response Body for Chapter List API
type ChapterListResponse struct {
	Chapters []ChapterWithSections `json:"chapters"`
}
