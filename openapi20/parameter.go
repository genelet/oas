// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi20

import "encoding/json"

// Parameter describes a single operation parameter.
// In Swagger 2.0, parameters can be:
// - Body parameters (in=body, with schema)
// - Non-body parameters (in=query/header/path/formData, with type, format, items, etc.)
type Parameter struct {
	// Reference field
	Ref string `json:"$ref,omitempty"`

	// Common fields
	Name        string `json:"name,omitempty"`
	In          string `json:"in,omitempty"` // query, header, path, formData, body
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`

	// Body parameter fields
	Schema *Schema `json:"schema,omitempty"`

	// Non-body parameter fields
	Type             string          `json:"type,omitempty"` // string, number, integer, boolean, array, file
	Format           string          `json:"format,omitempty"`
	AllowEmptyValue  bool            `json:"allowEmptyValue,omitempty"`
	Items            *Items          `json:"items,omitempty"`
	CollectionFormat string          `json:"collectionFormat,omitempty"` // csv, ssv, tsv, pipes, multi
	Default          any             `json:"default,omitempty"`
	Maximum          *float64        `json:"maximum,omitempty"`
	ExclusiveMaximum bool            `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64        `json:"minimum,omitempty"`
	ExclusiveMinimum bool            `json:"exclusiveMinimum,omitempty"`
	MaxLength        *int            `json:"maxLength,omitempty"`
	MinLength        *int            `json:"minLength,omitempty"`
	Pattern          string          `json:"pattern,omitempty"`
	MaxItems         *int            `json:"maxItems,omitempty"`
	MinItems         *int            `json:"minItems,omitempty"`
	UniqueItems      bool            `json:"uniqueItems,omitempty"`
	Enum             []any           `json:"enum,omitempty"`
	MultipleOf       *float64        `json:"multipleOf,omitempty"`
	Extensions       map[string]any  `json:"-"`
}

var parameterKnownFields = []string{
	"$ref", "name", "in", "description", "required", "schema",
	"type", "format", "allowEmptyValue", "items", "collectionFormat",
	"default", "maximum", "exclusiveMaximum", "minimum", "exclusiveMinimum",
	"maxLength", "minLength", "pattern", "maxItems", "minItems",
	"uniqueItems", "enum", "multipleOf",
}

// IsReference checks if this parameter is actually a reference ($ref)
func (p *Parameter) IsReference() bool {
	return p != nil && p.Ref != ""
}

// IsBodyParameter returns true if this is a body parameter
func (p *Parameter) IsBodyParameter() bool {
	return p != nil && p.In == "body"
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

// Items describes the type of items in an array parameter
// Used for non-body parameters with type=array
type Items struct {
	Type             string         `json:"type,omitempty"` // string, number, integer, boolean, array
	Format           string         `json:"format,omitempty"`
	Items            *Items         `json:"items,omitempty"` // for nested arrays
	CollectionFormat string         `json:"collectionFormat,omitempty"`
	Default          any            `json:"default,omitempty"`
	Maximum          *float64       `json:"maximum,omitempty"`
	ExclusiveMaximum bool           `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64       `json:"minimum,omitempty"`
	ExclusiveMinimum bool           `json:"exclusiveMinimum,omitempty"`
	MaxLength        *int           `json:"maxLength,omitempty"`
	MinLength        *int           `json:"minLength,omitempty"`
	Pattern          string         `json:"pattern,omitempty"`
	MaxItems         *int           `json:"maxItems,omitempty"`
	MinItems         *int           `json:"minItems,omitempty"`
	UniqueItems      bool           `json:"uniqueItems,omitempty"`
	Enum             []any          `json:"enum,omitempty"`
	MultipleOf       *float64       `json:"multipleOf,omitempty"`
	Extensions       map[string]any `json:"-"`
}

var itemsKnownFields = []string{
	"type", "format", "items", "collectionFormat", "default",
	"maximum", "exclusiveMaximum", "minimum", "exclusiveMinimum",
	"maxLength", "minLength", "pattern", "maxItems", "minItems",
	"uniqueItems", "enum", "multipleOf",
}

type itemsAlias Items

func (i *Items) UnmarshalJSON(data []byte) error {
	var alias itemsAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*i = Items(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	i.Extensions = extractExtensions(raw, itemsKnownFields)
	return nil
}

func (i Items) MarshalJSON() ([]byte, error) {
	alias := itemsAlias(i)
	return marshalWithExtensions(&alias, i.Extensions)
}

// Header represents a header in a response
type Header struct {
	Type             string         `json:"type"` // string, number, integer, boolean, array
	Format           string         `json:"format,omitempty"`
	Description      string         `json:"description,omitempty"`
	Items            *Items         `json:"items,omitempty"`
	CollectionFormat string         `json:"collectionFormat,omitempty"`
	Default          any            `json:"default,omitempty"`
	Maximum          *float64       `json:"maximum,omitempty"`
	ExclusiveMaximum bool           `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64       `json:"minimum,omitempty"`
	ExclusiveMinimum bool           `json:"exclusiveMinimum,omitempty"`
	MaxLength        *int           `json:"maxLength,omitempty"`
	MinLength        *int           `json:"minLength,omitempty"`
	Pattern          string         `json:"pattern,omitempty"`
	MaxItems         *int           `json:"maxItems,omitempty"`
	MinItems         *int           `json:"minItems,omitempty"`
	UniqueItems      bool           `json:"uniqueItems,omitempty"`
	Enum             []any          `json:"enum,omitempty"`
	MultipleOf       *float64       `json:"multipleOf,omitempty"`
	Extensions       map[string]any `json:"-"`
}

var headerKnownFields = []string{
	"type", "format", "description", "items", "collectionFormat",
	"default", "maximum", "exclusiveMaximum", "minimum", "exclusiveMinimum",
	"maxLength", "minLength", "pattern", "maxItems", "minItems",
	"uniqueItems", "enum", "multipleOf",
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
