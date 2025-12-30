// Package unified provides unified interfaces for OpenAPI documents across versions
// Copyright (c) Greetingland LLC

package unified

// Document is a unified interface for OpenAPI documents of any version (2.0, 3.0, 3.1).
type Document interface {
	// Version returns the OpenAPI/Swagger version string (e.g., "2.0", "3.0.0", "3.1.0")
	Version() string

	// GetServerURL returns the base server URL
	GetServerURL() string

	// GetInfo returns document metadata
	GetInfo() DocumentInfo

	// GetPaths returns all path items keyed by path string
	GetPaths() map[string]PathItem

	// GetSecuritySchemes returns all security scheme definitions
	GetSecuritySchemes() map[string]SecurityScheme

	// GetGlobalSecurity returns document-level security requirements
	GetGlobalSecurity() []SecurityRequirement

	// GetExtensions returns the extensions map (x-...)
	GetExtensions() map[string]any
}

// DocumentInfo provides metadata about the API
type DocumentInfo interface {
	GetTitle() string
	GetVersion() string
	GetDescription() string
	GetExtensions() map[string]any
}

// PathItem abstracts a path item across OpenAPI versions
type PathItem interface {
	HasRef() bool
	GetRef() string
	GetOperation(method string) Operation
	GetAllOperations() map[string]Operation
	GetParameters() []Parameter
	GetExtensions() map[string]any
}

// Operation abstracts an operation across OpenAPI versions
type Operation interface {
	IsNil() bool
	GetOperationID() string
	GetSummary() string
	GetDescription() string
	GetParameters() []Parameter
	GetRequestBody() RequestBody
	GetResponses() Responses
	GetSecurity() []SecurityRequirement
	GetTags() []string
	GetExtensions() map[string]any
	GetExternalDocs() ExternalDocumentation
	GetDeprecated() bool
}

// Parameter abstracts a parameter across OpenAPI versions
type Parameter interface {
	GetName() string
	GetIn() string
	GetRequired() bool
	GetDescription() string
	GetSchema() Schema
	// For Swagger 2.0 body parameters
	IsBodyParameter() bool
	GetExtensions() map[string]any
	GetDeprecated() bool
	GetAllowEmptyValue() bool
	GetStyle() string
	GetExplode() bool
	GetAllowReserved() bool
}

// RequestBody abstracts a request body across OpenAPI versions
type RequestBody interface {
	IsNil() bool
	GetRequired() bool
	GetContent() map[string]MediaType
	GetDescription() string
	GetExtensions() map[string]any
}

// MediaType abstracts media type content
type MediaType interface {
	GetSchema() Schema
	GetExtensions() map[string]any
}

// Responses abstracts the responses object
type Responses interface {
	GetDefault() Response
	GetStatusCodes() map[string]Response
	GetExtensions() map[string]any
}

// Response abstracts a single response
type Response interface {
	IsNil() bool
	HasRef() bool
	GetRef() string
	GetDescription() string
	GetHeaders() map[string]Header
	GetContent() map[string]MediaType
	// For Swagger 2.0 which has schema directly on response
	GetSchema() Schema
	GetExtensions() map[string]any
}

// Header abstracts a response header
type Header interface {
	GetSchema() Schema
	GetRequired() bool
	GetDescription() string
	GetExtensions() map[string]any
	GetDeprecated() bool
}

// Schema abstracts a schema across OpenAPI versions
type Schema interface {
	IsNil() bool
	GetRef() string
	GetType() string
	GetFormat() string
	GetDescription() string
	GetProperties() map[string]Schema
	GetItems() Schema
	GetRequired() []string
	GetAllOf() []Schema
	GetOneOf() []Schema
	GetAnyOf() []Schema
	// For boolean schemas (additionalProperties: false)
	IsBooleanSchema() bool
	GetBooleanValue() *bool
	GetExtensions() map[string]any

	GetDiscriminator() Discriminator
	GetXML() XML
	GetExternalDocs() ExternalDocumentation
	GetExample() any
	GetDefault() any
	GetDeprecated() bool
	GetReadOnly() bool
	GetWriteOnly() bool
}

// SecurityScheme abstracts a security scheme across OpenAPI versions
type SecurityScheme interface {
	GetType() string   // apiKey, http, oauth2, openIdConnect (3.0+), basic (2.0)
	GetName() string   // For apiKey
	GetIn() string     // For apiKey: query, header, cookie
	GetScheme() string // For http: bearer, basic, etc.
	GetDescription() string
	// For OAuth2
	GetFlow() string // For 2.0: implicit, password, application, accessCode; For 3.0+: implicit, password, clientCredentials, authorizationCode
	GetAuthorizationURL() string
	GetTokenURL() string
	GetScopes() map[string]string
	GetExtensions() map[string]any
}

// SecurityRequirement is a map of security scheme names to required scopes
type SecurityRequirement = map[string][]string

// Discriminator abstracts polymorphism discriminator
type Discriminator interface {
	GetPropertyName() string
	GetMapping() map[string]string
}

// XML abstracts XML metadata
type XML interface {
	GetName() string
	GetNamespace() string
	GetPrefix() string
	GetAttribute() bool
	GetWrapped() bool
}

// ExternalDocumentation abstracts external documentation
type ExternalDocumentation interface {
	GetDescription() string
	GetURL() string
}

// NilRequestBody is returned when there is no request body
type NilRequestBody struct{}

func (n NilRequestBody) IsNil() bool                      { return true }
func (n NilRequestBody) GetRequired() bool                { return false }
func (n NilRequestBody) GetContent() map[string]MediaType { return nil }
func (n NilRequestBody) GetDescription() string           { return "" }
func (n NilRequestBody) GetExtensions() map[string]any    { return nil }

// NilResponse is returned when there is no response
type NilResponse struct{}

func (n NilResponse) IsNil() bool                      { return true }
func (n NilResponse) HasRef() bool                     { return false }
func (n NilResponse) GetRef() string                   { return "" }
func (n NilResponse) GetDescription() string           { return "" }
func (n NilResponse) GetHeaders() map[string]Header    { return nil }
func (n NilResponse) GetContent() map[string]MediaType { return nil }
func (n NilResponse) GetSchema() Schema                { return nil }
func (n NilResponse) GetExtensions() map[string]any    { return nil }

// NilSchema is returned when there is no schema
type NilSchema struct{}

func (n NilSchema) IsNil() bool                            { return true }
func (n NilSchema) GetRef() string                         { return "" }
func (n NilSchema) GetType() string                        { return "" }
func (n NilSchema) GetFormat() string                      { return "" }
func (n NilSchema) GetDescription() string                 { return "" }
func (n NilSchema) GetProperties() map[string]Schema       { return nil }
func (n NilSchema) GetItems() Schema                       { return nil }
func (n NilSchema) GetRequired() []string                  { return nil }
func (n NilSchema) GetAllOf() []Schema                     { return nil }
func (n NilSchema) GetOneOf() []Schema                     { return nil }
func (n NilSchema) GetAnyOf() []Schema                     { return nil }
func (n NilSchema) IsBooleanSchema() bool                  { return false }
func (n NilSchema) GetBooleanValue() *bool                 { return nil }
func (n NilSchema) GetExtensions() map[string]any          { return nil }
func (n NilSchema) GetDiscriminator() Discriminator        { return nil }
func (n NilSchema) GetXML() XML                            { return nil }
func (n NilSchema) GetExternalDocs() ExternalDocumentation { return nil }
func (n NilSchema) GetExample() any                        { return nil }
func (n NilSchema) GetDefault() any                        { return nil }
func (n NilSchema) GetDeprecated() bool                    { return false }
func (n NilSchema) GetReadOnly() bool                      { return false }
func (n NilSchema) GetWriteOnly() bool                     { return false }

// NilOperation is returned when there is no operation
type NilOperation struct{}

func (n NilOperation) IsNil() bool                            { return true }
func (n NilOperation) GetOperationID() string                 { return "" }
func (n NilOperation) GetSummary() string                     { return "" }
func (n NilOperation) GetDescription() string                 { return "" }
func (n NilOperation) GetParameters() []Parameter             { return nil }
func (n NilOperation) GetRequestBody() RequestBody            { return NilRequestBody{} }
func (n NilOperation) GetResponses() Responses                { return nil }
func (n NilOperation) GetSecurity() []SecurityRequirement     { return nil }
func (n NilOperation) GetTags() []string                      { return nil }
func (n NilOperation) GetExtensions() map[string]any          { return nil }
func (n NilOperation) GetExternalDocs() ExternalDocumentation { return nil }
func (n NilOperation) GetDeprecated() bool                    { return false }

// NilResponses is returned when there are no responses
type NilResponses struct{}

func (n NilResponses) GetDefault() Response                { return NilResponse{} }
func (n NilResponses) GetStatusCodes() map[string]Response { return nil }
func (n NilResponses) GetExtensions() map[string]any       { return nil }

// BaseDiscriminator provides a default implementation for Discriminator
type BaseDiscriminator struct {
	PropertyName string
	Mapping      map[string]string
}

func (d BaseDiscriminator) GetPropertyName() string       { return d.PropertyName }
func (d BaseDiscriminator) GetMapping() map[string]string { return d.Mapping }

// BaseXML provides a default implementation for XML
type BaseXML struct {
	Name      string
	Namespace string
	Prefix    string
	Attribute bool
	Wrapped   bool
}

func (x BaseXML) GetName() string      { return x.Name }
func (x BaseXML) GetNamespace() string { return x.Namespace }
func (x BaseXML) GetPrefix() string    { return x.Prefix }
func (x BaseXML) GetAttribute() bool   { return x.Attribute }
func (x BaseXML) GetWrapped() bool     { return x.Wrapped }

// BaseExternalDocs provides a default implementation for ExternalDocumentation
type BaseExternalDocs struct {
	Description string
	URL         string
}

func (e BaseExternalDocs) GetDescription() string { return e.Description }
func (e BaseExternalDocs) GetURL() string         { return e.URL }
