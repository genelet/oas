// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

// Package openapi30 provides Go types for OpenAPI Specification v3.0.x
// Based on https://spec.openapis.org/oas/3.0/schema/2024-10-18
package openapi30

import "encoding/json"

// OpenAPI is the root object of an OpenAPI v3.0.x document
type OpenAPI struct {
	OpenAPI      string                 `json:"openapi"`
	Info         *Info                  `json:"info"`
	Servers      []*Server              `json:"servers,omitempty"`
	Paths        *Paths                 `json:"paths"`
	Components   *Components            `json:"components,omitempty"`
	Security     []SecurityRequirement  `json:"security,omitempty"`
	Tags         []*Tag                 `json:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	Extensions   map[string]any         `json:"-"`
}

var openapiKnownFields = []string{
	"openapi", "info", "servers", "paths", "components", "security", "tags", "externalDocs",
}

type openapiAlias OpenAPI

func (o *OpenAPI) UnmarshalJSON(data []byte) error {
	var alias openapiAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*o = OpenAPI(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	o.Extensions = extractExtensions(raw, openapiKnownFields)
	return nil
}

func (o OpenAPI) MarshalJSON() ([]byte, error) {
	alias := openapiAlias(o)
	return marshalWithExtensions(&alias, o.Extensions)
}
