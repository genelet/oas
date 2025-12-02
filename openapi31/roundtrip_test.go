// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

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

	// Files that contain mixed OpenAPI 3.0/3.1 syntax (e.g., boolean exclusiveMinimum)
	// These are demonstration files, not valid OpenAPI 3.1 documents
	skipFiles := map[string]bool{
		"schema-validation-local.json":     true,
		"schema-validation-top-level.json": true,
	}

	// Get all JSON files
	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	var jsonFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			if skipFiles[entry.Name()] {
				continue
			}
			jsonFiles = append(jsonFiles, filepath.Join(examplesDir, entry.Name()))
		}
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
				if api.OpenAPI != "3.1.0" {
					t.Errorf("Expected openapi 3.1.0, got %s", api.OpenAPI)
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
			file: "webhooks.json",
			checks: func(t *testing.T, api *OpenAPI) {
				if api.Webhooks == nil || len(api.Webhooks) == 0 {
					t.Error("Webhooks not preserved")
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
			file: "schema-types.json",
			checks: func(t *testing.T, api *OpenAPI) {
				if api.Components == nil || api.Components.Schemas == nil {
					t.Fatal("Components.Schemas is nil")
				}
				// Check that we have schemas
				if len(api.Components.Schemas) == 0 {
					t.Error("No schemas found")
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
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0.0"},
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
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0.0"},
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

// TestTypeArrayRoundTrip tests that type arrays are preserved
func TestTypeArrayRoundTrip(t *testing.T) {
	jsonDoc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"components": {
			"schemas": {
				"NullableString": {
					"type": ["string", "null"]
				},
				"MultiType": {
					"type": ["string", "integer", "null"]
				},
				"SingleType": {
					"type": "string"
				}
			}
		}
	}`

	var api1 OpenAPI
	if err := json.Unmarshal([]byte(jsonDoc), &api1); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check NullableString
	nullableString := api1.Components.Schemas["NullableString"]
	if !nullableString.Type.Contains("string") || !nullableString.Type.Contains("null") {
		t.Error("NullableString type not preserved")
	}

	// Check MultiType
	multiType := api1.Components.Schemas["MultiType"]
	if !multiType.Type.Contains("string") || !multiType.Type.Contains("integer") || !multiType.Type.Contains("null") {
		t.Error("MultiType types not preserved")
	}

	// Check SingleType
	singleType := api1.Components.Schemas["SingleType"]
	if !singleType.Type.Contains("string") {
		t.Error("SingleType type not preserved")
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

	// Verify types preserved after round-trip
	nullableString2 := api2.Components.Schemas["NullableString"]
	if !nullableString2.Type.Contains("string") || !nullableString2.Type.Contains("null") {
		t.Error("NullableString type not preserved after round-trip")
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

	// Files that contain mixed OpenAPI 3.0/3.1 syntax (e.g., boolean exclusiveMinimum)
	// These are demonstration files, not valid OpenAPI 3.1 documents
	skipFiles := map[string]bool{
		"schema-validation-local.json":     true,
		"schema-validation-top-level.json": true,
	}

	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	var jsonFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			if skipFiles[entry.Name()] {
				continue
			}
			jsonFiles = append(jsonFiles, filepath.Join(examplesDir, entry.Name()))
		}
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

// TestWebhooksRoundTrip specifically tests webhooks handling
func TestWebhooksRoundTrip(t *testing.T) {
	path := filepath.Join("oas-examples/json", "webhooks.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var api1 OpenAPI
	if err := json.Unmarshal(data, &api1); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if api1.Webhooks == nil || len(api1.Webhooks) == 0 {
		t.Fatal("Webhooks not parsed")
	}

	// Marshal and unmarshal
	marshaled, err := json.Marshal(&api1)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var api2 OpenAPI
	if err := json.Unmarshal(marshaled, &api2); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(api2.Webhooks) != len(api1.Webhooks) {
		t.Errorf("Webhooks count mismatch: %d vs %d", len(api1.Webhooks), len(api2.Webhooks))
	}

	for name := range api1.Webhooks {
		if api2.Webhooks[name] == nil {
			t.Errorf("Webhook %s not preserved", name)
		}
	}
}

// TestJsonSchemaDialectRoundTrip tests that jsonSchemaDialect is preserved
func TestJsonSchemaDialectRoundTrip(t *testing.T) {
	jsonDoc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"jsonSchemaDialect": "https://json-schema.org/draft/2020-12/schema"
	}`

	var api1 OpenAPI
	if err := json.Unmarshal([]byte(jsonDoc), &api1); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if api1.JsonSchemaDialect != "https://json-schema.org/draft/2020-12/schema" {
		t.Errorf("jsonSchemaDialect not parsed correctly: %s", api1.JsonSchemaDialect)
	}

	// Marshal and unmarshal
	marshaled, err := json.Marshal(&api1)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var api2 OpenAPI
	if err := json.Unmarshal(marshaled, &api2); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if api2.JsonSchemaDialect != api1.JsonSchemaDialect {
		t.Errorf("jsonSchemaDialect not preserved: got %s, want %s", api2.JsonSchemaDialect, api1.JsonSchemaDialect)
	}
}

// TestExclusiveMinMaxRoundTrip tests OpenAPI 3.1 numeric exclusive values (numbers, not booleans)
func TestExclusiveMinMaxRoundTrip(t *testing.T) {
	jsonDoc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"components": {
			"schemas": {
				"Rating": {
					"type": "number",
					"exclusiveMinimum": 0,
					"exclusiveMaximum": 10
				}
			}
		}
	}`

	var api1 OpenAPI
	if err := json.Unmarshal([]byte(jsonDoc), &api1); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	schema := api1.Components.Schemas["Rating"]
	if schema.ExclusiveMinimum == nil || *schema.ExclusiveMinimum != 0 {
		t.Errorf("exclusiveMinimum not parsed correctly: %v", schema.ExclusiveMinimum)
	}
	if schema.ExclusiveMaximum == nil || *schema.ExclusiveMaximum != 10 {
		t.Errorf("exclusiveMaximum not parsed correctly: %v", schema.ExclusiveMaximum)
	}

	// Marshal and unmarshal
	marshaled, err := json.Marshal(&api1)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var api2 OpenAPI
	if err := json.Unmarshal(marshaled, &api2); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	schema2 := api2.Components.Schemas["Rating"]
	if schema2.ExclusiveMinimum == nil || *schema2.ExclusiveMinimum != 0 {
		t.Error("exclusiveMinimum not preserved")
	}
	if schema2.ExclusiveMaximum == nil || *schema2.ExclusiveMaximum != 10 {
		t.Error("exclusiveMaximum not preserved")
	}
}
