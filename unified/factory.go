// Package unified provides unified interfaces for OpenAPI documents
// Copyright (c) Genelet

package unified

import (
	"encoding/json"
	"fmt"
	"strings"

	oa2 "github.com/genelet/oas/openapi20"
	oa3 "github.com/genelet/oas/openapi30"
	oa31 "github.com/genelet/oas/openapi31"
)

// versionProbe is used to detect OpenAPI version
type versionProbe struct {
	Swagger string `json:"swagger"`
	OpenAPI string `json:"openapi"`
}

// NewDocument parses a JSON-encoded OpenAPI document and returns a unified Document interface.
// It automatically detects the version (2.0, 3.0.x, 3.1.x) and returns the appropriate adapter.
// Note: This function only accepts JSON. YAML must be converted to JSON before calling this.
func NewDocument(data []byte) (Document, error) {
	// Detect version
	var probe versionProbe
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("failed to detect OpenAPI version: %w", err)
	}

	switch {
	case probe.Swagger == "2.0":
		return parseOpenAPI20(data)

	case strings.HasPrefix(probe.OpenAPI, "3.0"):
		return parseOpenAPI30(data)

	case strings.HasPrefix(probe.OpenAPI, "3.1"):
		return parseOpenAPI31(data)

	default:
		return nil, fmt.Errorf("unsupported OpenAPI version: swagger=%q openapi=%q",
			probe.Swagger, probe.OpenAPI)
	}
}

// parseOpenAPI30 parses an OpenAPI 3.0 document
func parseOpenAPI30(jsonData []byte) (Document, error) {
	var doc oa3.OpenAPI
	if err := json.Unmarshal(jsonData, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI 3.0: %w", err)
	}
	return NewDocument30(&doc), nil
}

// parseOpenAPI20 parses an OpenAPI 2.0 (Swagger) document
func parseOpenAPI20(jsonData []byte) (Document, error) {
	var doc oa2.Swagger
	if err := json.Unmarshal(jsonData, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI 2.0: %w", err)
	}
	return NewDocument20(&doc), nil
}

// parseOpenAPI31 parses an OpenAPI 3.1 document
func parseOpenAPI31(jsonData []byte) (Document, error) {
	var doc oa31.OpenAPI
	if err := json.Unmarshal(jsonData, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI 3.1: %w", err)
	}
	return NewDocument31(&doc), nil
}
