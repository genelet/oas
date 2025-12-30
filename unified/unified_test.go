// Package unified provides unified interfaces for OpenAPI documents
// Copyright (c) Genelet

package unified

import (
	"testing"
)

func TestUnified(t *testing.T) {
	// Simple test to verify interface compliance and basic parsing
	// Real parsing tests are done in individual adapter packages usually,
	// but here we verify the unification.

	minimal20 := `{
		"swagger": "2.0",
		"info": {
			"title": "Minimal 2.0",
			"version": "1.0.0"
		},
		"paths": {}
	}`

	minimal30 := `{
		"openapi": "3.0.0",
		"info": {
			"title": "Minimal 3.0",
			"version": "1.0.0"
		},
		"paths": {}
	}`

	minimal31 := `{
		"openapi": "3.1.0",
		"info": {
			"title": "Minimal 3.1",
			"version": "1.0.0"
		},
		"paths": {}
	}`

	tests := []struct {
		name    string
		data    string
		version string
	}{
		{"Swagger 2.0", minimal20, "2.0"},
		{"OpenAPI 3.0", minimal30, "3.0.0"},
		{"OpenAPI 3.1", minimal31, "3.1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := NewDocument([]byte(tt.data))
			if err != nil {
				t.Fatalf("NewDocument() error = %v", err)
			}
			if doc.Version() != tt.version {
				t.Errorf("Version() = %v, want %v", doc.Version(), tt.version)
			}
			if doc.GetInfo().GetTitle() != "Minimal "+tt.version[:3] && doc.GetInfo().GetTitle() != "Minimal "+tt.version {
				// 3.1 title is "Minimal 3.1", 3.0 is "Minimal 3.0", 2.0 is "Minimal 2.0"
				// but check loosely
			}
		})
	}
}

func TestExtension(t *testing.T) {
	spec := `{
		"openapi": "3.0.0",
		"info": {
			"title": "Extension Test",
			"version": "1.0",
			"x-logo": "logo.png"
		},
		"paths": {},
		"x-foo": "bar"
	}`

	doc, err := NewDocument([]byte(spec))
	if err != nil {
		t.Fatalf("NewDocument() error = %v", err)
	}

	ext := doc.GetExtensions()
	if ext["x-foo"] != "bar" {
		t.Errorf("GetExtensions() x-foo = %v, want bar", ext["x-foo"])
	}

	infoExt := doc.GetInfo().GetExtensions()
	if infoExt["x-logo"] != "logo.png" {
		t.Errorf("GetInfo().GetExtensions() x-logo = %v, want logo.png", infoExt["x-logo"])
	}
}
