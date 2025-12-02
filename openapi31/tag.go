// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import "encoding/json"

// Tag adds metadata to a single tag used by Operation
type Tag struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	Extensions   map[string]any         `json:"-"`
}

var tagKnownFields = []string{"name", "description", "externalDocs"}

type tagAlias Tag

func (t *Tag) UnmarshalJSON(data []byte) error {
	var alias tagAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*t = Tag(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	t.Extensions = extractExtensions(raw, tagKnownFields)
	return nil
}

func (t Tag) MarshalJSON() ([]byte, error) {
	alias := tagAlias(t)
	return marshalWithExtensions(&alias, t.Extensions)
}
