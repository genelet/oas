// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi30

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Response describes a single response from an API operation.
// It can also represent a Reference (when Ref is set).
type Response struct {
	// Reference field
	Ref string `json:"$ref,omitempty"`

	// Response fields
	Description string                `json:"description,omitempty"`
	Headers     map[string]*Header    `json:"headers,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
	Links       map[string]*Link      `json:"links,omitempty"`
	Extensions  map[string]any        `json:"-"`
}

var responseKnownFields = []string{"$ref", "description", "headers", "content", "links"}

// IsReference checks if this response is actually a reference ($ref)
func (r *Response) IsReference() bool {
	return r != nil && r.Ref != ""
}

// NewResponseReference creates a response that is actually a reference
func NewResponseReference(ref string) *Response {
	return &Response{Ref: ref}
}

type responseAlias Response

func (r *Response) UnmarshalJSON(data []byte) error {
	var alias responseAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*r = Response(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.Extensions = extractExtensions(raw, responseKnownFields)
	return nil
}

func (r Response) MarshalJSON() ([]byte, error) {
	alias := responseAlias(r)
	return marshalWithExtensions(&alias, r.Extensions)
}

// Responses is a container for the expected responses of an operation
type Responses struct {
	Default    *Response            `json:"-"`
	StatusCode map[string]*Response `json:"-"` // HTTP status codes (e.g., "200", "4XX")
	Extensions map[string]any       `json:"-"`
}

var statusCodePattern = regexp.MustCompile(`^[1-5](?:[0-9]{2}|XX)$`)

func (r *Responses) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.StatusCode = make(map[string]*Response)
	r.Extensions = make(map[string]any)

	for key, value := range raw {
		if key == "default" {
			r.Default = &Response{}
			if err := json.Unmarshal(value, r.Default); err != nil {
				return err
			}
		} else if statusCodePattern.MatchString(key) {
			resp := &Response{}
			if err := json.Unmarshal(value, resp); err != nil {
				return err
			}
			r.StatusCode[key] = resp
		} else if strings.HasPrefix(key, "x-") {
			var ext any
			if err := json.Unmarshal(value, &ext); err != nil {
				return err
			}
			r.Extensions[key] = ext
		}
	}

	if len(r.Extensions) == 0 {
		r.Extensions = nil
	}
	return nil
}

func (r Responses) MarshalJSON() ([]byte, error) {
	result := make(map[string]any)
	if r.Default != nil {
		result["default"] = r.Default
	}
	for key, value := range r.StatusCode {
		result[key] = value
	}
	for key, value := range r.Extensions {
		result[key] = value
	}
	return json.Marshal(result)
}

// Get returns the Response for the given status code
func (r *Responses) Get(statusCode string) *Response {
	if r == nil {
		return nil
	}
	if statusCode == "default" {
		return r.Default
	}
	if r.StatusCode == nil {
		return nil
	}
	return r.StatusCode[statusCode]
}

// GetDefault returns the default response
func (r *Responses) GetDefault() *Response {
	if r == nil {
		return nil
	}
	return r.Default
}
