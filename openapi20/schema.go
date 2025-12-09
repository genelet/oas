// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi20

import "encoding/json"

// Schema represents a JSON Schema subset for Swagger 2.0.
// This type supports both object schemas and boolean schemas (true/false).
// In Swagger 2.0, additionalProperties can be a boolean or a schema object.
type Schema struct {
	// Boolean schema marker (private)
	boolValue *bool

	// Reference
	Ref string `json:"$ref,omitempty"`

	// JSON Schema keywords
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Default     any    `json:"default,omitempty"`
	Format      string `json:"format,omitempty"`

	// Type - can be string or array of strings in Swagger 2.0
	Type string `json:"type,omitempty"`

	// Enum and validation
	Enum []any `json:"enum,omitempty"`

	// Numeric validation
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMaximum bool     `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	ExclusiveMinimum bool     `json:"exclusiveMinimum,omitempty"`

	// String validation
	MaxLength *int   `json:"maxLength,omitempty"`
	MinLength *int   `json:"minLength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`

	// Array validation
	MaxItems    *int    `json:"maxItems,omitempty"`
	MinItems    *int    `json:"minItems,omitempty"`
	UniqueItems bool    `json:"uniqueItems,omitempty"`
	Items       *Schema `json:"items,omitempty"`

	// Object validation
	MaxProperties        *int               `json:"maxProperties,omitempty"`
	MinProperties        *int               `json:"minProperties,omitempty"`
	Required             []string           `json:"required,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	AdditionalProperties *Schema            `json:"additionalProperties,omitempty"` // Can be boolean schema

	// Composition keywords - Swagger 2.0 only supports allOf
	AllOf []*Schema `json:"allOf,omitempty"`

	// Swagger 2.0 specific
	Discriminator string                 `json:"discriminator,omitempty"` // Property name (string in 2.0, object in 3.0)
	ReadOnly      bool                   `json:"readOnly,omitempty"`
	XML           *XML                   `json:"xml,omitempty"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs,omitempty"`
	Example       any                    `json:"example,omitempty"`

	Extensions map[string]any `json:"-"`
}

var schemaKnownFields = []string{
	"$ref", "title", "description", "default", "format", "type", "enum",
	"multipleOf", "maximum", "exclusiveMaximum", "minimum", "exclusiveMinimum",
	"maxLength", "minLength", "pattern",
	"maxItems", "minItems", "uniqueItems", "items",
	"maxProperties", "minProperties", "required", "properties", "additionalProperties",
	"allOf", "discriminator", "readOnly", "xml", "externalDocs", "example",
}

type schemaAlias Schema

func (s *Schema) UnmarshalJSON(data []byte) error {
	// Try boolean first
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		s.boolValue = &b
		return nil
	}

	// Otherwise unmarshal as object
	var alias schemaAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*s = Schema(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	s.Extensions = extractExtensions(raw, schemaKnownFields)
	return nil
}

func (s Schema) MarshalJSON() ([]byte, error) {
	// Handle boolean schema
	if s.boolValue != nil {
		return json.Marshal(*s.boolValue)
	}

	alias := schemaAlias(s)
	return marshalWithExtensions(&alias, s.Extensions)
}

// IsBooleanSchema returns true if this is a boolean schema (true or false)
func (s *Schema) IsBooleanSchema() bool {
	return s != nil && s.boolValue != nil
}

// BooleanValue returns the boolean value if this is a boolean schema
func (s *Schema) BooleanValue() *bool {
	if s == nil {
		return nil
	}
	return s.boolValue
}

// NewBooleanSchema creates a boolean schema
func NewBooleanSchema(value bool) *Schema {
	return &Schema{boolValue: &value}
}

// IsReference checks if this schema is actually a reference ($ref)
func (s *Schema) IsReference() bool {
	return s != nil && s.Ref != ""
}

// NewSchemaReference creates a schema that is actually a reference
func NewSchemaReference(ref string) *Schema {
	return &Schema{Ref: ref}
}

// XML provides metadata for XML representation
type XML struct {
	Name       string         `json:"name,omitempty"`
	Namespace  string         `json:"namespace,omitempty"`
	Prefix     string         `json:"prefix,omitempty"`
	Attribute  bool           `json:"attribute,omitempty"`
	Wrapped    bool           `json:"wrapped,omitempty"`
	Extensions map[string]any `json:"-"`
}

var xmlKnownFields = []string{"name", "namespace", "prefix", "attribute", "wrapped"}

type xmlAlias XML

func (x *XML) UnmarshalJSON(data []byte) error {
	var alias xmlAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*x = XML(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	x.Extensions = extractExtensions(raw, xmlKnownFields)
	return nil
}

func (x XML) MarshalJSON() ([]byte, error) {
	alias := xmlAlias(x)
	return marshalWithExtensions(&alias, x.Extensions)
}
