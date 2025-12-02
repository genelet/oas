# OpenAPI 3.0 Go Package

A Go package for parsing, manipulating, and validating OpenAPI 3.0 specifications.

## Features

- Full OpenAPI 3.0.x specification support
- JSON marshaling/unmarshaling with round-trip preservation
- Boolean schema support for `additionalProperties: true/false`
- Extension fields (`x-*`) support on all applicable types
- Comprehensive validation against OpenAPI 3.0 specification
- Reference (`$ref`) support for all referenceable types

## Installation

```bash
go get github.com/genelet/oas/openapi30
```

## Usage

### Parsing an OpenAPI Document

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/genelet/oas/openapi30"
)

func main() {
    data, _ := os.ReadFile("openapi.json")

    var api openapi30.OpenAPI
    if err := json.Unmarshal(data, &api); err != nil {
        panic(err)
    }

    // Access the parsed document
    println("API Title:", api.Info.Title)
    println("Version:", api.Info.Version)
}
```

### Creating an OpenAPI Document

```go
api := &openapi30.OpenAPI{
    OpenAPI: "3.0.0",
    Info: &openapi30.Info{
        Title:   "My API",
        Version: "1.0.0",
    },
    Paths: &openapi30.Paths{
        Paths: map[string]*openapi30.PathItem{
            "/users": {
                Get: &openapi30.Operation{
                    Summary: "List users",
                    Responses: &openapi30.Responses{
                        StatusCode: map[string]*openapi30.Response{
                            "200": {
                                Description: "Successful response",
                                Content: map[string]*openapi30.MediaType{
                                    "application/json": {
                                        Schema: &openapi30.Schema{
                                            Type: "array",
                                            Items: &openapi30.Schema{
                                                Ref: "#/components/schemas/User",
                                            },
                                        },
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

// Marshal to JSON
data, _ := json.MarshalIndent(api, "", "  ")
println(string(data))
```

### Validation

```go
api := &openapi30.OpenAPI{
    // ... your API definition
}

result := api.Validate()
if !result.Valid() {
    for _, err := range result.Errors {
        fmt.Printf("%s: %s\n", err.Path, err.Message)
    }
}
```

### Boolean Schemas

OpenAPI 3.0 allows `additionalProperties` to be a boolean:

```go
// Create a schema with additionalProperties: false
schema := &openapi30.Schema{
    Type: "object",
    Properties: map[string]*openapi30.Schema{
        "name": {Type: "string"},
    },
    AdditionalProperties: openapi30.NewBooleanSchema(false),
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
paramRef := openapi30.NewParameterReference("#/components/parameters/PageSize")

// Check if something is a reference
if param.IsReference() {
    fmt.Printf("Reference: %s\n", param.Ref)
}
```

### Extension Fields

All types that support extensions in the OpenAPI spec have an `Extensions` field:

```go
api := &openapi30.OpenAPI{
    OpenAPI: "3.0.0",
    Info: &openapi30.Info{
        Title:   "My API",
        Version: "1.0.0",
        Extensions: map[string]any{
            "x-logo": map[string]any{
                "url": "https://example.com/logo.png",
            },
        },
    },
    // ...
}
```

## Type Reference

### Core Types

| Type | Description |
|------|-------------|
| `OpenAPI` | Root document object |
| `Info` | API metadata (title, version, description, etc.) |
| `Contact` | Contact information |
| `License` | License information |
| `Server` | Server URL and variables |
| `ServerVariable` | Server URL template variable |
| `Paths` | Container for path items |
| `PathItem` | Operations on a single path |
| `Operation` | Single API operation |
| `ExternalDocumentation` | External documentation link |
| `Tag` | Tag for API documentation |

### Request/Response Types

| Type | Description |
|------|-------------|
| `Parameter` | Operation parameter (path, query, header, cookie) |
| `Header` | Response header |
| `RequestBody` | Request body definition |
| `MediaType` | Media type with schema and examples |
| `Encoding` | Encoding for multipart request bodies |
| `Responses` | Container for response definitions |
| `Response` | Single response definition |
| `Link` | Design-time link for responses |
| `Callback` | Callback definition |
| `Example` | Example value |

### Schema Types

| Type | Description |
|------|-------------|
| `Schema` | JSON Schema with OpenAPI extensions |
| `Discriminator` | Polymorphism support |
| `XML` | XML serialization metadata |

### Security Types

| Type | Description |
|------|-------------|
| `SecurityScheme` | Security scheme definition |
| `OAuthFlows` | OAuth 2.0 flow configurations |
| `OAuthFlow` | Single OAuth flow configuration |
| `SecurityRequirement` | Security requirement for operations |

### Component Types

| Type | Description |
|------|-------------|
| `Components` | Reusable component definitions |

## Validation Rules

The `Validate()` method checks for:

### Required Fields
- `openapi` version string (must be 3.0.x)
- `info.title` and `info.version`
- `paths` object
- `response.description` (for non-reference responses)
- `requestBody.content` (for non-reference request bodies)

### Parameter Constraints
- Path parameters must have `required: true`
- Valid `style` values per parameter location:
  - `path`: matrix, label, simple
  - `query`: form, spaceDelimited, pipeDelimited, deepObject
  - `header`: simple
  - `cookie`: form
- Must have either `schema` or `content`, not both

### Schema Constraints
- Array type must have `items`
- Valid type values: string, number, integer, boolean, array, object
- `minimum` <= `maximum`
- `minLength` <= `maxLength`
- `minItems` <= `maxItems`
- `minProperties` <= `maxProperties`
- Valid regex patterns
- Required properties must exist in `properties`

### Security Scheme Constraints
- `apiKey`: requires `name` and `in`
- `http`: requires `scheme`
- `oauth2`: requires `flows` with appropriate URLs
- `openIdConnect`: requires `openIdConnectUrl`

### Other Constraints
- Link cannot have both `operationId` and `operationRef`
- Cannot have both `example` and `examples`
- Responses must contain at least one response
- Component names must match pattern `^[a-zA-Z0-9.\-_]+$`

## OpenAPI 3.0 vs 3.1 Differences

This package implements OpenAPI 3.0 (JSON Schema Draft 4). Key differences from OpenAPI 3.1:

| Feature | OpenAPI 3.0 | OpenAPI 3.1 |
|---------|-------------|-------------|
| `type` | Single string | String or array |
| `exclusiveMinimum/Maximum` | Boolean | Number |
| `nullable` | Supported | Use `type: ["string", "null"]` |
| `example` | Supported | Deprecated, use `examples` |
| JSON Schema | Draft 4 subset | Draft 2020-12 |

## Testing

```bash
# Run all tests
go test ./openapi30/...

# Run with verbose output
go test -v ./openapi30/...

# Run specific test
go test -run TestJSONRoundTrip ./openapi30/...
```

## License

See the repository root for license information.
