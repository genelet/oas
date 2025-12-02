// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

// Package openapi31 provides Go types for OpenAPI Specification v3.1.x
// Generated from https://spec.openapis.org/oas/3.1/schema/2025-09-15
package openapi31

import "encoding/json"

// OpenAPI is the root object of an OpenAPI v3.1.x document
type OpenAPI struct {
	OpenAPI           string                 `json:"openapi"`
	Info              *Info                  `json:"info"`
	JsonSchemaDialect string                 `json:"jsonSchemaDialect,omitempty"`
	Servers           []*Server              `json:"servers,omitempty"`
	Paths             *Paths                 `json:"paths,omitempty"`
	Webhooks          map[string]*PathItem   `json:"webhooks,omitempty"`
	Components        *Components            `json:"components,omitempty"`
	Security          []SecurityRequirement  `json:"security,omitempty"`
	Tags              []*Tag                 `json:"tags,omitempty"`
	ExternalDocs      *ExternalDocumentation `json:"externalDocs,omitempty"`
	Extensions        map[string]any         `json:"-"`
}

var openapiKnownFields = []string{
	"openapi", "info", "jsonSchemaDialect", "servers", "paths", "webhooks",
	"components", "security", "tags", "externalDocs",
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
