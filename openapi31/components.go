// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import "encoding/json"

// Components holds a set of reusable objects for different aspects of the OAS
type Components struct {
	Schemas         map[string]*Schema          `json:"schemas,omitempty"`
	Responses       map[string]*Response        `json:"responses,omitempty"`
	Parameters      map[string]*Parameter       `json:"parameters,omitempty"`
	Examples        map[string]*Example         `json:"examples,omitempty"`
	RequestBodies   map[string]*RequestBody     `json:"requestBodies,omitempty"`
	Headers         map[string]*Header          `json:"headers,omitempty"`
	SecuritySchemes map[string]*SecurityScheme  `json:"securitySchemes,omitempty"`
	Links           map[string]*Link            `json:"links,omitempty"`
	Callbacks       map[string]*Callback        `json:"callbacks,omitempty"`
	PathItems       map[string]*PathItem        `json:"pathItems,omitempty"`
	Extensions      map[string]any              `json:"-"`
}

var componentsKnownFields = []string{
	"schemas", "responses", "parameters", "examples", "requestBodies",
	"headers", "securitySchemes", "links", "callbacks", "pathItems",
}

type componentsAlias Components

func (c *Components) UnmarshalJSON(data []byte) error {
	var alias componentsAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*c = Components(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	c.Extensions = extractExtensions(raw, componentsKnownFields)
	return nil
}

func (c Components) MarshalJSON() ([]byte, error) {
	alias := componentsAlias(c)
	return marshalWithExtensions(&alias, c.Extensions)
}
