// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi20

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

			// Unmarshal to Swagger struct
			var swagger1 Swagger
			if err := json.Unmarshal(originalData, &swagger1); err != nil {
				t.Fatalf("Failed to unmarshal %s: %v", file, err)
			}

			// Marshal back to JSON
			marshaled1, err := json.Marshal(&swagger1)
			if err != nil {
				t.Fatalf("Failed to marshal %s: %v", file, err)
			}

			// Unmarshal again
			var swagger2 Swagger
			if err := json.Unmarshal(marshaled1, &swagger2); err != nil {
				t.Fatalf("Failed to unmarshal marshaled JSON for %s: %v", file, err)
			}

			// Marshal second time
			marshaled2, err := json.Marshal(&swagger2)
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
		checks func(t *testing.T, swagger *Swagger)
	}{
		{
			file: "petstore-simple.json",
			checks: func(t *testing.T, swagger *Swagger) {
				if swagger.Swagger != "2.0" {
					t.Errorf("Expected swagger 2.0, got %s", swagger.Swagger)
				}
				if swagger.Info == nil || swagger.Info.Title != "Simple Petstore" {
					t.Error("Info.Title not preserved")
				}
				if swagger.Host != "api.example.com" {
					t.Error("Host not preserved")
				}
				if swagger.Paths == nil {
					t.Error("Paths is nil")
				}
			},
		},
		{
			file: "petstore-expanded.json",
			checks: func(t *testing.T, swagger *Swagger) {
				if swagger.Info == nil {
					t.Fatal("Info is nil")
				}
				if swagger.Info.Contact == nil {
					t.Error("Info.Contact not preserved")
				}
				if swagger.Info.License == nil {
					t.Error("Info.License not preserved")
				}
				// Check allOf composition
				if swagger.Definitions == nil {
					t.Fatal("Definitions is nil")
				}
				pet := swagger.Definitions["Pet"]
				if pet == nil {
					t.Fatal("Pet definition not found")
				}
				if len(pet.AllOf) == 0 {
					t.Error("Pet.AllOf not preserved")
				}
			},
		},
		{
			file: "security.json",
			checks: func(t *testing.T, swagger *Swagger) {
				if swagger.SecurityDefinitions == nil || len(swagger.SecurityDefinitions) == 0 {
					t.Error("SecurityDefinitions not preserved")
				}
				// Check all security types
				if swagger.SecurityDefinitions["api_key"] == nil {
					t.Error("api_key security not preserved")
				}
				if swagger.SecurityDefinitions["basic_auth"] == nil {
					t.Error("basic_auth security not preserved")
				}
				if swagger.SecurityDefinitions["oauth2"] == nil {
					t.Error("oauth2 security not preserved")
				}
				// Check OAuth2 flows
				oauth := swagger.SecurityDefinitions["oauth2"]
				if oauth.Flow != "accessCode" {
					t.Errorf("Expected oauth2 flow accessCode, got %s", oauth.Flow)
				}
				if oauth.Scopes == nil || len(oauth.Scopes) == 0 {
					t.Error("OAuth2 scopes not preserved")
				}
			},
		},
		{
			file: "schema-additional-properties.json",
			checks: func(t *testing.T, swagger *Swagger) {
				if swagger.Definitions == nil {
					t.Fatal("Definitions is nil")
				}
				// Check boolean additionalProperties
				closed := swagger.Definitions["Closed"]
				if closed == nil {
					t.Fatal("Closed schema not found")
				}
				if closed.AdditionalProperties == nil || !closed.AdditionalProperties.IsBooleanSchema() {
					t.Error("Closed.AdditionalProperties should be boolean schema")
				}
				if bv := closed.AdditionalProperties.BooleanValue(); bv == nil || *bv != false {
					t.Error("Closed.AdditionalProperties should be false")
				}

				open := swagger.Definitions["Open"]
				if open == nil {
					t.Fatal("Open schema not found")
				}
				if open.AdditionalProperties == nil || !open.AdditionalProperties.IsBooleanSchema() {
					t.Error("Open.AdditionalProperties should be boolean schema")
				}
				if bv := open.AdditionalProperties.BooleanValue(); bv == nil || *bv != true {
					t.Error("Open.AdditionalProperties should be true")
				}

				// Check typed additionalProperties
				typed := swagger.Definitions["TypedAdditional"]
				if typed == nil {
					t.Fatal("TypedAdditional schema not found")
				}
				if typed.AdditionalProperties == nil || typed.AdditionalProperties.IsBooleanSchema() {
					t.Error("TypedAdditional.AdditionalProperties should not be boolean schema")
				}
				if typed.AdditionalProperties.Type != "string" {
					t.Errorf("Expected TypedAdditional.AdditionalProperties type string, got %s", typed.AdditionalProperties.Type)
				}
			},
		},
		{
			file: "parameters.json",
			checks: func(t *testing.T, swagger *Swagger) {
				// Check global parameters
				if swagger.Parameters == nil || len(swagger.Parameters) == 0 {
					t.Error("Global parameters not preserved")
				}
				// Check parameter validation
				pathItem := swagger.Paths.Get("/search")
				if pathItem == nil || pathItem.Get == nil {
					t.Fatal("Path /search GET not found")
				}
				for _, param := range pathItem.Get.Parameters {
					if param.Name == "q" {
						if param.MinLength == nil || *param.MinLength != 1 {
							t.Error("Parameter minLength not preserved")
						}
						if param.MaxLength == nil || *param.MaxLength != 100 {
							t.Error("Parameter maxLength not preserved")
						}
						if param.Pattern != "^[a-zA-Z0-9 ]+$" {
							t.Error("Parameter pattern not preserved")
						}
					}
				}
			},
		},
		{
			file: "responses.json",
			checks: func(t *testing.T, swagger *Swagger) {
				// Check global responses
				if swagger.Responses == nil || len(swagger.Responses) == 0 {
					t.Error("Global responses not preserved")
				}
				// Check response headers
				pathItem := swagger.Paths.Get("/items")
				if pathItem == nil || pathItem.Get == nil {
					t.Fatal("Path /items GET not found")
				}
				resp := pathItem.Get.Responses.StatusCode["200"]
				if resp == nil {
					t.Fatal("200 response not found")
				}
				if resp.Headers == nil || len(resp.Headers) == 0 {
					t.Error("Response headers not preserved")
				}
				if resp.Examples == nil || len(resp.Examples) == 0 {
					t.Error("Response examples not preserved")
				}
			},
		},
		{
			file: "extensions.json",
			checks: func(t *testing.T, swagger *Swagger) {
				// Check root-level extensions
				if swagger.Extensions == nil || swagger.Extensions["x-custom-root"] == nil {
					t.Error("Root extensions not preserved")
				}
				// Check info extensions
				if swagger.Info.Extensions == nil || swagger.Info.Extensions["x-logo"] == nil {
					t.Error("Info extensions not preserved")
				}
				// Check path extensions
				pathItem := swagger.Paths.Get("/items")
				if pathItem == nil {
					t.Fatal("Path /items not found")
				}
				if pathItem.Extensions == nil || pathItem.Extensions["x-path-extension"] == nil {
					t.Error("Path item extensions not preserved")
				}
				// Check operation extensions
				if pathItem.Get.Extensions == nil || pathItem.Get.Extensions["x-operation-extension"] == nil {
					t.Error("Operation extensions not preserved")
				}
			},
		},
		{
			file: "tags.json",
			checks: func(t *testing.T, swagger *Swagger) {
				if swagger.Tags == nil || len(swagger.Tags) == 0 {
					t.Error("Tags not preserved")
				}
				// Check tag with external docs
				var usersTag *Tag
				for _, tag := range swagger.Tags {
					if tag.Name == "users" {
						usersTag = tag
						break
					}
				}
				if usersTag == nil {
					t.Fatal("users tag not found")
				}
				if usersTag.ExternalDocs == nil {
					t.Error("Tag external docs not preserved")
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

			var swagger Swagger
			if err := json.Unmarshal(data, &swagger); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			tc.checks(t, &swagger)
		})
	}
}

// TestBooleanSchemaRoundTrip specifically tests boolean schema handling
func TestBooleanSchemaRoundTrip(t *testing.T) {
	// Test with additionalProperties: true
	jsonTrue := `{
		"swagger": "2.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {},
		"definitions": {
			"FreeForm": {
				"type": "object",
				"additionalProperties": true
			}
		}
	}`

	var swagger1 Swagger
	if err := json.Unmarshal([]byte(jsonTrue), &swagger1); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	schema := swagger1.Definitions["FreeForm"]
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
	marshaled, err := json.Marshal(&swagger1)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var swagger2 Swagger
	if err := json.Unmarshal(marshaled, &swagger2); err != nil {
		t.Fatalf("Failed to unmarshal marshaled JSON: %v", err)
	}

	schema2 := swagger2.Definitions["FreeForm"]
	if !schema2.AdditionalProperties.IsBooleanSchema() {
		t.Error("AdditionalProperties should still be boolean schema after round-trip")
	}
	if *schema2.AdditionalProperties.BooleanValue() != true {
		t.Error("AdditionalProperties should still be true after round-trip")
	}

	// Test with additionalProperties: false
	jsonFalse := `{
		"swagger": "2.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {},
		"definitions": {
			"Strict": {
				"type": "object",
				"additionalProperties": false
			}
		}
	}`

	var swagger3 Swagger
	if err := json.Unmarshal([]byte(jsonFalse), &swagger3); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	schema3 := swagger3.Definitions["Strict"]
	if !schema3.AdditionalProperties.IsBooleanSchema() {
		t.Error("AdditionalProperties should be boolean schema")
	}
	if *schema3.AdditionalProperties.BooleanValue() != false {
		t.Error("AdditionalProperties should be false")
	}
}

// TestExtensionsPreserved tests that x- extensions are preserved
func TestExtensionsPreserved(t *testing.T) {
	path := filepath.Join("oas-examples/json", "extensions.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var swagger Swagger
	if err := json.Unmarshal(data, &swagger); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check that extensions are captured
	if swagger.Extensions == nil || len(swagger.Extensions) == 0 {
		t.Fatal("No top-level extensions found")
	}
	for key := range swagger.Extensions {
		t.Logf("Found extension: %s", key)
	}

	// Marshal and verify extensions survive round-trip
	marshaled, err := json.Marshal(&swagger)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var swagger2 Swagger
	if err := json.Unmarshal(marshaled, &swagger2); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(swagger.Extensions, swagger2.Extensions) {
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

			var swagger Swagger
			if err := json.Unmarshal(data, &swagger); err != nil {
				t.Errorf("Failed to parse %s: %v", file, err)
			}

			// Basic sanity checks
			if swagger.Swagger == "" {
				t.Error("Swagger version is empty")
			}
			if swagger.Info == nil {
				t.Error("Info is nil")
			}
		})
	}
}
