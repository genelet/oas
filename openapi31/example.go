// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import "encoding/json"

// Example represents an example of a media type.
// It can also represent a Reference (when isReference is true).
type Example struct {
	// Internal marker for reference
	isReference bool

	// Reference fields (used when isReference is true)
	Ref     string `json:"$ref,omitempty"`
	Summary string `json:"summary,omitempty"` // Shared with example summary

	// Example fields
	Description   string         `json:"description,omitempty"`
	Value         any            `json:"value,omitempty"`
	ExternalValue string         `json:"externalValue,omitempty"`
	Extensions    map[string]any `json:"-"`
}

var exampleKnownFields = []string{"$ref", "summary", "description", "value", "externalValue"}

// IsReference checks if this example is actually a reference ($ref)
func (e *Example) IsReference() bool {
	if e == nil {
		return false
	}
	return e.isReference
}

// NewExampleReference creates an example that is actually a reference
func NewExampleReference(ref string) *Example {
	return &Example{isReference: true, Ref: ref}
}

type exampleAlias Example

type exampleRefOnly struct {
	Ref     string `json:"$ref"`
	Summary string `json:"summary,omitempty"`
}

func (e *Example) UnmarshalJSON(data []byte) error {
	var alias exampleAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*e = Example(alias)
	if e.Ref != "" {
		e.isReference = true
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Extensions = extractExtensions(raw, exampleKnownFields)
	return nil
}

func (e Example) MarshalJSON() ([]byte, error) {
	if e.IsReference() {
		ref := exampleRefOnly{
			Ref:     e.Ref,
			Summary: e.Summary,
		}
		return marshalWithExtensions(&ref, e.Extensions)
	}
	alias := exampleAlias(e)
	return marshalWithExtensions(&alias, e.Extensions)
}
