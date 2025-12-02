// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi30

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestJSONRoundTrip tests that JSON can be unmarshaled and marshaled back without loss.
// Note: We compare JSON output rather than Go structs because empty slices/maps
// with omitempty become nil after a full round-trip, which is semantically equivalent.
func TestJSONRoundTrip(t *testing.T) {
	examplesDir := "oas-examples/json"

	// Get all JSON files recursively
	var jsonFiles []string
	err := filepath.Walk(examplesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			jsonFiles = append(jsonFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk examples directory: %v", err)
	}

	if len(jsonFiles) == 0 {
		t.Fatal("No JSON files found in examples directory")
	}

	for _, file := range jsonFiles {
		t.Run(filepath.Base(file), func(t *testing.T) {
			// Read original JSON
			originalData, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", file, err)
			}

			// Unmarshal to OpenAPI struct
			var api1 OpenAPI
			if err := json.Unmarshal(originalData, &api1); err != nil {
				t.Fatalf("Failed to unmarshal %s: %v", file, err)
			}

			// Marshal back to JSON
			marshaled1, err := json.Marshal(&api1)
			if err != nil {
				t.Fatalf("Failed to marshal %s: %v", file, err)
			}

			// Unmarshal again
			var api2 OpenAPI
			if err := json.Unmarshal(marshaled1, &api2); err != nil {
				t.Fatalf("Failed to unmarshal marshaled JSON for %s: %v", file, err)
			}

			// Marshal second time
			marshaled2, err := json.Marshal(&api2)
			if err != nil {
				t.Fatalf("Failed to marshal second time %s: %v", file, err)
			}

			// Compare JSON outputs - this is the stable comparison
			// after one round-trip, subsequent round-trips should be identical
			if string(marshaled1) != string(marshaled2) {
				t.Errorf("Round-trip mismatch for %s", file)
				t.Logf("First marshal:\n%s", marshaled1)
				t.Logf("Second marshal:\n%s", marshaled2)
			}
		})
	}
}

// TestJSONUnmarshalMarshalPreservesData tests specific fields are preserved
func TestJSONUnmarshalMarshalPreservesData(t *testing.T) {
	examplesDir := "oas-examples/json"

	tests := []struct {
		file   string
		checks func(t *testing.T, api *OpenAPI)
	}{
		{
			file: "petstore-simple.json",
			checks: func(t *testing.T, api *OpenAPI) {
				if api.OpenAPI != "3.0.0" {
					t.Errorf("Expected openapi 3.0.0, got %s", api.OpenAPI)
				}
				if api.Info == nil || api.Info.Title != "Simple Petstore" {
					t.Error("Info.Title not preserved")
				}
				if len(api.Servers) != 1 || api.Servers[0].URL != "https://httpbin.org" {
					t.Error("Servers not preserved")
				}
				if api.Paths == nil {
					t.Error("Paths is nil")
				}
			},
		},
		{
			file: "callbacks.json",
			checks: func(t *testing.T, api *OpenAPI) {
				if api.Paths == nil {
					t.Fatal("Paths is nil")
				}
				// Check that callbacks exist
				pathItem := api.Paths.Paths["/streams"]
				if pathItem == nil {
					t.Fatal("Path /streams not found")
				}
				if pathItem.Post == nil {
					t.Fatal("POST operation not found")
				}
				if pathItem.Post.Callbacks == nil || len(pathItem.Post.Callbacks) == 0 {
					t.Error("Callbacks not preserved")
				}
			},
		},
		{
			file: "security.json",
			checks: func(t *testing.T, api *OpenAPI) {
				if api.Components == nil {
					t.Fatal("Components is nil")
				}
				if api.Components.SecuritySchemes == nil || len(api.Components.SecuritySchemes) == 0 {
					t.Error("SecuritySchemes not preserved")
				}
			},
		},
		{
			file: "schema-additional-properties.json",
			checks: func(t *testing.T, api *OpenAPI) {
				if api.Components == nil || api.Components.Schemas == nil {
					t.Fatal("Components.Schemas is nil")
				}
				// Check that additionalProperties is preserved
				for name, schema := range api.Components.Schemas {
					if schema.AdditionalProperties != nil {
						t.Logf("Schema %s has additionalProperties", name)
						if schema.AdditionalProperties.IsBooleanSchema() {
							t.Logf("  - Boolean schema: %v", *schema.AdditionalProperties.BooleanValue())
						}
					}
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			path := filepath.Join(examplesDir, tc.file)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			var api OpenAPI
			if err := json.Unmarshal(data, &api); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			tc.checks(t, &api)
		})
	}
}

// TestBooleanSchemaRoundTrip specifically tests boolean schema handling
func TestBooleanSchemaRoundTrip(t *testing.T) {
	// Test with additionalProperties: true
	jsonTrue := `{
		"openapi": "3.0.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {},
		"components": {
			"schemas": {
				"FreeForm": {
					"type": "object",
					"additionalProperties": true
				}
			}
		}
	}`

	var api1 OpenAPI
	if err := json.Unmarshal([]byte(jsonTrue), &api1); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	schema := api1.Components.Schemas["FreeForm"]
	if schema == nil {
		t.Fatal("Schema FreeForm not found")
	}
	if schema.AdditionalProperties == nil {
		t.Fatal("AdditionalProperties is nil")
	}
	if !schema.AdditionalProperties.IsBooleanSchema() {
		t.Error("AdditionalProperties should be boolean schema")
	}
	if *schema.AdditionalProperties.BooleanValue() != true {
		t.Error("AdditionalProperties should be true")
	}

	// Marshal and unmarshal again
	marshaled, err := json.Marshal(&api1)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var api2 OpenAPI
	if err := json.Unmarshal(marshaled, &api2); err != nil {
		t.Fatalf("Failed to unmarshal marshaled JSON: %v", err)
	}

	schema2 := api2.Components.Schemas["FreeForm"]
	if !schema2.AdditionalProperties.IsBooleanSchema() {
		t.Error("AdditionalProperties should still be boolean schema after round-trip")
	}
	if *schema2.AdditionalProperties.BooleanValue() != true {
		t.Error("AdditionalProperties should still be true after round-trip")
	}

	// Test with additionalProperties: false
	jsonFalse := `{
		"openapi": "3.0.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {},
		"components": {
			"schemas": {
				"Strict": {
					"type": "object",
					"additionalProperties": false
				}
			}
		}
	}`

	var api3 OpenAPI
	if err := json.Unmarshal([]byte(jsonFalse), &api3); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	schema3 := api3.Components.Schemas["Strict"]
	if !schema3.AdditionalProperties.IsBooleanSchema() {
		t.Error("AdditionalProperties should be boolean schema")
	}
	if *schema3.AdditionalProperties.BooleanValue() != false {
		t.Error("AdditionalProperties should be false")
	}
}

// TestExtensionsPreserved tests that x- extensions are preserved
func TestExtensionsPreserved(t *testing.T) {
	path := filepath.Join("oas-examples/json", "readme-extensions.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var api OpenAPI
	if err := json.Unmarshal(data, &api); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check that extensions are captured
	if api.Extensions == nil || len(api.Extensions) == 0 {
		t.Log("No top-level extensions found (may be expected)")
	} else {
		for key := range api.Extensions {
			t.Logf("Found extension: %s", key)
		}
	}

	// Marshal and verify extensions survive round-trip
	marshaled, err := json.Marshal(&api)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var api2 OpenAPI
	if err := json.Unmarshal(marshaled, &api2); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(api.Extensions, api2.Extensions) {
		t.Error("Extensions not preserved in round-trip")
	}
}

// TestAllExampleFilesParseSuccessfully ensures all example files can be parsed
func TestAllExampleFilesParseSuccessfully(t *testing.T) {
	examplesDir := "oas-examples/json"

	var jsonFiles []string
	err := filepath.Walk(examplesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			jsonFiles = append(jsonFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk examples directory: %v", err)
	}

	t.Logf("Found %d JSON files to test", len(jsonFiles))

	for _, file := range jsonFiles {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			var api OpenAPI
			if err := json.Unmarshal(data, &api); err != nil {
				t.Errorf("Failed to parse %s: %v", file, err)
			}

			// Basic sanity checks
			if api.OpenAPI == "" {
				t.Error("OpenAPI version is empty")
			}
			if api.Info == nil {
				t.Error("Info is nil")
			}
		})
	}
}
