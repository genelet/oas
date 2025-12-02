// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import "encoding/json"

// Link represents a possible design-time link for a response.
// It can also represent a Reference (when isReference is true).
type Link struct {
	// Internal marker for reference
	isReference bool

	// Reference fields (used when isReference is true)
	Ref         string `json:"$ref,omitempty"`
	Summary     string `json:"summary,omitempty"`     // Reference summary
	Description string `json:"description,omitempty"` // Shared with link description

	// Link fields
	OperationRef string            `json:"operationRef,omitempty"`
	OperationId  string            `json:"operationId,omitempty"`
	Parameters   map[string]string `json:"parameters,omitempty"`
	RequestBody  any               `json:"requestBody,omitempty"`
	Server       *Server           `json:"server,omitempty"`
	Extensions   map[string]any    `json:"-"`
}

var linkKnownFields = []string{
	"$ref", "summary", "description", "operationRef", "operationId", "parameters", "requestBody", "server",
}

// IsReference checks if this link is actually a reference ($ref)
func (l *Link) IsReference() bool {
	if l == nil {
		return false
	}
	return l.isReference
}

// NewLinkReference creates a link that is actually a reference
func NewLinkReference(ref string) *Link {
	return &Link{isReference: true, Ref: ref}
}

type linkAlias Link

type linkRefOnly struct {
	Ref         string `json:"$ref"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
}

func (l *Link) UnmarshalJSON(data []byte) error {
	var alias linkAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*l = Link(alias)
	if l.Ref != "" {
		l.isReference = true
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	l.Extensions = extractExtensions(raw, linkKnownFields)
	return nil
}

func (l Link) MarshalJSON() ([]byte, error) {
	if l.IsReference() {
		ref := linkRefOnly{
			Ref:         l.Ref,
			Summary:     l.Summary,
			Description: l.Description,
		}
		return marshalWithExtensions(&ref, l.Extensions)
	}
	alias := linkAlias(l)
	return marshalWithExtensions(&alias, l.Extensions)
}
