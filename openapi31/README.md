# OpenAPI 3.1 Go Package

A Go package for parsing, manipulating, and validating OpenAPI 3.1 specifications.

## Features

- Full OpenAPI 3.1.x specification support
- Complete JSON Schema Draft 2020-12 support
- JSON marshaling/unmarshaling with round-trip preservation
- Boolean schema support for `additionalProperties: true/false`
- Type arrays support (`["string", "null"]`)
- Webhooks support (OpenAPI 3.1 feature)
- Extension fields (`x-*`) support on all applicable types
- Comprehensive validation against OpenAPI 3.1 specification
- Reference (`$ref`) support with summary and description

## Installation

```bash
go get github.com/genelet/oas/openapi31
```

## Usage

### Parsing an OpenAPI Document

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/genelet/oas/openapi31"
)

func main() {
    data, _ := os.ReadFile("openapi.json")

    var api openapi31.OpenAPI
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
api := &openapi31.OpenAPI{
    OpenAPI: "3.1.0",
    Info: &openapi31.Info{
        Title:   "My API",
        Version: "1.0.0",
    },
    Paths: &openapi31.Paths{
        Paths: map[string]*openapi31.PathItem{
            "/users": {
                Get: &openapi31.Operation{
                    Summary: "List users",
                    Responses: &openapi31.Responses{
                        StatusCode: map[string]*openapi31.Response{
                            "200": {
                                Description: "Successful response",
                                Content: map[string]*openapi31.MediaType{
                                    "application/json": {
                                        Schema: &openapi31.Schema{
                                            Type: &openapi31.StringOrStringArray{String: "array"},
                                            Items: &openapi31.Schema{
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
api := &openapi31.OpenAPI{
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

OpenAPI 3.1 allows schemas to be boolean values:

```go
// Create a schema with additionalProperties: false
schema := &openapi31.Schema{
    Type: &openapi31.StringOrStringArray{String: "object"},
    Properties: map[string]*openapi31.Schema{
        "name": {Type: &openapi31.StringOrStringArray{String: "string"}},
    },
    AdditionalProperties: openapi31.NewBooleanSchema(false),
}

// Check if a schema is a boolean schema
if schema.AdditionalProperties.IsBooleanSchema() {
    value := schema.AdditionalProperties.BooleanValue()
    fmt.Printf("additionalProperties: %v\n", *value)
}
```

### Type Arrays (Nullable Types)

OpenAPI 3.1 uses JSON Schema 2020-12 style for nullable types:

```go
// Create a nullable string schema
schema := &openapi31.Schema{
    Type: &openapi31.StringOrStringArray{
        Array: []string{"string", "null"},
    },
}

// Check if type contains a specific value
if schema.Type.Contains("null") {
    println("Schema is nullable")
}
```

### Webhooks

OpenAPI 3.1 introduces webhooks:

```go
api := &openapi31.OpenAPI{
    OpenAPI: "3.1.0",
    Info:    &openapi31.Info{Title: "My API", Version: "1.0"},
    Webhooks: map[string]*openapi31.PathItem{
        "newUser": {
            Post: &openapi31.Operation{
                Summary: "New user webhook",
                RequestBody: &openapi31.RequestBody{
                    Content: map[string]*openapi31.MediaType{
                        "application/json": {
                            Schema: &openapi31.Schema{
                                Ref: "#/components/schemas/User",
                            },
                        },
                    },
                },
                Responses: &openapi31.Responses{
                    StatusCode: map[string]*openapi31.Response{
                        "200": {Description: "Webhook processed"},
                    },
                },
            },
        },
    },
}
```

### Working with References

```go
// Create a reference with summary and description (OpenAPI 3.1 feature)
paramRef := openapi31.NewParameterReference("#/components/parameters/PageSize")

// Check if something is a reference
if param.IsReference() {
    fmt.Printf("Reference: %s\n", param.Ref)
}
```

### Extension Fields

All types that support extensions in the OpenAPI spec have an `Extensions` field:

```go
api := &openapi31.OpenAPI{
    OpenAPI: "3.1.0",
    Info: &openapi31.Info{
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
| `Info` | API metadata (title, version, summary, description, etc.) |
| `Contact` | Contact information |
| `License` | License information (with identifier support) |
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
| `Schema` | JSON Schema Draft 2020-12 with OpenAPI extensions |
| `StringOrStringArray` | Union type for `type` field |
| `Discriminator` | Polymorphism support |
| `XML` | XML serialization metadata |

### Security Types

| Type | Description |
|------|-------------|
| `SecurityScheme` | Security scheme definition (includes mutualTLS) |
| `OAuthFlows` | OAuth 2.0 flow configurations |
| `OAuthFlow` | Single OAuth flow configuration |
| `SecurityRequirement` | Security requirement for operations |

### Component Types

| Type | Description |
|------|-------------|
| `Components` | Reusable component definitions (includes pathItems) |

## JSON Schema 2020-12 Keywords

The Schema type supports all JSON Schema Draft 2020-12 keywords:

### Core Keywords
- `$id`, `$schema`, `$ref`, `$anchor`, `$dynamicRef`, `$dynamicAnchor`
- `$defs`, `$comment`, `$vocabulary`

### Applicator Keywords
- `allOf`, `anyOf`, `oneOf`, `not`
- `if`, `then`, `else`
- `dependentSchemas`
- `prefixItems`, `items`, `contains`
- `properties`, `patternProperties`, `additionalProperties`, `propertyNames`
- `unevaluatedItems`, `unevaluatedProperties`

### Validation Keywords
- `type` (string or array), `enum`, `const`
- `multipleOf`, `maximum`, `exclusiveMaximum`, `minimum`, `exclusiveMinimum`
- `maxLength`, `minLength`, `pattern`
- `maxItems`, `minItems`, `uniqueItems`, `maxContains`, `minContains`
- `maxProperties`, `minProperties`, `required`, `dependentRequired`

### Other Keywords
- `title`, `description`, `default`, `deprecated`, `readOnly`, `writeOnly`
- `examples`, `format`
- `contentEncoding`, `contentMediaType`, `contentSchema`

## Validation Rules

The `Validate()` method checks for:

### Required Fields
- `openapi` version string (must be 3.1.x)
- `info.title` and `info.version`
- At least one of: `paths`, `webhooks`, or `components`
- `response.description` (for non-reference responses)
- `requestBody.content` (for non-reference request bodies)

### License Constraints
- `identifier` and `url` are mutually exclusive

### Parameter Constraints
- Path parameters must have `required: true`
- Valid `style` values per parameter location
- Must have either `schema` or `content`, not both

### Schema Constraints
- Valid type values: string, number, integer, boolean, array, object, null
- `minimum` <= `maximum`, `exclusiveMinimum` < `exclusiveMaximum`
- `minLength` <= `maxLength`, `minItems` <= `maxItems`
- Valid regex patterns
- Required properties must exist in `properties`

### Security Scheme Constraints
- `apiKey`: requires `name` and `in`
- `http`: requires `scheme`
- `oauth2`: requires `flows` with appropriate URLs
- `openIdConnect`: requires `openIdConnectUrl`
- `mutualTLS`: no additional requirements

### Other Constraints
- Link cannot have both `operationId` and `operationRef`
- Cannot have both `example` and `examples`
- Responses must contain at least one response
- Component names must match pattern `^[a-zA-Z0-9.\-_]+$`

## OpenAPI 3.1 vs 3.0 Differences

| Feature | OpenAPI 3.0 | OpenAPI 3.1 |
|---------|-------------|-------------|
| JSON Schema | Draft 4 subset | Draft 2020-12 |
| `type` | Single string | String or array |
| `exclusiveMinimum/Maximum` | Boolean | Number |
| `nullable` | Supported | Use `type: ["string", "null"]` |
| `example` | Preferred | Deprecated, use `examples` |
| Webhooks | Not supported | Supported |
| `pathItems` in Components | Not supported | Supported |
| License `identifier` | Not supported | Supported (SPDX) |
| `mutualTLS` security | Not supported | Supported |
| `$defs` | Not supported | Supported |
| `prefixItems` | Not supported | Supported |

## Testing

```bash
# Run all tests
go test ./openapi31/...

# Run with verbose output
go test -v ./openapi31/...

# Run specific test
go test -run TestJSONRoundTrip ./openapi31/...

# Run validation tests
go test -run TestValidate ./openapi31/...
```

## License

See the repository root for license information.
