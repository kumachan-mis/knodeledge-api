/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// GraphSectionalizeErrorResponse - Error Response Body for Graph Sectionalize API
type GraphSectionalizeErrorResponse struct {
	User UserOnlyIdError `json:"user,omitempty"`

	Project ProjectOnlyIdError `json:"project,omitempty"`

	Chapter ChapterOnlyIdError `json:"chapter,omitempty"`

	Sections SectionWithoutAutofieldListError `json:"sections,omitempty"`
}
