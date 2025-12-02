// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import (
	"encoding/json"
	"strings"
)

// Paths holds the relative paths to the individual endpoints and their operations
type Paths struct {
	Paths      map[string]*PathItem `json:"-"`
	Extensions map[string]any       `json:"-"`
}

// Get returns the PathItem for the given path
func (p *Paths) Get(path string) *PathItem {
	if p == nil || p.Paths == nil {
		return nil
	}
	return p.Paths[path]
}

// Set sets the PathItem for the given path
func (p *Paths) Set(path string, item *PathItem) {
	if p.Paths == nil {
		p.Paths = make(map[string]*PathItem)
	}
	p.Paths[path] = item
}

func (p *Paths) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	p.Paths = make(map[string]*PathItem)
	p.Extensions = make(map[string]any)

	for key, value := range raw {
		if strings.HasPrefix(key, "x-") {
			var ext any
			if err := json.Unmarshal(value, &ext); err != nil {
				return err
			}
			p.Extensions[key] = ext
		} else if strings.HasPrefix(key, "/") {
			var pathItem PathItem
			if err := json.Unmarshal(value, &pathItem); err != nil {
				return err
			}
			p.Paths[key] = &pathItem
		}
	}

	if len(p.Extensions) == 0 {
		p.Extensions = nil
	}
	return nil
}

func (p Paths) MarshalJSON() ([]byte, error) {
	result := make(map[string]any)
	for key, value := range p.Paths {
		result[key] = value
	}
	for key, value := range p.Extensions {
		result[key] = value
	}
	return json.Marshal(result)
}

// PathItem describes the operations available on a single path
type PathItem struct {
	Ref         string         `json:"$ref,omitempty"`
	Summary     string         `json:"summary,omitempty"`
	Description string         `json:"description,omitempty"`
	Servers     []*Server      `json:"servers,omitempty"`
	Parameters  []*Parameter   `json:"parameters,omitempty"`
	Get         *Operation     `json:"get,omitempty"`
	Put         *Operation     `json:"put,omitempty"`
	Post        *Operation     `json:"post,omitempty"`
	Delete      *Operation     `json:"delete,omitempty"`
	Options     *Operation     `json:"options,omitempty"`
	Head        *Operation     `json:"head,omitempty"`
	Patch       *Operation     `json:"patch,omitempty"`
	Trace       *Operation     `json:"trace,omitempty"`
	Extensions  map[string]any `json:"-"`
}

var pathItemKnownFields = []string{
	"$ref", "summary", "description", "servers", "parameters",
	"get", "put", "post", "delete", "options", "head", "patch", "trace",
}

type pathItemAlias PathItem

func (pi *PathItem) UnmarshalJSON(data []byte) error {
	var alias pathItemAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*pi = PathItem(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	pi.Extensions = extractExtensions(raw, pathItemKnownFields)
	return nil
}

func (pi PathItem) MarshalJSON() ([]byte, error) {
	alias := pathItemAlias(pi)
	return marshalWithExtensions(&alias, pi.Extensions)
}

// HasRef returns true if this PathItem has a $ref
func (pi *PathItem) HasRef() bool {
	return pi != nil && pi.Ref != ""
}

// GetRef returns the $ref value
func (pi *PathItem) GetRef() string {
	if pi == nil {
		return ""
	}
	return pi.Ref
}
