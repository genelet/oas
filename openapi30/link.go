// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi30

import "encoding/json"

// Link represents a possible design-time link for a response.
// It can also represent a Reference (when Ref is set).
type Link struct {
	// Reference field
	Ref string `json:"$ref,omitempty"`

	// Link fields
	OperationRef string         `json:"operationRef,omitempty"`
	OperationId  string         `json:"operationId,omitempty"`
	Parameters   map[string]any `json:"parameters,omitempty"`
	RequestBody  any            `json:"requestBody,omitempty"`
	Description  string         `json:"description,omitempty"`
	Server       *Server        `json:"server,omitempty"`
	Extensions   map[string]any `json:"-"`
}

var linkKnownFields = []string{
	"$ref", "operationRef", "operationId", "parameters", "requestBody", "description", "server",
}

// IsReference checks if this link is actually a reference ($ref)
func (l *Link) IsReference() bool {
	return l != nil && l.Ref != ""
}

// NewLinkReference creates a link that is actually a reference
func NewLinkReference(ref string) *Link {
	return &Link{Ref: ref}
}

type linkAlias Link

func (l *Link) UnmarshalJSON(data []byte) error {
	var alias linkAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*l = Link(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	l.Extensions = extractExtensions(raw, linkKnownFields)
	return nil
}

func (l Link) MarshalJSON() ([]byte, error) {
	alias := linkAlias(l)
	return marshalWithExtensions(&alias, l.Extensions)
}
