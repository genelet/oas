// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import (
	"encoding/json"
	"testing"
)

func TestParseSimpleOpenAPI(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {
			"title": "Test API",
			"version": "1.0.0"
		},
		"paths": {}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if api.OpenAPI != "3.1.0" {
		t.Errorf("Expected openapi '3.1.0', got %s", api.OpenAPI)
	}
	if api.Info.Title != "Test API" {
		t.Errorf("Expected title 'Test API', got %s", api.Info.Title)
	}
	if api.Info.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", api.Info.Version)
	}
}

func TestParseWithSchema(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0"},
		"paths": {},
		"components": {
			"schemas": {
				"Pet": {
					"type": "object",
					"properties": {
						"name": {"type": "string"},
						"age": {"type": "integer"}
					},
					"required": ["name"]
				}
			}
		}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	petSchema := api.Components.Schemas["Pet"]
	if petSchema == nil {
		t.Fatal("Expected Pet schema")
	}

	if !petSchema.Type.Contains("object") {
		t.Errorf("Expected type 'object', got %v", petSchema.Type)
	}

	if len(petSchema.Required) != 1 || petSchema.Required[0] != "name" {
		t.Errorf("Expected required ['name'], got %v", petSchema.Required)
	}
}

func TestBooleanSchema(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0"},
		"paths": {},
		"components": {
			"schemas": {
				"Any": true,
				"None": false
			}
		}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	anySchema := api.Components.Schemas["Any"]
	if !anySchema.IsBooleanSchema() || *anySchema.BooleanValue() != true {
		t.Error("Expected boolean true schema")
	}

	noneSchema := api.Components.Schemas["None"]
	if !noneSchema.IsBooleanSchema() || *noneSchema.BooleanValue() != false {
		t.Error("Expected boolean false schema")
	}
}

func TestTypeArray(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0"},
		"paths": {},
		"components": {
			"schemas": {
				"NullableString": {
					"type": ["string", "null"]
				}
			}
		}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	schema := api.Components.Schemas["NullableString"]
	if !schema.Type.Contains("string") {
		t.Error("Expected type to contain 'string'")
	}
	if !schema.Type.Contains("null") {
		t.Error("Expected type to contain 'null'")
	}
}

func TestExtensions(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {
			"title": "Test",
			"version": "1.0",
			"x-custom": "value"
		},
		"paths": {},
		"x-root-extension": true
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if api.Extensions["x-root-extension"] != true {
		t.Errorf("Expected x-root-extension true, got %v", api.Extensions["x-root-extension"])
	}

	if api.Info.Extensions["x-custom"] != "value" {
		t.Errorf("Expected x-custom 'value', got %v", api.Info.Extensions["x-custom"])
	}
}

func TestWebhooks(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0"},
		"webhooks": {
			"newPet": {
				"post": {
					"summary": "New pet webhook",
					"responses": {
						"200": {"description": "OK"}
					}
				}
			}
		}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if api.Webhooks == nil || api.Webhooks["newPet"] == nil {
		t.Fatal("Expected newPet webhook")
	}

	if api.Webhooks["newPet"].Post == nil {
		t.Error("Expected post operation")
	}

	if api.Webhooks["newPet"].Post.Summary != "New pet webhook" {
		t.Errorf("Expected summary 'New pet webhook', got %s", api.Webhooks["newPet"].Post.Summary)
	}
}

func TestRoundTrip(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test API", "version": "1.0.0"},
		"jsonSchemaDialect": "https://json-schema.org/draft/2020-12/schema",
		"paths": {
			"/pets": {
				"get": {
					"summary": "List pets",
					"responses": {
						"200": {"description": "OK"}
					}
				}
			}
		},
		"components": {
			"schemas": {
				"Pet": {
					"type": "object",
					"properties": {
						"name": {"type": "string"}
					}
				}
			}
		}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Marshal back
	out, err := json.Marshal(&api)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Parse again
	var api2 OpenAPI
	if err := json.Unmarshal(out, &api2); err != nil {
		t.Fatalf("Failed to re-parse: %v", err)
	}

	// Verify
	if api2.OpenAPI != "3.1.0" {
		t.Error("Round-trip failed: openapi version mismatch")
	}
	if api2.JsonSchemaDialect != "https://json-schema.org/draft/2020-12/schema" {
		t.Error("Round-trip failed: jsonSchemaDialect mismatch")
	}
	if api2.Paths.Get("/pets") == nil {
		t.Error("Round-trip failed: path not found")
	}
	if api2.Components.Schemas["Pet"] == nil {
		t.Error("Round-trip failed: schema not found")
	}
}

func TestReference(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0"},
		"paths": {
			"/pets": {
				"get": {
					"parameters": [
						{"$ref": "#/components/parameters/limit"},
						{"name": "offset", "in": "query", "schema": {"type": "integer"}}
					],
					"responses": {
						"200": {"description": "OK"}
					}
				}
			}
		}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	params := api.Paths.Get("/pets").Get.Parameters
	if len(params) != 2 {
		t.Fatalf("Expected 2 parameters, got %d", len(params))
	}

	// First is reference
	if !params[0].IsReference() {
		t.Error("Expected first parameter to be a reference")
	}
	if params[0].Ref != "#/components/parameters/limit" {
		t.Errorf("Expected ref '#/components/parameters/limit', got %s", params[0].Ref)
	}

	// Second is inline
	if params[1].IsReference() {
		t.Error("Expected second parameter to NOT be a reference")
	}
	if params[1].Name != "offset" {
		t.Errorf("Expected name 'offset', got %s", params[1].Name)
	}
}
