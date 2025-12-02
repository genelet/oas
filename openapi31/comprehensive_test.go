// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import (
	"encoding/json"
	"testing"
)

// TestComprehensiveOpenAPI31 tests a full-featured OpenAPI 3.1 document
func TestComprehensiveOpenAPI31(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {
			"title": "Comprehensive Test API",
			"summary": "A test API",
			"description": "Testing all OpenAPI 3.1 features",
			"termsOfService": "https://example.com/terms",
			"contact": {
				"name": "API Support",
				"url": "https://example.com/support",
				"email": "support@example.com",
				"x-contact-ext": "contact extension"
			},
			"license": {
				"name": "Apache 2.0",
				"identifier": "Apache-2.0",
				"url": "https://www.apache.org/licenses/LICENSE-2.0",
				"x-license-ext": true
			},
			"version": "1.0.0",
			"x-info-ext": {"nested": "value"}
		},
		"jsonSchemaDialect": "https://json-schema.org/draft/2020-12/schema",
		"servers": [
			{
				"url": "https://api.example.com/{version}",
				"description": "Production server",
				"variables": {
					"version": {
						"default": "v1",
						"enum": ["v1", "v2"],
						"description": "API version"
					}
				},
				"x-server-ext": 123
			}
		],
		"paths": {
			"/pets": {
				"summary": "Pet operations",
				"description": "Operations on pets",
				"parameters": [
					{"$ref": "#/components/parameters/limitParam"}
				],
				"get": {
					"tags": ["pets"],
					"summary": "List pets",
					"description": "Returns all pets",
					"operationId": "listPets",
					"parameters": [
						{
							"name": "status",
							"in": "query",
							"description": "Filter by status",
							"required": false,
							"schema": {"type": "string", "enum": ["available", "pending", "sold"]},
							"x-param-ext": "value"
						}
					],
					"responses": {
						"200": {
							"description": "Successful response",
							"headers": {
								"X-Rate-Limit": {
									"description": "Rate limit",
									"schema": {"type": "integer"}
								}
							},
							"content": {
								"application/json": {
									"schema": {"$ref": "#/components/schemas/PetList"},
									"examples": {
										"dogs": {
											"summary": "Dogs example",
											"value": [{"name": "Fido", "type": "dog"}]
										}
									}
								}
							},
							"links": {
								"GetPetById": {
									"operationId": "getPet",
									"parameters": {"petId": "$response.body#/0/id"}
								}
							}
						},
						"default": {"$ref": "#/components/responses/Error"}
					},
					"callbacks": {
						"petCallback": {
							"{$request.body#/callbackUrl}": {
								"post": {
									"responses": {
										"200": {"description": "Callback received"}
									}
								}
							}
						}
					},
					"security": [{"apiKey": []}],
					"x-operation-ext": true
				},
				"post": {
					"summary": "Create pet",
					"operationId": "createPet",
					"requestBody": {
						"description": "Pet to create",
						"required": true,
						"content": {
							"application/json": {
								"schema": {"$ref": "#/components/schemas/Pet"}
							}
						}
					},
					"responses": {
						"201": {"description": "Created"}
					}
				},
				"x-path-ext": "path extension"
			},
			"/pets/{petId}": {
				"$ref": "#/components/pathItems/PetItem"
			}
		},
		"webhooks": {
			"newPet": {
				"post": {
					"summary": "New pet webhook",
					"requestBody": {
						"content": {
							"application/json": {
								"schema": {"$ref": "#/components/schemas/Pet"}
							}
						}
					},
					"responses": {
						"200": {"description": "Webhook received"}
					}
				}
			}
		},
		"components": {
			"schemas": {
				"Pet": {
					"type": "object",
					"required": ["name"],
					"properties": {
						"id": {"type": "integer", "format": "int64", "readOnly": true},
						"name": {"type": "string", "minLength": 1, "maxLength": 100},
						"type": {"type": "string", "enum": ["dog", "cat", "bird"]},
						"tags": {
							"type": "array",
							"items": {"type": "string"},
							"uniqueItems": true
						},
						"metadata": {
							"type": "object",
							"additionalProperties": true
						}
					},
					"x-schema-ext": "schema extension"
				},
				"PetList": {
					"type": "array",
					"items": {"$ref": "#/components/schemas/Pet"}
				},
				"NullableString": {
					"type": ["string", "null"],
					"description": "A nullable string using type array"
				},
				"AnyValue": true,
				"NeverValid": false,
				"Error": {
					"type": "object",
					"properties": {
						"code": {"type": "integer"},
						"message": {"type": "string"}
					}
				}
			},
			"responses": {
				"Error": {
					"description": "Error response",
					"content": {
						"application/json": {
							"schema": {"$ref": "#/components/schemas/Error"}
						}
					}
				}
			},
			"parameters": {
				"limitParam": {
					"name": "limit",
					"in": "query",
					"schema": {"type": "integer", "minimum": 1, "maximum": 100, "default": 20}
				}
			},
			"examples": {
				"dogExample": {
					"summary": "A dog",
					"value": {"name": "Buddy", "type": "dog"}
				}
			},
			"requestBodies": {
				"PetBody": {
					"description": "Pet request body",
					"content": {
						"application/json": {
							"schema": {"$ref": "#/components/schemas/Pet"}
						}
					}
				}
			},
			"headers": {
				"X-Request-Id": {
					"description": "Request ID",
					"schema": {"type": "string", "format": "uuid"}
				}
			},
			"securitySchemes": {
				"apiKey": {
					"type": "apiKey",
					"name": "X-API-Key",
					"in": "header"
				},
				"oauth2": {
					"type": "oauth2",
					"flows": {
						"authorizationCode": {
							"authorizationUrl": "https://example.com/oauth/authorize",
							"tokenUrl": "https://example.com/oauth/token",
							"scopes": {
								"read:pets": "Read pets",
								"write:pets": "Write pets"
							}
						}
					}
				}
			},
			"links": {
				"GetPetById": {
					"operationId": "getPet",
					"parameters": {"petId": "$response.body#/id"}
				}
			},
			"callbacks": {
				"petEvent": {
					"{$request.body#/callbackUrl}": {
						"post": {
							"responses": {"200": {"description": "OK"}}
						}
					}
				}
			},
			"pathItems": {
				"PetItem": {
					"get": {
						"summary": "Get pet by ID",
						"operationId": "getPet",
						"parameters": [
							{"name": "petId", "in": "path", "required": true, "schema": {"type": "string"}}
						],
						"responses": {
							"200": {
								"description": "Pet found",
								"content": {
									"application/json": {
										"schema": {"$ref": "#/components/schemas/Pet"}
									}
								}
							}
						}
					}
				}
			}
		},
		"security": [
			{"apiKey": []},
			{"oauth2": ["read:pets"]}
		],
		"tags": [
			{
				"name": "pets",
				"description": "Pet operations",
				"externalDocs": {
					"description": "Find more info here",
					"url": "https://example.com/docs"
				},
				"x-tag-ext": "tag extension"
			}
		],
		"externalDocs": {
			"description": "Full documentation",
			"url": "https://example.com/docs",
			"x-docs-ext": true
		},
		"x-root-ext": "root extension"
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Test root level
	if api.OpenAPI != "3.1.0" {
		t.Errorf("Expected openapi '3.1.0', got %s", api.OpenAPI)
	}
	if api.JsonSchemaDialect != "https://json-schema.org/draft/2020-12/schema" {
		t.Errorf("Expected jsonSchemaDialect, got %s", api.JsonSchemaDialect)
	}
	if api.Extensions["x-root-ext"] != "root extension" {
		t.Errorf("Expected root extension")
	}

	// Test Info
	if api.Info.Title != "Comprehensive Test API" {
		t.Errorf("Expected title")
	}
	if api.Info.Summary != "A test API" {
		t.Errorf("Expected summary")
	}
	if api.Info.License.Identifier != "Apache-2.0" {
		t.Errorf("Expected license identifier")
	}
	if api.Info.Contact.Extensions["x-contact-ext"] != "contact extension" {
		t.Errorf("Expected contact extension")
	}

	// Test Servers
	if len(api.Servers) != 1 {
		t.Fatalf("Expected 1 server")
	}
	if api.Servers[0].Variables["version"].Default != "v1" {
		t.Errorf("Expected server variable default")
	}

	// Test Paths
	petsPath := api.Paths.Get("/pets")
	if petsPath == nil {
		t.Fatal("Expected /pets path")
	}
	if petsPath.Summary != "Pet operations" {
		t.Errorf("Expected path summary")
	}

	// Test path-level parameter reference
	if len(petsPath.Parameters) != 1 {
		t.Fatalf("Expected 1 path parameter")
	}
	if !petsPath.Parameters[0].IsReference() {
		t.Error("Expected path parameter to be reference")
	}

	// Test Operation
	getOp := petsPath.Get
	if getOp == nil {
		t.Fatal("Expected GET operation")
	}
	if getOp.OperationID != "listPets" {
		t.Errorf("Expected operationId")
	}
	if getOp.Extensions["x-operation-ext"] != true {
		t.Errorf("Expected operation extension")
	}

	// Test inline parameter
	if len(getOp.Parameters) != 1 {
		t.Fatalf("Expected 1 operation parameter")
	}
	if getOp.Parameters[0].Name != "status" {
		t.Errorf("Expected parameter name 'status'")
	}
	if getOp.Parameters[0].IsReference() {
		t.Error("Expected inline parameter, not reference")
	}

	// Test Response reference
	defaultResp := getOp.Responses.GetDefault()
	if defaultResp == nil {
		t.Fatal("Expected default response")
	}
	if !defaultResp.IsReference() {
		t.Error("Expected default response to be reference")
	}
	if defaultResp.Ref != "#/components/responses/Error" {
		t.Errorf("Expected response ref")
	}

	// Test inline response with headers and links
	resp200 := getOp.Responses.Get("200")
	if resp200 == nil {
		t.Fatal("Expected 200 response")
	}
	if resp200.Headers["X-Rate-Limit"] == nil {
		t.Error("Expected header")
	}
	if resp200.Links["GetPetById"] == nil {
		t.Error("Expected link")
	}

	// Test Callbacks
	if getOp.Callbacks == nil || getOp.Callbacks["petCallback"] == nil {
		t.Error("Expected callback")
	}

	// Test RequestBody
	postOp := petsPath.Post
	if postOp.RequestBody == nil {
		t.Fatal("Expected request body")
	}
	if !postOp.RequestBody.Required {
		t.Error("Expected required request body")
	}

	// Test PathItem $ref
	petIdPath := api.Paths.Get("/pets/{petId}")
	if petIdPath == nil {
		t.Fatal("Expected /pets/{petId} path")
	}
	if petIdPath.Ref != "#/components/pathItems/PetItem" {
		t.Errorf("Expected path item ref")
	}

	// Test Webhooks
	if api.Webhooks == nil || api.Webhooks["newPet"] == nil {
		t.Error("Expected webhook")
	}
	if api.Webhooks["newPet"].Post == nil {
		t.Error("Expected webhook POST operation")
	}

	// Test Components Schemas
	petSchema := api.Components.Schemas["Pet"]
	if petSchema == nil {
		t.Fatal("Expected Pet schema")
	}
	if !petSchema.Type.Contains("object") {
		t.Errorf("Expected object type")
	}
	if len(petSchema.Required) != 1 || petSchema.Required[0] != "name" {
		t.Errorf("Expected required field")
	}

	// Test nullable type array
	nullableSchema := api.Components.Schemas["NullableString"]
	if nullableSchema == nil {
		t.Fatal("Expected NullableString schema")
	}
	if !nullableSchema.Type.Contains("string") || !nullableSchema.Type.Contains("null") {
		t.Errorf("Expected type array with string and null")
	}

	// Test boolean schemas
	anySchema := api.Components.Schemas["AnyValue"]
	if !anySchema.IsBooleanSchema() || *anySchema.BooleanValue() != true {
		t.Error("Expected boolean true schema")
	}
	neverSchema := api.Components.Schemas["NeverValid"]
	if !neverSchema.IsBooleanSchema() || *neverSchema.BooleanValue() != false {
		t.Error("Expected boolean false schema")
	}

	// Test Components other types
	if api.Components.Responses["Error"] == nil {
		t.Error("Expected Error response component")
	}
	if api.Components.Parameters["limitParam"] == nil {
		t.Error("Expected limitParam component")
	}
	if api.Components.Examples["dogExample"] == nil {
		t.Error("Expected dogExample component")
	}
	if api.Components.RequestBodies["PetBody"] == nil {
		t.Error("Expected PetBody component")
	}
	if api.Components.Headers["X-Request-Id"] == nil {
		t.Error("Expected header component")
	}
	if api.Components.SecuritySchemes["apiKey"] == nil {
		t.Error("Expected apiKey security scheme")
	}
	if api.Components.SecuritySchemes["oauth2"] == nil {
		t.Error("Expected oauth2 security scheme")
	}
	if api.Components.Links["GetPetById"] == nil {
		t.Error("Expected link component")
	}
	if api.Components.Callbacks["petEvent"] == nil {
		t.Error("Expected callback component")
	}
	if api.Components.PathItems["PetItem"] == nil {
		t.Error("Expected pathItem component")
	}

	// Test OAuth2 flows
	oauth2 := api.Components.SecuritySchemes["oauth2"]
	if oauth2.Flows == nil || oauth2.Flows.AuthorizationCode == nil {
		t.Error("Expected OAuth2 authorization code flow")
	}
	if oauth2.Flows.AuthorizationCode.Scopes["read:pets"] != "Read pets" {
		t.Error("Expected OAuth2 scope")
	}

	// Test Security
	if len(api.Security) != 2 {
		t.Errorf("Expected 2 security requirements")
	}

	// Test Tags
	if len(api.Tags) != 1 || api.Tags[0].Name != "pets" {
		t.Error("Expected pets tag")
	}
	if api.Tags[0].ExternalDocs == nil {
		t.Error("Expected tag external docs")
	}

	// Test ExternalDocs
	if api.ExternalDocs == nil || api.ExternalDocs.URL != "https://example.com/docs" {
		t.Error("Expected external docs")
	}
}

// TestRoundTripComprehensive tests marshal/unmarshal round-trip
func TestRoundTripComprehensive(t *testing.T) {
	original := &OpenAPI{
		OpenAPI: "3.1.0",
		Info: &Info{
			Title:   "Round Trip Test",
			Version: "1.0.0",
			License: &License{
				Name:       "MIT",
				Identifier: "MIT",
			},
		},
		JsonSchemaDialect: "https://json-schema.org/draft/2020-12/schema",
		Paths: &Paths{
			Paths: map[string]*PathItem{
				"/test": {
					Get: &Operation{
						OperationID: "test",
						Parameters: []*Parameter{
							{Name: "q", In: "query", Schema: &Schema{Type: &StringOrStringArray{String: "string"}}},
							NewParameterReference("#/components/parameters/common"),
						},
						Responses: &Responses{
							StatusCode: map[string]*Response{
								"200": {Description: "OK"},
							},
							Default: NewResponseReference("#/components/responses/Error"),
						},
					},
				},
			},
		},
		Components: &Components{
			Schemas: map[string]*Schema{
				"Test": {
					Type:       &StringOrStringArray{String: "object"},
					Properties: map[string]*Schema{
						"name": {Type: &StringOrStringArray{String: "string"}},
						"nullable": {Type: &StringOrStringArray{Array: []string{"string", "null"}}},
					},
					AdditionalProperties: NewBooleanSchema(false),
				},
				"Any":  NewBooleanSchema(true),
				"None": NewBooleanSchema(false),
			},
			Parameters: map[string]*Parameter{
				"common": {Name: "common", In: "query"},
			},
			Responses: map[string]*Response{
				"Error": {Description: "Error"},
			},
		},
		Webhooks: map[string]*PathItem{
			"event": {
				Post: &Operation{
					Responses: &Responses{
						StatusCode: map[string]*Response{
							"200": {Description: "OK"},
						},
					},
				},
			},
		},
		Extensions: map[string]any{
			"x-custom": "value",
		},
	}

	// Marshal
	data, err := json.MarshalIndent(original, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var parsed OpenAPI
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify
	if parsed.OpenAPI != "3.1.0" {
		t.Error("Round-trip failed: openapi")
	}
	if parsed.JsonSchemaDialect != "https://json-schema.org/draft/2020-12/schema" {
		t.Error("Round-trip failed: jsonSchemaDialect")
	}
	if parsed.Info.License.Identifier != "MIT" {
		t.Error("Round-trip failed: license identifier")
	}

	// Test parameters
	params := parsed.Paths.Get("/test").Get.Parameters
	if len(params) != 2 {
		t.Fatalf("Round-trip failed: expected 2 params, got %d", len(params))
	}
	if params[0].IsReference() {
		t.Error("Round-trip failed: first param should not be reference")
	}
	if params[0].Name != "q" {
		t.Error("Round-trip failed: first param name")
	}
	if !params[1].IsReference() {
		t.Error("Round-trip failed: second param should be reference")
	}
	if params[1].Ref != "#/components/parameters/common" {
		t.Error("Round-trip failed: second param ref")
	}

	// Test response reference
	defaultResp := parsed.Paths.Get("/test").Get.Responses.GetDefault()
	if !defaultResp.IsReference() {
		t.Error("Round-trip failed: default response should be reference")
	}

	// Test boolean schemas
	testSchema := parsed.Components.Schemas["Test"]
	if testSchema.AdditionalProperties == nil || !testSchema.AdditionalProperties.IsBooleanSchema() {
		t.Error("Round-trip failed: additionalProperties boolean schema")
	}
	if *testSchema.AdditionalProperties.BooleanValue() != false {
		t.Error("Round-trip failed: additionalProperties should be false")
	}

	anySchema := parsed.Components.Schemas["Any"]
	if !anySchema.IsBooleanSchema() || *anySchema.BooleanValue() != true {
		t.Error("Round-trip failed: Any schema")
	}

	// Test type array
	nullableProp := testSchema.Properties["nullable"]
	if nullableProp.Type.Array == nil || len(nullableProp.Type.Array) != 2 {
		t.Error("Round-trip failed: nullable type array")
	}

	// Test webhooks
	if parsed.Webhooks == nil || parsed.Webhooks["event"] == nil {
		t.Error("Round-trip failed: webhooks")
	}

	// Test extensions
	if parsed.Extensions["x-custom"] != "value" {
		t.Error("Round-trip failed: extensions")
	}
}

// TestSchemaJSONSchema2020Features tests JSON Schema 2020-12 specific features
func TestSchemaJSONSchema2020Features(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0"},
		"components": {
			"schemas": {
				"Advanced": {
					"$id": "https://example.com/schemas/advanced",
					"$schema": "https://json-schema.org/draft/2020-12/schema",
					"$anchor": "advancedSchema",
					"$comment": "This is a test schema",
					"$defs": {
						"subSchema": {"type": "string"}
					},
					"type": "object",
					"properties": {
						"name": {"type": "string"},
						"items": {
							"type": "array",
							"prefixItems": [
								{"type": "string"},
								{"type": "integer"}
							],
							"items": {"type": "boolean"},
							"contains": {"type": "string", "minLength": 1},
							"minContains": 1,
							"maxContains": 5
						},
						"conditional": {
							"if": {"properties": {"type": {"const": "premium"}}},
							"then": {"required": ["premium_field"]},
							"else": {"required": ["basic_field"]}
						}
					},
					"dependentSchemas": {
						"credit_card": {
							"required": ["billing_address"]
						}
					},
					"dependentRequired": {
						"email": ["name"]
					},
					"propertyNames": {"pattern": "^[a-z]+$"},
					"unevaluatedProperties": false,
					"unevaluatedItems": false
				}
			}
		}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	schema := api.Components.Schemas["Advanced"]
	if schema == nil {
		t.Fatal("Expected Advanced schema")
	}

	// Test core keywords
	if schema.ID != "https://example.com/schemas/advanced" {
		t.Errorf("Expected $id")
	}
	if schema.Schema != "https://json-schema.org/draft/2020-12/schema" {
		t.Errorf("Expected $schema")
	}
	if schema.Anchor != "advancedSchema" {
		t.Errorf("Expected $anchor")
	}
	if schema.Comment != "This is a test schema" {
		t.Errorf("Expected $comment")
	}

	// Test $defs
	if schema.Defs == nil || schema.Defs["subSchema"] == nil {
		t.Error("Expected $defs")
	}

	// Test prefixItems
	itemsSchema := schema.Properties["items"]
	if itemsSchema == nil {
		t.Fatal("Expected items property")
	}
	if len(itemsSchema.PrefixItems) != 2 {
		t.Errorf("Expected 2 prefixItems")
	}

	// Test contains with minContains/maxContains
	if itemsSchema.Contains == nil {
		t.Error("Expected contains")
	}
	if itemsSchema.MinContains == nil || *itemsSchema.MinContains != 1 {
		t.Error("Expected minContains")
	}
	if itemsSchema.MaxContains == nil || *itemsSchema.MaxContains != 5 {
		t.Error("Expected maxContains")
	}

	// Test if/then/else
	condSchema := schema.Properties["conditional"]
	if condSchema == nil {
		t.Fatal("Expected conditional property")
	}
	if condSchema.If == nil || condSchema.Then == nil || condSchema.Else == nil {
		t.Error("Expected if/then/else")
	}

	// Test dependentSchemas
	if schema.DependentSchemas == nil || schema.DependentSchemas["credit_card"] == nil {
		t.Error("Expected dependentSchemas")
	}

	// Test dependentRequired
	if schema.DependentRequired == nil || schema.DependentRequired["email"] == nil {
		t.Error("Expected dependentRequired")
	}

	// Test propertyNames
	if schema.PropertyNames == nil {
		t.Error("Expected propertyNames")
	}

	// Test unevaluated*
	if schema.UnevaluatedProperties == nil || !schema.UnevaluatedProperties.IsBooleanSchema() {
		t.Error("Expected unevaluatedProperties")
	}
	if schema.UnevaluatedItems == nil || !schema.UnevaluatedItems.IsBooleanSchema() {
		t.Error("Expected unevaluatedItems")
	}
}

// TestAllReferenceTypes tests all types that support $ref
func TestAllReferenceTypes(t *testing.T) {
	doc := `{
		"openapi": "3.1.0",
		"info": {"title": "Test", "version": "1.0"},
		"paths": {
			"/test": {
				"get": {
					"parameters": [
						{"$ref": "#/components/parameters/test", "summary": "Param ref"}
					],
					"requestBody": {"$ref": "#/components/requestBodies/test"},
					"responses": {
						"200": {"$ref": "#/components/responses/test"},
						"default": {
							"description": "Error",
							"headers": {
								"X-Test": {"$ref": "#/components/headers/test"}
							},
							"links": {
								"test": {"$ref": "#/components/links/test"}
							}
						}
					},
					"callbacks": {
						"test": {"$ref": "#/components/callbacks/test"}
					}
				}
			}
		},
		"components": {
			"securitySchemes": {
				"ref": {"$ref": "#/components/securitySchemes/actual"},
				"actual": {"type": "apiKey", "name": "key", "in": "header"}
			},
			"parameters": {"test": {"name": "test", "in": "query"}},
			"requestBodies": {"test": {"content": {"application/json": {}}}},
			"responses": {"test": {"description": "Test"}},
			"headers": {"test": {"schema": {"type": "string"}}},
			"links": {"test": {"operationId": "test"}},
			"callbacks": {"test": {"/callback": {"post": {"responses": {"200": {"description": "OK"}}}}}}
		}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	op := api.Paths.Get("/test").Get

	// Parameter reference
	if !op.Parameters[0].IsReference() {
		t.Error("Expected parameter reference")
	}
	if op.Parameters[0].Ref != "#/components/parameters/test" {
		t.Error("Expected parameter ref value")
	}

	// RequestBody reference
	if !op.RequestBody.IsReference() {
		t.Error("Expected requestBody reference")
	}

	// Response reference
	if !op.Responses.Get("200").IsReference() {
		t.Error("Expected response reference")
	}

	// Header reference
	defaultResp := op.Responses.GetDefault()
	if !defaultResp.Headers["X-Test"].IsReference() {
		t.Error("Expected header reference")
	}

	// Link reference
	if !defaultResp.Links["test"].IsReference() {
		t.Error("Expected link reference")
	}

	// Callback reference
	if !op.Callbacks["test"].IsReference() {
		t.Error("Expected callback reference")
	}

	// SecurityScheme reference
	if !api.Components.SecuritySchemes["ref"].IsReference() {
		t.Error("Expected securityScheme reference")
	}
}

// TestExampleReferenceInMediaType tests example references in media types
func TestExampleReferenceInMediaType(t *testing.T) {
	doc := `{
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
									"examples": {
										"inline": {"value": {"test": "data"}},
										"ref": {"$ref": "#/components/examples/test"}
									}
								}
							}
						}
					}
				}
			}
		},
		"components": {
			"examples": {
				"test": {"value": {"referenced": "example"}}
			}
		}
	}`

	var api OpenAPI
	if err := json.Unmarshal([]byte(doc), &api); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	examples := api.Paths.Get("/test").Get.Responses.Get("200").Content["application/json"].Examples

	// Inline example
	if examples["inline"].IsReference() {
		t.Error("Expected inline example, not reference")
	}

	// Reference example
	if !examples["ref"].IsReference() {
		t.Error("Expected example reference")
	}
	if examples["ref"].Ref != "#/components/examples/test" {
		t.Error("Expected example ref value")
	}
}
