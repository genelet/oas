// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import "encoding/json"

// Parameter describes a single operation parameter.
// It can also represent a Reference (when isReference is true).
type Parameter struct {
	// Internal marker for reference
	isReference bool

	// Reference fields (used when isReference is true)
	Ref         string `json:"$ref,omitempty"`
	Summary     string `json:"summary,omitempty"`     // Reference summary
	Description string `json:"description,omitempty"` // Shared with parameter description

	// Parameter fields
	Name            string                `json:"name,omitempty"`
	In              string                `json:"in,omitempty"` // query, header, path, cookie
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
	"$ref", "summary", "description", "name", "in", "required", "deprecated",
	"allowEmptyValue", "style", "explode", "allowReserved", "schema",
	"content", "example", "examples",
}

// IsReference checks if this parameter is actually a reference ($ref)
func (p *Parameter) IsReference() bool {
	if p == nil {
		return false
	}
	return p.isReference
}

// NewParameterReference creates a parameter that is actually a reference
func NewParameterReference(ref string) *Parameter {
	return &Parameter{isReference: true, Ref: ref}
}

type parameterAlias Parameter

type parameterRefOnly struct {
	Ref         string `json:"$ref"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
}

func (p *Parameter) UnmarshalJSON(data []byte) error {
	var alias parameterAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*p = Parameter(alias)
	if p.Ref != "" {
		p.isReference = true
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Extensions = extractExtensions(raw, parameterKnownFields)
	return nil
}

func (p Parameter) MarshalJSON() ([]byte, error) {
	if p.IsReference() {
		ref := parameterRefOnly{
			Ref:         p.Ref,
			Summary:     p.Summary,
			Description: p.Description,
		}
		return marshalWithExtensions(&ref, p.Extensions)
	}
	alias := parameterAlias(p)
	return marshalWithExtensions(&alias, p.Extensions)
}

// Header represents a header parameter (similar to Parameter but without name and in).
// It can also represent a Reference.
type Header struct {
	// Internal marker for reference
	isReference bool

	// Reference fields
	Ref         string `json:"$ref,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`

	// Header fields
	Required   bool                  `json:"required,omitempty"`
	Deprecated bool                  `json:"deprecated,omitempty"`
	Style      string                `json:"style,omitempty"`
	Explode    *bool                 `json:"explode,omitempty"`
	Schema     *Schema               `json:"schema,omitempty"`
	Content    map[string]*MediaType `json:"content,omitempty"`
	Example    any                   `json:"example,omitempty"`
	Examples   map[string]*Example   `json:"examples,omitempty"`
	Extensions map[string]any        `json:"-"`
}

var headerKnownFields = []string{
	"$ref", "summary", "description", "required", "deprecated", "style", "explode",
	"schema", "content", "example", "examples",
}

// IsReference checks if this header is actually a reference ($ref)
func (h *Header) IsReference() bool {
	if h == nil {
		return false
	}
	return h.isReference
}

// NewHeaderReference creates a header that is actually a reference
func NewHeaderReference(ref string) *Header {
	return &Header{isReference: true, Ref: ref}
}

type headerAlias Header

type headerRefOnly struct {
	Ref         string `json:"$ref"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
}

func (h *Header) UnmarshalJSON(data []byte) error {
	var alias headerAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*h = Header(alias)
	if h.Ref != "" {
		h.isReference = true
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	h.Extensions = extractExtensions(raw, headerKnownFields)
	return nil
}

func (h Header) MarshalJSON() ([]byte, error) {
	if h.IsReference() {
		ref := headerRefOnly{
			Ref:         h.Ref,
			Summary:     h.Summary,
			Description: h.Description,
		}
		return marshalWithExtensions(&ref, h.Extensions)
	}
	alias := headerAlias(h)
	return marshalWithExtensions(&alias, h.Extensions)
}
