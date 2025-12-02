// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import "encoding/json"

// RequestBody describes a single request body.
// It can also represent a Reference (when isReference is true).
type RequestBody struct {
	// Internal marker for reference
	isReference bool

	// Reference fields (used when isReference is true)
	Ref         string `json:"$ref,omitempty"`
	Summary     string `json:"summary,omitempty"`     // Reference summary
	Description string `json:"description,omitempty"` // Shared with request body description

	// RequestBody fields
	Content    map[string]*MediaType `json:"content,omitempty"`
	Required   bool                  `json:"required,omitempty"`
	Extensions map[string]any        `json:"-"`
}

var requestBodyKnownFields = []string{"$ref", "summary", "description", "content", "required"}

// IsReference checks if this request body is actually a reference ($ref)
func (rb *RequestBody) IsReference() bool {
	if rb == nil {
		return false
	}
	return rb.isReference
}

// NewRequestBodyReference creates a request body that is actually a reference
func NewRequestBodyReference(ref string) *RequestBody {
	return &RequestBody{isReference: true, Ref: ref}
}

type requestBodyAlias RequestBody

type requestBodyRefOnly struct {
	Ref         string `json:"$ref"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
}

func (rb *RequestBody) UnmarshalJSON(data []byte) error {
	var alias requestBodyAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*rb = RequestBody(alias)
	if rb.Ref != "" {
		rb.isReference = true
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	rb.Extensions = extractExtensions(raw, requestBodyKnownFields)
	return nil
}

func (rb RequestBody) MarshalJSON() ([]byte, error) {
	if rb.IsReference() {
		ref := requestBodyRefOnly{
			Ref:         rb.Ref,
			Summary:     rb.Summary,
			Description: rb.Description,
		}
		return marshalWithExtensions(&ref, rb.Extensions)
	}
	alias := requestBodyAlias(rb)
	return marshalWithExtensions(&alias, rb.Extensions)
}
