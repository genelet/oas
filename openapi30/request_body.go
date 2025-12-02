// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi30

import "encoding/json"

// RequestBody describes a single request body.
// It can also represent a Reference (when Ref is set).
type RequestBody struct {
	// Reference field
	Ref string `json:"$ref,omitempty"`

	// RequestBody fields
	Description string                `json:"description,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
	Required    bool                  `json:"required,omitempty"`
	Extensions  map[string]any        `json:"-"`
}

var requestBodyKnownFields = []string{"$ref", "description", "content", "required"}

// IsReference checks if this request body is actually a reference ($ref)
func (rb *RequestBody) IsReference() bool {
	return rb != nil && rb.Ref != ""
}

// NewRequestBodyReference creates a request body that is actually a reference
func NewRequestBodyReference(ref string) *RequestBody {
	return &RequestBody{Ref: ref}
}

type requestBodyAlias RequestBody

func (rb *RequestBody) UnmarshalJSON(data []byte) error {
	var alias requestBodyAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*rb = RequestBody(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	rb.Extensions = extractExtensions(raw, requestBodyKnownFields)
	return nil
}

func (rb RequestBody) MarshalJSON() ([]byte, error) {
	alias := requestBodyAlias(rb)
	return marshalWithExtensions(&alias, rb.Extensions)
}
