// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

// Package openapi20 provides Go types for Swagger/OpenAPI Specification v2.0
// Based on https://swagger.io/specification/v2/
package openapi20

import "encoding/json"

// Swagger is the root object of a Swagger 2.0 document
type Swagger struct {
	Swagger             string                      `json:"swagger"`
	Info                *Info                       `json:"info"`
	Host                string                      `json:"host,omitempty"`
	BasePath            string                      `json:"basePath,omitempty"`
	Schemes             []string                    `json:"schemes,omitempty"`
	Consumes            []string                    `json:"consumes,omitempty"`
	Produces            []string                    `json:"produces,omitempty"`
	Paths               *Paths                      `json:"paths"`
	Definitions         map[string]*Schema          `json:"definitions,omitempty"`
	Parameters          map[string]*Parameter       `json:"parameters,omitempty"`
	Responses           map[string]*Response        `json:"responses,omitempty"`
	SecurityDefinitions map[string]*SecurityScheme  `json:"securityDefinitions,omitempty"`
	Security            []SecurityRequirement       `json:"security,omitempty"`
	Tags                []*Tag                      `json:"tags,omitempty"`
	ExternalDocs        *ExternalDocumentation      `json:"externalDocs,omitempty"`
	Extensions          map[string]any              `json:"-"`
}

var swaggerKnownFields = []string{
	"swagger", "info", "host", "basePath", "schemes", "consumes", "produces",
	"paths", "definitions", "parameters", "responses", "securityDefinitions",
	"security", "tags", "externalDocs",
}

type swaggerAlias Swagger

func (s *Swagger) UnmarshalJSON(data []byte) error {
	var alias swaggerAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*s = Swagger(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	s.Extensions = extractExtensions(raw, swaggerKnownFields)
	return nil
}

func (s Swagger) MarshalJSON() ([]byte, error) {
	alias := swaggerAlias(s)
	return marshalWithExtensions(&alias, s.Extensions)
}
