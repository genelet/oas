# OAS - OpenAPI Specification Go Packages

Go packages for parsing, manipulating, and validating OpenAPI/Swagger specifications.

[![GoDoc](https://godoc.org/github.com/genelet/oas?status.svg)](https://godoc.org/github.com/genelet/oas)

## Overview

This repository provides Go packages for working with OpenAPI specifications across all major versions:

| Package | OpenAPI Version | JSON Schema Base |
|---------|-----------------|------------------|
| [openapi20](./openapi20/) | 2.0 (Swagger) | JSON Schema Draft 4 subset |
| [openapi30](./openapi30/) | 3.0.x | JSON Schema Draft 4 |
| [openapi31](./openapi31/) | 3.1.x | JSON Schema Draft 2020-12 |

All packages are built directly from the official JSON Schema specifications of their respective OpenAPI versions, ensuring complete and accurate type definitions.

## Installation

```bash
# Swagger 2.0 (OpenAPI 2.0)
go get github.com/genelet/oas/openapi20

# OpenAPI 3.0
go get github.com/genelet/oas/openapi30

# OpenAPI 3.1
go get github.com/genelet/oas/openapi31
```

## Features

- Full OpenAPI specification support for 2.0, 3.0, and 3.1
- Zero external dependencies - uses only Go standard library
- JSON marshaling/unmarshaling with round-trip preservation
- Boolean schema support (`additionalProperties: true/false`)
- Extension fields (`x-*`) on all applicable types
- Comprehensive validation against specifications (3.0, 3.1)
- Reference (`$ref`) support for all referenceable types

### OpenAPI 3.1 Specific Features

- Complete JSON Schema Draft 2020-12 support
- Type arrays (`["string", "null"]`)
- Webhooks
- `mutualTLS` security scheme
- License `identifier` field (SPDX)
- `pathItems` in Components

### Swagger 2.0 Specific Features

- `host`, `basePath`, `schemes` for server configuration
- Body parameters with `in: body`
- `consumes` and `produces` for content types
- `securityDefinitions` with basic, apiKey, and oauth2 flows
- `definitions` for reusable schemas

## Quick Start

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/genelet/oas/openapi31"
)

func main() {
    // Parse an OpenAPI document
    data, _ := os.ReadFile("openapi.json")

    var api openapi31.OpenAPI
    if err := json.Unmarshal(data, &api); err != nil {
        panic(err)
    }

    // Validate the document
    result := api.Validate()
    if !result.Valid() {
        for _, err := range result.Errors {
            fmt.Printf("%s: %s\n", err.Path, err.Message)
        }
        return
    }

    fmt.Printf("API: %s v%s\n", api.Info.Title, api.Info.Version)
}
```

## Documentation

- [Swagger 2.0 Package Documentation](./openapi20/README.md)
- [OpenAPI 3.0 Package Documentation](./openapi30/README.md)
- [OpenAPI 3.1 Package Documentation](./openapi31/README.md)

## Code Generation

Most of the code in this repository was generated using [Claude Code](https://claude.ai/code), Anthropic's AI-powered coding assistant. The packages were built by:

1. Analyzing the official OpenAPI JSON Schema specifications (`schema.json`)
2. Generating Go struct definitions with appropriate JSON tags
3. Implementing custom marshal/unmarshal methods for complex types
4. Creating comprehensive validation logic
5. Writing round-trip tests against real-world OpenAPI examples

## Specification References

- [Swagger 2.0 Specification](https://swagger.io/specification/v2/)
- [OpenAPI 3.0 Specification](https://spec.openapis.org/oas/v3.0.3)
- [OpenAPI 3.1 Specification](https://spec.openapis.org/oas/v3.1.0)
- [JSON Schema Draft 4](https://json-schema.org/specification-links#draft-4)
- [JSON Schema Draft 2020-12](https://json-schema.org/specification-links#2020-12)

## License

See [LICENSE](./LICENSE) for details.
