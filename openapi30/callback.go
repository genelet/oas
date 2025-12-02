// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi30

import (
	"encoding/json"
	"strings"
)

// Callback is a map of possible out-of band callbacks related to the parent operation.
// It can also represent a Reference (when Ref is set).
type Callback struct {
	// Reference field
	Ref string `json:"-"`

	// Callback fields
	Paths      map[string]*PathItem `json:"-"`
	Extensions map[string]any       `json:"-"`
}

// IsReference checks if this callback is actually a reference ($ref)
func (c *Callback) IsReference() bool {
	return c != nil && c.Ref != ""
}

// NewCallbackReference creates a callback that is actually a reference
func NewCallbackReference(ref string) *Callback {
	return &Callback{Ref: ref}
}

type callbackRefOnly struct {
	Ref string `json:"$ref"`
}

func (c *Callback) UnmarshalJSON(data []byte) error {
	// Check if this is a reference first
	var ref struct {
		Ref string `json:"$ref"`
	}
	if err := json.Unmarshal(data, &ref); err == nil && ref.Ref != "" {
		c.Ref = ref.Ref
		return nil
	}

	// Otherwise unmarshal as callback (map of paths)
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.Paths = make(map[string]*PathItem)
	c.Extensions = make(map[string]any)

	for key, value := range raw {
		if strings.HasPrefix(key, "x-") {
			var ext any
			if err := json.Unmarshal(value, &ext); err != nil {
				return err
			}
			c.Extensions[key] = ext
		} else {
			var pathItem PathItem
			if err := json.Unmarshal(value, &pathItem); err != nil {
				return err
			}
			c.Paths[key] = &pathItem
		}
	}

	if len(c.Extensions) == 0 {
		c.Extensions = nil
	}
	return nil
}

func (c Callback) MarshalJSON() ([]byte, error) {
	if c.IsReference() {
		return json.Marshal(callbackRefOnly{Ref: c.Ref})
	}

	result := make(map[string]any)
	for key, value := range c.Paths {
		result[key] = value
	}
	for key, value := range c.Extensions {
		result[key] = value
	}
	return json.Marshal(result)
}

// Get returns the PathItem for the given expression
func (c *Callback) Get(expression string) *PathItem {
	if c == nil || c.Paths == nil {
		return nil
	}
	return c.Paths[expression]
}
