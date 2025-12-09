// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi20

import (
	"encoding/json"
	"testing"
)

// Sample Swagger 2.0 document for testing
var petstoreSwagger = `{
  "swagger": "2.0",
  "info": {
    "title": "Petstore",
    "version": "1.0.0",
    "description": "A sample Petstore API",
    "contact": {
      "name": "API Support",
      "email": "support@example.com"
    },
    "license": {
      "name": "MIT"
    }
  },
  "host": "api.example.com",
  "basePath": "/v1",
  "schemes": ["https"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "paths": {
    "/pets": {
      "get": {
        "operationId": "listPets",
        "summary": "List all pets",
        "tags": ["pets"],
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "type": "integer",
            "format": "int32"
          }
        ],
        "responses": {
          "200": {
            "description": "A list of pets",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Pet"
              }
            }
          },
          "default": {
            "description": "Error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "operationId": "createPet",
        "summary": "Create a pet",
        "tags": ["pets"],
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Pet"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Pet created"
          }
        }
      }
    },
    "/pets/{petId}": {
      "get": {
        "operationId": "getPet",
        "summary": "Get a pet by ID",
        "parameters": [
          {
            "name": "petId",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "A pet",
            "schema": {
              "$ref": "#/definitions/Pet"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Pet": {
      "type": "object",
      "required": ["name"],
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
        },
        "tag": {
          "type": "string"
        }
      }
    },
    "Error": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer"
        },
        "message": {
          "type": "string"
        }
      }
    }
  },
  "securityDefinitions": {
    "api_key": {
      "type": "apiKey",
      "name": "X-API-Key",
      "in": "header"
    },
    "oauth2": {
      "type": "oauth2",
      "flow": "implicit",
      "authorizationUrl": "https://example.com/oauth/authorize",
      "scopes": {
        "read:pets": "Read pets",
        "write:pets": "Write pets"
      }
    }
  },
  "tags": [
    {
      "name": "pets",
      "description": "Pet operations"
    }
  ]
}`

func TestUnmarshalSwagger(t *testing.T) {
	var swagger Swagger
	if err := json.Unmarshal([]byte(petstoreSwagger), &swagger); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check basic fields
	if swagger.Swagger != "2.0" {
		t.Errorf("Expected swagger 2.0, got %s", swagger.Swagger)
	}
	if swagger.Info == nil {
		t.Fatal("Info is nil")
	}
	if swagger.Info.Title != "Petstore" {
		t.Errorf("Expected title Petstore, got %s", swagger.Info.Title)
	}
	if swagger.Host != "api.example.com" {
		t.Errorf("Expected host api.example.com, got %s", swagger.Host)
	}
	if swagger.BasePath != "/v1" {
		t.Errorf("Expected basePath /v1, got %s", swagger.BasePath)
	}

	// Check paths
	if swagger.Paths == nil {
		t.Fatal("Paths is nil")
	}
	petsPath := swagger.Paths.Get("/pets")
	if petsPath == nil {
		t.Fatal("Path /pets not found")
	}
	if petsPath.Get == nil {
		t.Fatal("GET operation on /pets not found")
	}
	if petsPath.Get.OperationID != "listPets" {
		t.Errorf("Expected operationId listPets, got %s", petsPath.Get.OperationID)
	}

	// Check parameters
	if len(petsPath.Get.Parameters) != 1 {
		t.Fatalf("Expected 1 parameter, got %d", len(petsPath.Get.Parameters))
	}
	param := petsPath.Get.Parameters[0]
	if param.Name != "limit" {
		t.Errorf("Expected parameter name limit, got %s", param.Name)
	}
	if param.In != "query" {
		t.Errorf("Expected parameter in query, got %s", param.In)
	}

	// Check definitions
	if swagger.Definitions == nil {
		t.Fatal("Definitions is nil")
	}
	petSchema := swagger.Definitions["Pet"]
	if petSchema == nil {
		t.Fatal("Pet definition not found")
	}
	if petSchema.Type != "object" {
		t.Errorf("Expected Pet type object, got %s", petSchema.Type)
	}
	if len(petSchema.Required) != 1 || petSchema.Required[0] != "name" {
		t.Error("Pet required fields not preserved")
	}

	// Check security definitions
	if swagger.SecurityDefinitions == nil {
		t.Fatal("SecurityDefinitions is nil")
	}
	apiKey := swagger.SecurityDefinitions["api_key"]
	if apiKey == nil {
		t.Fatal("api_key security not found")
	}
	if !apiKey.IsAPIKey() {
		t.Error("Expected api_key to be API key type")
	}
	if apiKey.In != "header" {
		t.Errorf("Expected api_key in header, got %s", apiKey.In)
	}

	oauth := swagger.SecurityDefinitions["oauth2"]
	if oauth == nil {
		t.Fatal("oauth2 security not found")
	}
	if !oauth.IsOAuth2() {
		t.Error("Expected oauth2 to be OAuth2 type")
	}
	if oauth.Flow != "implicit" {
		t.Errorf("Expected flow implicit, got %s", oauth.Flow)
	}
}

func TestRoundTrip(t *testing.T) {
	// First unmarshal
	var swagger1 Swagger
	if err := json.Unmarshal([]byte(petstoreSwagger), &swagger1); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Marshal
	marshaled1, err := json.Marshal(&swagger1)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Second unmarshal
	var swagger2 Swagger
	if err := json.Unmarshal(marshaled1, &swagger2); err != nil {
		t.Fatalf("Failed to unmarshal marshaled JSON: %v", err)
	}

	// Second marshal
	marshaled2, err := json.Marshal(&swagger2)
	if err != nil {
		t.Fatalf("Failed to marshal second time: %v", err)
	}

	// Compare - after one round-trip, subsequent should be identical
	if string(marshaled1) != string(marshaled2) {
		t.Error("Round-trip mismatch")
		t.Logf("First marshal:\n%s", marshaled1)
		t.Logf("Second marshal:\n%s", marshaled2)
	}
}

// Test boolean additionalProperties
var schemaWithBoolAdditionalProperties = `{
  "swagger": "2.0",
  "info": {
    "title": "Test",
    "version": "1.0.0"
  },
  "paths": {},
  "definitions": {
    "Closed": {
      "type": "object",
      "properties": {
        "name": {"type": "string"}
      },
      "additionalProperties": false
    },
    "Open": {
      "type": "object",
      "properties": {
        "name": {"type": "string"}
      },
      "additionalProperties": true
    },
    "Typed": {
      "type": "object",
      "properties": {
        "name": {"type": "string"}
      },
      "additionalProperties": {
        "type": "string"
      }
    }
  }
}`

func TestBooleanAdditionalProperties(t *testing.T) {
	var swagger Swagger
	if err := json.Unmarshal([]byte(schemaWithBoolAdditionalProperties), &swagger); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if swagger.Definitions == nil {
		t.Fatal("Definitions is nil")
	}

	// Test additionalProperties: false
	closed := swagger.Definitions["Closed"]
	if closed == nil {
		t.Fatal("Closed schema not found")
	}
	if closed.AdditionalProperties == nil {
		t.Fatal("Closed.AdditionalProperties is nil")
	}
	if !closed.AdditionalProperties.IsBooleanSchema() {
		t.Error("Expected Closed.AdditionalProperties to be boolean schema")
	}
	if bv := closed.AdditionalProperties.BooleanValue(); bv == nil || *bv != false {
		t.Error("Expected Closed.AdditionalProperties to be false")
	}

	// Test additionalProperties: true
	open := swagger.Definitions["Open"]
	if open == nil {
		t.Fatal("Open schema not found")
	}
	if open.AdditionalProperties == nil {
		t.Fatal("Open.AdditionalProperties is nil")
	}
	if !open.AdditionalProperties.IsBooleanSchema() {
		t.Error("Expected Open.AdditionalProperties to be boolean schema")
	}
	if bv := open.AdditionalProperties.BooleanValue(); bv == nil || *bv != true {
		t.Error("Expected Open.AdditionalProperties to be true")
	}

	// Test additionalProperties: {type: string}
	typed := swagger.Definitions["Typed"]
	if typed == nil {
		t.Fatal("Typed schema not found")
	}
	if typed.AdditionalProperties == nil {
		t.Fatal("Typed.AdditionalProperties is nil")
	}
	if typed.AdditionalProperties.IsBooleanSchema() {
		t.Error("Expected Typed.AdditionalProperties to NOT be boolean schema")
	}
	if typed.AdditionalProperties.Type != "string" {
		t.Errorf("Expected Typed.AdditionalProperties type string, got %s", typed.AdditionalProperties.Type)
	}

	// Test round-trip preserves boolean additionalProperties
	marshaled, err := json.Marshal(&swagger)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var swagger2 Swagger
	if err := json.Unmarshal(marshaled, &swagger2); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify after round-trip
	closed2 := swagger2.Definitions["Closed"]
	if closed2.AdditionalProperties == nil || !closed2.AdditionalProperties.IsBooleanSchema() {
		t.Error("Closed.AdditionalProperties not preserved as boolean after round-trip")
	}
	if bv := closed2.AdditionalProperties.BooleanValue(); bv == nil || *bv != false {
		t.Error("Closed.AdditionalProperties value not preserved after round-trip")
	}

	open2 := swagger2.Definitions["Open"]
	if open2.AdditionalProperties == nil || !open2.AdditionalProperties.IsBooleanSchema() {
		t.Error("Open.AdditionalProperties not preserved as boolean after round-trip")
	}
	if bv := open2.AdditionalProperties.BooleanValue(); bv == nil || *bv != true {
		t.Error("Open.AdditionalProperties value not preserved after round-trip")
	}
}

func TestNewBooleanSchema(t *testing.T) {
	trueSchema := NewBooleanSchema(true)
	if !trueSchema.IsBooleanSchema() {
		t.Error("Expected true schema to be boolean")
	}
	if bv := trueSchema.BooleanValue(); bv == nil || *bv != true {
		t.Error("Expected true schema value to be true")
	}

	falseSchema := NewBooleanSchema(false)
	if !falseSchema.IsBooleanSchema() {
		t.Error("Expected false schema to be boolean")
	}
	if bv := falseSchema.BooleanValue(); bv == nil || *bv != false {
		t.Error("Expected false schema value to be false")
	}

	// Test marshaling
	trueJSON, _ := json.Marshal(trueSchema)
	if string(trueJSON) != "true" {
		t.Errorf("Expected true schema to marshal as 'true', got %s", trueJSON)
	}

	falseJSON, _ := json.Marshal(falseSchema)
	if string(falseJSON) != "false" {
		t.Errorf("Expected false schema to marshal as 'false', got %s", falseJSON)
	}
}

// Test extensions
var swaggerWithExtensions = `{
  "swagger": "2.0",
  "info": {
    "title": "Test",
    "version": "1.0.0",
    "x-logo": "https://example.com/logo.png"
  },
  "x-custom": "value",
  "paths": {
    "/test": {
      "x-path-extension": "path-value",
      "get": {
        "operationId": "test",
        "x-operation-extension": "op-value",
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}`

func TestExtensions(t *testing.T) {
	var swagger Swagger
	if err := json.Unmarshal([]byte(swaggerWithExtensions), &swagger); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check root extension
	if swagger.Extensions == nil {
		t.Fatal("Root extensions is nil")
	}
	if swagger.Extensions["x-custom"] != "value" {
		t.Error("Root x-custom extension not preserved")
	}

	// Check info extension
	if swagger.Info.Extensions == nil {
		t.Fatal("Info extensions is nil")
	}
	if swagger.Info.Extensions["x-logo"] != "https://example.com/logo.png" {
		t.Error("Info x-logo extension not preserved")
	}

	// Check path item extension
	testPath := swagger.Paths.Get("/test")
	if testPath == nil {
		t.Fatal("Path /test not found")
	}
	if testPath.Extensions == nil {
		t.Fatal("PathItem extensions is nil")
	}
	if testPath.Extensions["x-path-extension"] != "path-value" {
		t.Error("PathItem x-path-extension not preserved")
	}

	// Check operation extension
	if testPath.Get.Extensions == nil {
		t.Fatal("Operation extensions is nil")
	}
	if testPath.Get.Extensions["x-operation-extension"] != "op-value" {
		t.Error("Operation x-operation-extension not preserved")
	}
}

func TestSchemaReference(t *testing.T) {
	ref := NewSchemaReference("#/definitions/Pet")
	if !ref.IsReference() {
		t.Error("Expected schema to be reference")
	}
	if ref.Ref != "#/definitions/Pet" {
		t.Errorf("Expected ref #/definitions/Pet, got %s", ref.Ref)
	}
}

func TestParameterReference(t *testing.T) {
	ref := NewParameterReference("#/parameters/limitParam")
	if !ref.IsReference() {
		t.Error("Expected parameter to be reference")
	}
	if ref.Ref != "#/parameters/limitParam" {
		t.Errorf("Expected ref #/parameters/limitParam, got %s", ref.Ref)
	}
}

func TestResponseReference(t *testing.T) {
	ref := NewResponseReference("#/responses/NotFound")
	if !ref.IsReference() {
		t.Error("Expected response to be reference")
	}
	if ref.Ref != "#/responses/NotFound" {
		t.Errorf("Expected ref #/responses/NotFound, got %s", ref.Ref)
	}
}

func TestBodyParameter(t *testing.T) {
	param := &Parameter{
		Name:     "body",
		In:       "body",
		Required: true,
		Schema:   &Schema{Type: "object"},
	}
	if !param.IsBodyParameter() {
		t.Error("Expected parameter to be body parameter")
	}
}

func TestSecuritySchemeTypes(t *testing.T) {
	basic := &SecurityScheme{Type: "basic"}
	if !basic.IsBasic() {
		t.Error("Expected basic scheme")
	}

	apiKey := &SecurityScheme{Type: "apiKey", Name: "X-API-Key", In: "header"}
	if !apiKey.IsAPIKey() {
		t.Error("Expected apiKey scheme")
	}

	oauth := &SecurityScheme{Type: "oauth2", Flow: "implicit"}
	if !oauth.IsOAuth2() {
		t.Error("Expected oauth2 scheme")
	}
}
