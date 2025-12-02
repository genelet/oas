// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import (
	"encoding/json"
	"strings"
)

// Callback is a map of possible out-of band callbacks related to the parent operation.
// It can also represent a Reference (when isReference is true).
type Callback struct {
	// Internal marker for reference
	isReference bool

	// Reference fields (used when isReference is true)
	Ref         string `json:"-"`
	Summary     string `json:"-"` // Reference summary
	Description string `json:"-"` // Reference description

	// Callback fields
	Paths      map[string]*PathItem `json:"-"`
	Extensions map[string]any       `json:"-"`
}

// IsReference checks if this callback is actually a reference ($ref)
func (c *Callback) IsReference() bool {
	if c == nil {
		return false
	}
	return c.isReference
}

// NewCallbackReference creates a callback that is actually a reference
func NewCallbackReference(ref string) *Callback {
	return &Callback{isReference: true, Ref: ref}
}

type callbackRefOnly struct {
	Ref         string `json:"$ref"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *Callback) UnmarshalJSON(data []byte) error {
	// Check if this is a reference first
	var ref struct {
		Ref string `json:"$ref"`
	}
	if err := json.Unmarshal(data, &ref); err == nil && ref.Ref != "" {
		c.isReference = true
		var refOnly callbackRefOnly
		if err := json.Unmarshal(data, &refOnly); err != nil {
			return err
		}
		c.Ref = refOnly.Ref
		c.Summary = refOnly.Summary
		c.Description = refOnly.Description
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
		ref := callbackRefOnly{
			Ref:         c.Ref,
			Summary:     c.Summary,
			Description: c.Description,
		}
		return marshalWithExtensions(&ref, c.Extensions)
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
