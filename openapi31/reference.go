// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

// Reference is a simple object to allow referencing other components in the OpenAPI document
type Reference struct {
	Ref         string `json:"$ref"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
}
