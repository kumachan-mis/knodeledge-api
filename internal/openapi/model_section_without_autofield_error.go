/*
 * Web API of kNODEledge
 *
 * App to Create Graphically-Summarized Notes in Three Steps
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// SectionWithoutAutofieldError - Error Message for SectionWithoutAutofield object
type SectionWithoutAutofieldError struct {

	// Error message for section name
	Name string `json:"name,omitempty"`

	// Error message for section content
	Content string `json:"content,omitempty"`
}
