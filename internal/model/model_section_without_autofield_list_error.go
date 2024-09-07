/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

// SectionWithoutAutofieldListError - Error Message for SectionWithoutAutofield list
type SectionWithoutAutofieldListError struct {

	// Error message for overall of sections
	Message string `json:"message,omitempty"`

	Items []SectionWithoutAutofieldError `json:"items,omitempty"`
}