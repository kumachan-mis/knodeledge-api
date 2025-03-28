/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// GraphFindRequest - Request Parameters for Graph Find API
type GraphFindRequest struct {

	// User ID
	UserId string `json:"userId" form:"userId"`

	// Auto-generated project ID
	ProjectId string `json:"projectId" form:"projectId"`

	// Auto-generated chapter ID
	ChapterId string `json:"chapterId" form:"chapterId"`

	// Auto-generated section ID
	SectionId string `json:"sectionId" form:"sectionId"`
}
