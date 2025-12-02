// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import "encoding/json"

// Schema represents a JSON Schema object.
// In OpenAPI 3.1, schemas are fully compatible with JSON Schema Draft 2020-12.
// This type supports both object schemas and boolean schemas (true/false).
type Schema struct {
	// Boolean schema marker
	boolValue *bool

	// Core JSON Schema keywords
	ID     string `json:"$id,omitempty"`
	Schema string `json:"$schema,omitempty"`
	Ref    string `json:"$ref,omitempty"`
	Anchor string `json:"$anchor,omitempty"`
	DynamicRef    string `json:"$dynamicRef,omitempty"`
	DynamicAnchor string `json:"$dynamicAnchor,omitempty"`
	Defs   map[string]*Schema `json:"$defs,omitempty"`
	Comment string `json:"$comment,omitempty"`

	// Vocabulary keywords
	Vocabulary map[string]bool `json:"$vocabulary,omitempty"`

	// Applicator keywords
	AllOf       []*Schema          `json:"allOf,omitempty"`
	AnyOf       []*Schema          `json:"anyOf,omitempty"`
	OneOf       []*Schema          `json:"oneOf,omitempty"`
	Not         *Schema            `json:"not,omitempty"`
	If          *Schema            `json:"if,omitempty"`
	Then        *Schema            `json:"then,omitempty"`
	Else        *Schema            `json:"else,omitempty"`
	DependentSchemas map[string]*Schema `json:"dependentSchemas,omitempty"`
	PrefixItems []*Schema          `json:"prefixItems,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
	Contains    *Schema            `json:"contains,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	PatternProperties map[string]*Schema `json:"patternProperties,omitempty"`
	AdditionalProperties *Schema   `json:"additionalProperties,omitempty"`
	PropertyNames *Schema          `json:"propertyNames,omitempty"`
	UnevaluatedItems *Schema       `json:"unevaluatedItems,omitempty"`
	UnevaluatedProperties *Schema  `json:"unevaluatedProperties,omitempty"`

	// Validation keywords - any instance type
	Type  *StringOrStringArray `json:"type,omitempty"`
	Enum  []any                `json:"enum,omitempty"`
	Const any                  `json:"const,omitempty"`

	// Validation keywords - numeric
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`

	// Validation keywords - strings
	MaxLength *int   `json:"maxLength,omitempty"`
	MinLength *int   `json:"minLength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`

	// Validation keywords - arrays
	MaxItems    *int `json:"maxItems,omitempty"`
	MinItems    *int `json:"minItems,omitempty"`
	UniqueItems bool `json:"uniqueItems,omitempty"`
	MaxContains *int `json:"maxContains,omitempty"`
	MinContains *int `json:"minContains,omitempty"`

	// Validation keywords - objects
	MaxProperties     *int              `json:"maxProperties,omitempty"`
	MinProperties     *int              `json:"minProperties,omitempty"`
	Required          []string          `json:"required,omitempty"`
	DependentRequired map[string][]string `json:"dependentRequired,omitempty"`

	// Format
	Format string `json:"format,omitempty"`

	// Content
	ContentEncoding  string  `json:"contentEncoding,omitempty"`
	ContentMediaType string  `json:"contentMediaType,omitempty"`
	ContentSchema    *Schema `json:"contentSchema,omitempty"`

	// Meta-data
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Default     any    `json:"default,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty"`
	ReadOnly    bool   `json:"readOnly,omitempty"`
	WriteOnly   bool   `json:"writeOnly,omitempty"`
	Examples    []any  `json:"examples,omitempty"`

	// OpenAPI specific
	Discriminator *Discriminator        `json:"discriminator,omitempty"`
	XML           *XML                  `json:"xml,omitempty"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs,omitempty"`
	Example       any                   `json:"example,omitempty"`

	Extensions map[string]any `json:"-"`
}

var schemaKnownFields = []string{
	"$id", "$schema", "$ref", "$anchor", "$dynamicRef", "$dynamicAnchor", "$defs", "$comment", "$vocabulary",
	"allOf", "anyOf", "oneOf", "not", "if", "then", "else", "dependentSchemas",
	"prefixItems", "items", "contains", "properties", "patternProperties",
	"additionalProperties", "propertyNames", "unevaluatedItems", "unevaluatedProperties",
	"type", "enum", "const", "multipleOf", "maximum", "exclusiveMaximum", "minimum", "exclusiveMinimum",
	"maxLength", "minLength", "pattern", "maxItems", "minItems", "uniqueItems", "maxContains", "minContains",
	"maxProperties", "minProperties", "required", "dependentRequired", "format",
	"contentEncoding", "contentMediaType", "contentSchema",
	"title", "description", "default", "deprecated", "readOnly", "writeOnly", "examples",
	"discriminator", "xml", "externalDocs", "example",
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

// Discriminator adds support for polymorphism
type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
	Extensions   map[string]any    `json:"-"`
}

var discriminatorKnownFields = []string{"propertyName", "mapping"}

type discriminatorAlias Discriminator

func (d *Discriminator) UnmarshalJSON(data []byte) error {
	var alias discriminatorAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*d = Discriminator(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	d.Extensions = extractExtensions(raw, discriminatorKnownFields)
	return nil
}

func (d Discriminator) MarshalJSON() ([]byte, error) {
	alias := discriminatorAlias(d)
	return marshalWithExtensions(&alias, d.Extensions)
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

// StringOrStringArray represents a value that can be either a string or an array of strings
type StringOrStringArray struct {
	String string
	Array  []string
}

func (s *StringOrStringArray) UnmarshalJSON(data []byte) error {
	// Try array first
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		s.Array = arr
		return nil
	}
	// Otherwise single string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	s.String = str
	return nil
}

func (s StringOrStringArray) MarshalJSON() ([]byte, error) {
	if len(s.Array) > 0 {
		return json.Marshal(s.Array)
	}
	return json.Marshal(s.String)
}

// Contains checks if the type contains the given value
func (s *StringOrStringArray) Contains(typ string) bool {
	if s == nil {
		return false
	}
	if s.String != "" {
		return s.String == typ
	}
	for _, t := range s.Array {
		if t == typ {
			return true
		}
	}
	return false
}

// IsEmpty returns true if neither String nor Array is set
func (s *StringOrStringArray) IsEmpty() bool {
	return s == nil || (s.String == "" && len(s.Array) == 0)
}
