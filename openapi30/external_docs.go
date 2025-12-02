// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi30

import "encoding/json"

// ExternalDocumentation allows referencing an external resource for extended documentation
type ExternalDocumentation struct {
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url"`
	Extensions  map[string]any `json:"-"`
}

var externalDocsKnownFields = []string{"description", "url"}

type externalDocsAlias ExternalDocumentation

func (ed *ExternalDocumentation) UnmarshalJSON(data []byte) error {
	var alias externalDocsAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*ed = ExternalDocumentation(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	ed.Extensions = extractExtensions(raw, externalDocsKnownFields)
	return nil
}

func (ed ExternalDocumentation) MarshalJSON() ([]byte, error) {
	alias := externalDocsAlias(ed)
	return marshalWithExtensions(&alias, ed.Extensions)
}
