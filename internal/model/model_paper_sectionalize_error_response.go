/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// PaperSectionalizeErrorResponse - Error Response Body for Paper Sectionalize API
type PaperSectionalizeErrorResponse struct {
	User UserOnlyIdError `json:"user,omitempty"`

	Project ProjectOnlyIdError `json:"project,omitempty"`

	Paper PaperOnlyIdError `json:"paper,omitempty"`

	Sections []SectionWithoutAutofieldError `json:"sections,omitempty"`
}
