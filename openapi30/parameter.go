// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi30

import "encoding/json"

// Parameter describes a single operation parameter.
// It can also represent a Reference (when Ref is set).
type Parameter struct {
	// Reference field
	Ref string `json:"$ref,omitempty"`

	// Parameter fields
	Name            string                `json:"name,omitempty"`
	In              string                `json:"in,omitempty"` // query, header, path, cookie
	Description     string                `json:"description,omitempty"`
	Required        bool                  `json:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty"`
	Style           string                `json:"style,omitempty"`
	Explode         *bool                 `json:"explode,omitempty"`
	AllowReserved   bool                  `json:"allowReserved,omitempty"`
	Schema          *Schema               `json:"schema,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty"`
	Example         any                   `json:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty"`
	Extensions      map[string]any        `json:"-"`
}

var parameterKnownFields = []string{
	"$ref", "name", "in", "description", "required", "deprecated",
	"allowEmptyValue", "style", "explode", "allowReserved", "schema",
	"content", "example", "examples",
}

// IsReference checks if this parameter is actually a reference ($ref)
func (p *Parameter) IsReference() bool {
	return p != nil && p.Ref != ""
}

// NewParameterReference creates a parameter that is actually a reference
func NewParameterReference(ref string) *Parameter {
	return &Parameter{Ref: ref}
}

type parameterAlias Parameter

func (p *Parameter) UnmarshalJSON(data []byte) error {
	var alias parameterAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*p = Parameter(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Extensions = extractExtensions(raw, parameterKnownFields)
	return nil
}

func (p Parameter) MarshalJSON() ([]byte, error) {
	alias := parameterAlias(p)
	return marshalWithExtensions(&alias, p.Extensions)
}

// Header represents a header parameter (similar to Parameter but without name and in).
// It can also represent a Reference (when Ref is set).
type Header struct {
	// Reference field
	Ref string `json:"$ref,omitempty"`

	// Header fields
	Description     string                `json:"description,omitempty"`
	Required        bool                  `json:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty"`
	Style           string                `json:"style,omitempty"`
	Explode         *bool                 `json:"explode,omitempty"`
	AllowReserved   bool                  `json:"allowReserved,omitempty"`
	Schema          *Schema               `json:"schema,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty"`
	Example         any                   `json:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty"`
	Extensions      map[string]any        `json:"-"`
}

var headerKnownFields = []string{
	"$ref", "description", "required", "deprecated", "allowEmptyValue",
	"style", "explode", "allowReserved", "schema", "content", "example", "examples",
}

// IsReference checks if this header is actually a reference ($ref)
func (h *Header) IsReference() bool {
	return h != nil && h.Ref != ""
}

// NewHeaderReference creates a header that is actually a reference
func NewHeaderReference(ref string) *Header {
	return &Header{Ref: ref}
}

type headerAlias Header

func (h *Header) UnmarshalJSON(data []byte) error {
	var alias headerAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*h = Header(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	h.Extensions = extractExtensions(raw, headerKnownFields)
	return nil
}

func (h Header) MarshalJSON() ([]byte, error) {
	alias := headerAlias(h)
	return marshalWithExtensions(&alias, h.Extensions)
}
