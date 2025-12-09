# OpenAPI 2.0 (Swagger) Go Package

A Go package for parsing, manipulating, and validating Swagger/OpenAPI 2.0 specifications.

## Features

- Full Swagger 2.0 (OpenAPI 2.0) specification support
- JSON marshaling/unmarshaling with round-trip preservation
- Boolean schema support for `additionalProperties: true/false`
- Extension fields (`x-*`) support on all applicable types
- Reference (`$ref`) support for all referenceable types
- Zero external dependencies - uses only Go standard library

## Installation

```bash
go get github.com/genelet/oas/openapi20
```

## Usage

### Parsing a Swagger Document

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/genelet/oas/openapi20"
)

func main() {
    data, _ := os.ReadFile("swagger.json")

    var swagger openapi20.Swagger
    if err := json.Unmarshal(data, &swagger); err != nil {
        panic(err)
    }

    // Access the parsed document
    println("API Title:", swagger.Info.Title)
    println("Version:", swagger.Info.Version)
    println("Host:", swagger.Host)
    println("BasePath:", swagger.BasePath)
}
```

### Creating a Swagger Document

```go
swagger := &openapi20.Swagger{
    Swagger: "2.0",
    Info: &openapi20.Info{
        Title:   "My API",
        Version: "1.0.0",
    },
    Host:     "api.example.com",
    BasePath: "/v1",
    Schemes:  []string{"https"},
    Consumes: []string{"application/json"},
    Produces: []string{"application/json"},
    Paths: &openapi20.Paths{
        Paths: map[string]*openapi20.PathItem{
            "/users": {
                Get: &openapi20.Operation{
                    Summary:     "List users",
                    OperationID: "listUsers",
                    Responses: &openapi20.Responses{
                        StatusCode: map[string]*openapi20.Response{
                            "200": {
                                Description: "Successful response",
                                Schema: &openapi20.Schema{
                                    Type: "array",
                                    Items: &openapi20.Schema{
                                        Ref: "#/definitions/User",
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    },
    Definitions: map[string]*openapi20.Schema{
        "User": {
            Type:     "object",
            Required: []string{"id", "name"},
            Properties: map[string]*openapi20.Schema{
                "id":   {Type: "integer", Format: "int64"},
                "name": {Type: "string"},
            },
        },
    },
}

// Marshal to JSON
data, _ := json.MarshalIndent(swagger, "", "  ")
println(string(data))
```

### Boolean Schemas

Swagger 2.0 allows `additionalProperties` to be a boolean:

```go
// Create a schema with additionalProperties: false
schema := &openapi20.Schema{
    Type: "object",
    Properties: map[string]*openapi20.Schema{
        "name": {Type: "string"},
    },
    AdditionalProperties: openapi20.NewBooleanSchema(false),
}

// Check if a schema is a boolean schema
if schema.AdditionalProperties.IsBooleanSchema() {
    value := schema.AdditionalProperties.BooleanValue()
    fmt.Printf("additionalProperties: %v\n", *value)
}
```

### Working with References

```go
// Create a reference
schemaRef := openapi20.NewSchemaReference("#/definitions/Pet")
paramRef := openapi20.NewParameterReference("#/parameters/limitParam")
responseRef := openapi20.NewResponseReference("#/responses/NotFound")

// Check if something is a reference
if schema.IsReference() {
    fmt.Printf("Reference: %s\n", schema.Ref)
}
```

### Parameters

Swagger 2.0 has different parameter locations:

```go
// Path parameter
pathParam := &openapi20.Parameter{
    Name:     "id",
    In:       "path",
    Required: true,
    Type:     "string",
}

// Query parameter
queryParam := &openapi20.Parameter{
    Name: "limit",
    In:   "query",
    Type: "integer",
}

// Body parameter (uses schema instead of type)
bodyParam := &openapi20.Parameter{
    Name:     "body",
    In:       "body",
    Required: true,
    Schema:   &openapi20.Schema{Ref: "#/definitions/Pet"},
}

// Check parameter type
if bodyParam.IsBodyParameter() {
    fmt.Println("This is a body parameter")
}
```

### Security Definitions

```go
swagger := &openapi20.Swagger{
    // ...
    SecurityDefinitions: map[string]*openapi20.SecurityScheme{
        "api_key": {
            Type: "apiKey",
            Name: "X-API-Key",
            In:   "header",
        },
        "basic": {
            Type:        "basic",
            Description: "Basic HTTP authentication",
        },
        "oauth2": {
            Type:             "oauth2",
            Flow:             "accessCode",
            AuthorizationUrl: "https://example.com/oauth/authorize",
            TokenUrl:         "https://example.com/oauth/token",
            Scopes: map[string]string{
                "read":  "Read access",
                "write": "Write access",
            },
        },
    },
}
```

### Extension Fields

All types that support extensions in the Swagger spec have an `Extensions` field:

```go
swagger := &openapi20.Swagger{
    Swagger: "2.0",
    Info: &openapi20.Info{
        Title:   "My API",
        Version: "1.0.0",
        Extensions: map[string]any{
            "x-logo": map[string]any{
                "url": "https://example.com/logo.png",
            },
        },
    },
    Extensions: map[string]any{
        "x-api-id": "12345",
    },
    // ...
}
```

## Type Reference

### Core Types

| Type | Description |
|------|-------------|
| `Swagger` | Root document object |
| `Info` | API metadata (title, version, description, etc.) |
| `Contact` | Contact information |
| `License` | License information |
| `Paths` | Container for path items |
| `PathItem` | Operations on a single path |
| `Operation` | Single API operation |
| `ExternalDocumentation` | External documentation link |
| `Tag` | Tag for API documentation |

### Request/Response Types

| Type | Description |
|------|-------------|
| `Parameter` | Operation parameter (path, query, header, formData, body) |
| `Items` | Item type for array parameters |
| `Header` | Response header |
| `Responses` | Container for response definitions |
| `Response` | Single response definition |

### Schema Types

| Type | Description |
|------|-------------|
| `Schema` | JSON Schema subset with Swagger extensions |
| `XML` | XML serialization metadata |

### Security Types

| Type | Description |
|------|-------------|
| `SecurityScheme` | Security scheme definition (basic, apiKey, oauth2) |
| `SecurityRequirement` | Security requirement for operations |

## Swagger 2.0 vs OpenAPI 3.0 Differences

This package implements Swagger 2.0 (OpenAPI 2.0). Key differences from OpenAPI 3.0:

| Feature | Swagger 2.0 | OpenAPI 3.0 |
|---------|-------------|-------------|
| Version field | `swagger: "2.0"` | `openapi: "3.0.x"` |
| Server | `host`, `basePath`, `schemes` | `servers` array |
| Request body | `in: body` parameter | `requestBody` object |
| Content types | `consumes`, `produces` | `content` in request/response |
| Schemas | `definitions` | `components/schemas` |
| Parameters | `parameters` | `components/parameters` |
| Responses | `responses` | `components/responses` |
| Security | `securityDefinitions` | `components/securitySchemes` |
| Composition | Only `allOf` | `allOf`, `oneOf`, `anyOf`, `not` |
| File upload | `type: file` in formData | Binary string in requestBody |
| Discriminator | String (property name) | Object with mapping |

## Testing

```bash
# Run all tests
go test ./openapi20/...

# Run with verbose output
go test -v ./openapi20/...

# Run specific test
go test -run TestJSONRoundTrip ./openapi20/...
```

## License

See the repository root for license information.
