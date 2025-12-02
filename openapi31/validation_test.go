// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidateMinimalValid(t *testing.T) {
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info: &Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: &Paths{
			Paths: map[string]*PathItem{
				"/test": {
					Get: &Operation{
						Responses: &Responses{
							StatusCode: map[string]*Response{
								"200": {Description: "OK"},
							},
						},
					},
				},
			},
		},
	}

	result := api.Validate()
	if !result.Valid() {
		t.Errorf("Expected valid document, got errors: %v", result.Error())
	}
}

func TestValidateNilDocument(t *testing.T) {
	var api *OpenAPI
	result := api.Validate()
	if result.Valid() {
		t.Error("Expected error for nil document")
	}
}

func TestValidateMissingRequiredFields(t *testing.T) {
	tests := []struct {
		name        string
		api         *OpenAPI
		expectError string
	}{
		{
			name:        "missing openapi",
			api:         &OpenAPI{Info: &Info{Title: "Test", Version: "1.0"}, Paths: &Paths{}},
			expectError: "openapi",
		},
		{
			name:        "missing info",
			api:         &OpenAPI{OpenAPI: "3.1.0", Paths: &Paths{}},
			expectError: "info",
		},
		{
			name: "missing info.title",
			api: &OpenAPI{
				OpenAPI: "3.1.0",
				Info:    &Info{Version: "1.0"},
				Paths:   &Paths{},
			},
			expectError: "info.title",
		},
		{
			name: "missing info.version",
			api: &OpenAPI{
				OpenAPI: "3.1.0",
				Info:    &Info{Title: "Test"},
				Paths:   &Paths{},
			},
			expectError: "info.version",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.api.Validate()
			if result.Valid() {
				t.Error("Expected validation error")
				return
			}
			found := false
			for _, err := range result.Errors {
				if strings.Contains(err.Path, tc.expectError) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error containing '%s', got: %v", tc.expectError, result.Error())
			}
		})
	}
}

func TestValidatePathsWebhooksOrComponents(t *testing.T) {
	// OpenAPI 3.1 requires at least one of: paths, webhooks, or components
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info:    &Info{Title: "Test", Version: "1.0"},
		// No paths, webhooks, or components
	}

	result := api.Validate()
	found := false
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "paths, webhooks, or components") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error about needing paths, webhooks, or components")
	}

	// With only webhooks should be valid
	api.Webhooks = map[string]*PathItem{
		"newPet": {
			Post: &Operation{
				Responses: &Responses{
					StatusCode: map[string]*Response{
						"200": {Description: "OK"},
					},
				},
			},
		},
	}
	result = api.Validate()
	if !result.Valid() {
		t.Errorf("Expected valid with webhooks, got: %v", result.Error())
	}
}

func TestValidateLicenseMutualExclusion(t *testing.T) {
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info: &Info{
			Title:   "Test",
			Version: "1.0",
			License: &License{
				Name:       "MIT",
				Identifier: "MIT",
				URL:        "https://opensource.org/licenses/MIT",
			},
		},
		Paths: &Paths{Paths: map[string]*PathItem{"/": {}}},
	}

	result := api.Validate()
	found := false
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "mutually exclusive") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error about identifier and url being mutually exclusive")
	}
}

func TestValidatePathParameterRequired(t *testing.T) {
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info:    &Info{Title: "Test", Version: "1.0"},
		Paths: &Paths{
			Paths: map[string]*PathItem{
				"/users/{id}": {
					Get: &Operation{
						Parameters: []*Parameter{
							{
								Name:     "id",
								In:       "path",
								Required: false, // This should fail
								Schema:   &Schema{Type: &StringOrStringArray{String: "string"}},
							},
						},
						Responses: &Responses{
							StatusCode: map[string]*Response{
								"200": {Description: "OK"},
							},
						},
					},
				},
			},
		},
	}

	result := api.Validate()
	found := false
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "path parameters must have required: true") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'path parameters must have required: true' error, got: %v", result.Error())
	}
}

func TestValidateSchemaXORContent(t *testing.T) {
	// Parameter with neither schema nor content
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info:    &Info{Title: "Test", Version: "1.0"},
		Paths: &Paths{
			Paths: map[string]*PathItem{
				"/test": {
					Get: &Operation{
						Parameters: []*Parameter{
							{Name: "q", In: "query"}, // Missing both schema and content
						},
						Responses: &Responses{
							StatusCode: map[string]*Response{
								"200": {Description: "OK"},
							},
						},
					},
				},
			},
		},
	}

	result := api.Validate()
	found := false
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "must have either 'schema' or 'content'") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for parameter without schema or content")
	}
}

func TestValidateLinkOperationXOR(t *testing.T) {
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info:    &Info{Title: "Test", Version: "1.0"},
		Paths: &Paths{
			Paths: map[string]*PathItem{
				"/test": {
					Get: &Operation{
						Responses: &Responses{
							StatusCode: map[string]*Response{
								"200": {
									Description: "OK",
									Links: map[string]*Link{
										"test": {
											OperationId:  "getUser",
											OperationRef: "#/paths/~1users/get", // Both set - invalid
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	result := api.Validate()
	found := false
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "cannot have both 'operationId' and 'operationRef'") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for link with both operationId and operationRef")
	}
}

func TestValidateResponseDescription(t *testing.T) {
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info:    &Info{Title: "Test", Version: "1.0"},
		Paths: &Paths{
			Paths: map[string]*PathItem{
				"/test": {
					Get: &Operation{
						Responses: &Responses{
							StatusCode: map[string]*Response{
								"200": {}, // Missing description
							},
						},
					},
				},
			},
		},
	}

	result := api.Validate()
	found := false
	for _, err := range result.Errors {
		if strings.Contains(err.Path, "description") && strings.Contains(err.Message, "required") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for response without description")
	}
}

func TestValidateSecurityScheme(t *testing.T) {
	tests := []struct {
		name        string
		scheme      *SecurityScheme
		expectError string
	}{
		{
			name:        "apiKey missing name",
			scheme:      &SecurityScheme{Type: "apiKey", In: "header"},
			expectError: "name",
		},
		{
			name:        "apiKey missing in",
			scheme:      &SecurityScheme{Type: "apiKey", Name: "X-API-Key"},
			expectError: "in",
		},
		{
			name:        "http missing scheme",
			scheme:      &SecurityScheme{Type: "http"},
			expectError: "scheme",
		},
		{
			name:        "oauth2 missing flows",
			scheme:      &SecurityScheme{Type: "oauth2"},
			expectError: "flows",
		},
		{
			name:        "openIdConnect missing url",
			scheme:      &SecurityScheme{Type: "openIdConnect"},
			expectError: "openIdConnectUrl",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &OpenAPI{
				OpenAPI: "3.1.0",
				Info:    &Info{Title: "Test", Version: "1.0"},
				Paths:   &Paths{Paths: map[string]*PathItem{"/": {}}},
				Components: &Components{
					SecuritySchemes: map[string]*SecurityScheme{
						"test": tc.scheme,
					},
				},
			}

			result := api.Validate()
			found := false
			for _, err := range result.Errors {
				if strings.Contains(err.Path, tc.expectError) || strings.Contains(err.Message, tc.expectError) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error about '%s', got: %v", tc.expectError, result.Error())
			}
		})
	}
}

func TestValidateMutualTLSSecurityScheme(t *testing.T) {
	// mutualTLS is valid in OpenAPI 3.1
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info:    &Info{Title: "Test", Version: "1.0"},
		Paths:   &Paths{Paths: map[string]*PathItem{"/": {}}},
		Components: &Components{
			SecuritySchemes: map[string]*SecurityScheme{
				"mtls": {Type: "mutualTLS"},
			},
		},
	}

	result := api.Validate()
	// Should not have errors about invalid type
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "mutualTLS") {
			t.Errorf("mutualTLS should be valid in OpenAPI 3.1, got error: %v", err)
		}
	}
}

func TestValidateSchemaTypeArray(t *testing.T) {
	// OpenAPI 3.1 allows type arrays like ["string", "null"]
	jsonDoc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0"},
		"paths": {
			"/test": {
				"get": {
					"responses": {
						"200": {
							"description": "OK",
							"content": {
								"application/json": {
									"schema": {
										"type": ["string", "null"]
									}
								}
							}
						}
					}
				}
			}
		}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(jsonDoc), &api); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	result := api.Validate()
	if !result.Valid() {
		t.Errorf("Expected valid document with type array, got: %v", result.Error())
	}
}

func TestValidateSchemaInvalidType(t *testing.T) {
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info:    &Info{Title: "Test", Version: "1.0"},
		Paths:   &Paths{Paths: map[string]*PathItem{"/": {}}},
		Components: &Components{
			Schemas: map[string]*Schema{
				"Invalid": {Type: &StringOrStringArray{String: "invalid"}},
			},
		},
	}

	result := api.Validate()
	found := false
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "invalid type") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for invalid type")
	}
}

func TestValidateSchemaConstraints(t *testing.T) {
	tests := []struct {
		name        string
		schema      *Schema
		expectError string
	}{
		{
			name:        "min > max",
			schema:      &Schema{Type: &StringOrStringArray{String: "number"}, Minimum: float64Ptr(10), Maximum: float64Ptr(5)},
			expectError: "minimum cannot be greater than maximum",
		},
		{
			name:        "minLength > maxLength",
			schema:      &Schema{Type: &StringOrStringArray{String: "string"}, MinLength: intPtr(10), MaxLength: intPtr(5)},
			expectError: "minLength cannot be greater than maxLength",
		},
		{
			name:        "invalid pattern",
			schema:      &Schema{Type: &StringOrStringArray{String: "string"}, Pattern: "[invalid(regex"},
			expectError: "invalid regex pattern",
		},
		{
			name: "required property not in properties",
			schema: &Schema{
				Type:       &StringOrStringArray{String: "object"},
				Required:   []string{"missing"},
				Properties: map[string]*Schema{"existing": {Type: &StringOrStringArray{String: "string"}}},
			},
			expectError: "required property 'missing' not defined",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &OpenAPI{
				OpenAPI: "3.1.0",
				Info:    &Info{Title: "Test", Version: "1.0"},
				Paths:   &Paths{Paths: map[string]*PathItem{"/": {}}},
				Components: &Components{
					Schemas: map[string]*Schema{"test": tc.schema},
				},
			}

			result := api.Validate()
			found := false
			for _, err := range result.Errors {
				if strings.Contains(err.Message, tc.expectError) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error containing '%s', got: %v", tc.expectError, result.Error())
			}
		})
	}
}

func TestValidateServerVariable(t *testing.T) {
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info:    &Info{Title: "Test", Version: "1.0"},
		Paths:   &Paths{Paths: map[string]*PathItem{"/": {}}},
		Servers: []*Server{
			{
				URL: "https://{env}.example.com",
				Variables: map[string]*ServerVariable{
					"env": {
						Enum:    []string{"dev", "staging", "prod"},
						Default: "invalid", // Not in enum
					},
				},
			},
		},
	}

	result := api.Validate()
	found := false
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "default value must be one of the enum values") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for server variable default not in enum")
	}
}

func TestValidateWebhooks(t *testing.T) {
	api := &OpenAPI{
		OpenAPI: "3.1.0",
		Info:    &Info{Title: "Test", Version: "1.0"},
		Webhooks: map[string]*PathItem{
			"newPet": {
				Post: &Operation{
					// Missing responses
				},
			},
		},
	}

	result := api.Validate()
	found := false
	for _, err := range result.Errors {
		// Error path should be like "webhooks[newPet].post.responses"
		if strings.Contains(err.Path, "webhooks") && strings.Contains(err.Path, "responses") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected error for webhook operation without responses, got: %v", result.Error())
	}
}

// Helper functions
func float64Ptr(v float64) *float64 { return &v }
func intPtr(v int) *int             { return &v }
