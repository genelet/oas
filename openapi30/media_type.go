// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi30

import "encoding/json"

// MediaType provides schema and examples for a media type
type MediaType struct {
	Schema     *Schema              `json:"schema,omitempty"`
	Example    any                  `json:"example,omitempty"`
	Examples   map[string]*Example  `json:"examples,omitempty"`
	Encoding   map[string]*Encoding `json:"encoding,omitempty"`
	Extensions map[string]any       `json:"-"`
}

var mediaTypeKnownFields = []string{"schema", "example", "examples", "encoding"}

type mediaTypeAlias MediaType

func (mt *MediaType) UnmarshalJSON(data []byte) error {
	var alias mediaTypeAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*mt = MediaType(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	mt.Extensions = extractExtensions(raw, mediaTypeKnownFields)
	return nil
}

func (mt MediaType) MarshalJSON() ([]byte, error) {
	alias := mediaTypeAlias(mt)
	return marshalWithExtensions(&alias, mt.Extensions)
}

// Encoding defines encoding for a single property
type Encoding struct {
	ContentType   string             `json:"contentType,omitempty"`
	Headers       map[string]*Header `json:"headers,omitempty"`
	Style         string             `json:"style,omitempty"`
	Explode       *bool              `json:"explode,omitempty"`
	AllowReserved bool               `json:"allowReserved,omitempty"`
	Extensions    map[string]any     `json:"-"`
}

var encodingKnownFields = []string{"contentType", "headers", "style", "explode", "allowReserved"}

type encodingAlias Encoding

func (e *Encoding) UnmarshalJSON(data []byte) error {
	var alias encodingAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*e = Encoding(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Extensions = extractExtensions(raw, encodingKnownFields)
	return nil
}

func (e Encoding) MarshalJSON() ([]byte, error) {
	alias := encodingAlias(e)
	return marshalWithExtensions(&alias, e.Extensions)
}
