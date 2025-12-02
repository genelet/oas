// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi30

import "encoding/json"

// Example represents an example object.
// It can also represent a Reference (when Ref is set).
type Example struct {
	// Reference field
	Ref string `json:"$ref,omitempty"`

	// Example fields
	Summary       string         `json:"summary,omitempty"`
	Description   string         `json:"description,omitempty"`
	Value         any            `json:"value,omitempty"`
	ExternalValue string         `json:"externalValue,omitempty"`
	Extensions    map[string]any `json:"-"`
}

var exampleKnownFields = []string{"$ref", "summary", "description", "value", "externalValue"}

// IsReference checks if this example is actually a reference ($ref)
func (e *Example) IsReference() bool {
	return e != nil && e.Ref != ""
}

// NewExampleReference creates an example that is actually a reference
func NewExampleReference(ref string) *Example {
	return &Example{Ref: ref}
}

type exampleAlias Example

func (e *Example) UnmarshalJSON(data []byte) error {
	var alias exampleAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*e = Example(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Extensions = extractExtensions(raw, exampleKnownFields)
	return nil
}

func (e Example) MarshalJSON() ([]byte, error) {
	alias := exampleAlias(e)
	return marshalWithExtensions(&alias, e.Extensions)
}
