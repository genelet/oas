// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import "encoding/json"

// Operation describes a single API operation on a path
type Operation struct {
	Tags         []string               `json:"tags,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	Description  string                 `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	OperationID  string                 `json:"operationId,omitempty"`
	Parameters   []*Parameter           `json:"parameters,omitempty"`
	RequestBody  *RequestBody           `json:"requestBody,omitempty"`
	Responses    *Responses             `json:"responses,omitempty"`
	Callbacks    map[string]*Callback   `json:"callbacks,omitempty"`
	Deprecated   bool                   `json:"deprecated,omitempty"`
	Security     []SecurityRequirement  `json:"security,omitempty"`
	Servers      []*Server              `json:"servers,omitempty"`
	Extensions   map[string]any         `json:"-"`
}

var operationKnownFields = []string{
	"tags", "summary", "description", "externalDocs", "operationId",
	"parameters", "requestBody", "responses", "callbacks", "deprecated",
	"security", "servers",
}

type operationAlias Operation

func (o *Operation) UnmarshalJSON(data []byte) error {
	var alias operationAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*o = Operation(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	o.Extensions = extractExtensions(raw, operationKnownFields)
	return nil
}

func (o Operation) MarshalJSON() ([]byte, error) {
	alias := operationAlias(o)
	return marshalWithExtensions(&alias, o.Extensions)
}
