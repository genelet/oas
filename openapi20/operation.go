// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi20

import "encoding/json"

// Operation describes a single API operation on a path
type Operation struct {
	Tags         []string               `json:"tags,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	Description  string                 `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	OperationID  string                 `json:"operationId,omitempty"`
	Consumes     []string               `json:"consumes,omitempty"`
	Produces     []string               `json:"produces,omitempty"`
	Parameters   []*Parameter           `json:"parameters,omitempty"`
	Responses    *Responses             `json:"responses"`
	Schemes      []string               `json:"schemes,omitempty"`
	Deprecated   bool                   `json:"deprecated,omitempty"`
	Security     []SecurityRequirement  `json:"security,omitempty"`
	Extensions   map[string]any         `json:"-"`
}

var operationKnownFields = []string{
	"tags", "summary", "description", "externalDocs", "operationId",
	"consumes", "produces", "parameters", "responses", "schemes",
	"deprecated", "security",
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
