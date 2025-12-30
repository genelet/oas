// Package unified provides an adapter for OpenAPI 2.0 (Swagger) documents
// Copyright (c) Genelet

package unified

import (
	"strings"

	oa2 "github.com/genelet/oas/openapi20"
)

// Document20 wraps an OpenAPI 2.0 document and implements Document
type Document20 struct {
	doc *oa2.Swagger
}

// NewDocument20 creates a new Document adapter for OpenAPI 2.0
func NewDocument20(doc *oa2.Swagger) Document {
	return &Document20{doc: doc}
}

// GetRaw returns the underlying OpenAPI 2.0 document
func (d *Document20) GetRaw() *oa2.Swagger {
	return d.doc
}

func (d *Document20) Version() string {
	return d.doc.Swagger
}

func (d *Document20) GetServerURL() string {
	// Construct from host, basePath, schemes
	scheme := "https"
	if len(d.doc.Schemes) > 0 {
		scheme = d.doc.Schemes[0]
	}
	host := d.doc.Host
	basePath := d.doc.BasePath
	if host == "" {
		return basePath
	}
	return scheme + "://" + host + basePath
}

func (d *Document20) GetInfo() DocumentInfo {
	if d.doc.Info == nil {
		return &documentInfo20{}
	}
	return &documentInfo20{info: d.doc.Info}
}

func (d *Document20) GetPaths() map[string]PathItem {
	if d.doc.Paths == nil || d.doc.Paths.Paths == nil {
		return nil
	}
	result := make(map[string]PathItem)
	for path, item := range d.doc.Paths.Paths {
		if item != nil {
			result[path] = &pathItem20{item: item, doc: d.doc}
		}
	}
	return result
}

func (d *Document20) GetSecuritySchemes() map[string]SecurityScheme {
	if d.doc.SecurityDefinitions == nil {
		return nil
	}
	result := make(map[string]SecurityScheme)
	for name, scheme := range d.doc.SecurityDefinitions {
		if scheme != nil {
			result[name] = &securityScheme20{scheme: scheme}
		}
	}
	return result
}

func (d *Document20) GetGlobalSecurity() []SecurityRequirement {
	if d.doc.Security == nil {
		return nil
	}
	result := make([]SecurityRequirement, len(d.doc.Security))
	for i, sec := range d.doc.Security {
		result[i] = SecurityRequirement(sec)
	}
	return result
}

func (d *Document20) GetExtensions() map[string]any {
	return d.doc.Extensions
}

// documentInfo20 wraps OpenAPI 2.0 Info
type documentInfo20 struct {
	info *oa2.Info
}

func (i *documentInfo20) GetTitle() string {
	if i.info == nil {
		return ""
	}
	return i.info.Title
}

func (i *documentInfo20) GetVersion() string {
	if i.info == nil {
		return ""
	}
	return i.info.Version
}

func (i *documentInfo20) GetDescription() string {
	if i.info == nil {
		return ""
	}
	return i.info.Description
}

func (i *documentInfo20) GetExtensions() map[string]any {
	if i.info == nil {
		return nil
	}
	return i.info.Extensions
}

// pathItem20 wraps OpenAPI 2.0 PathItem
type pathItem20 struct {
	item *oa2.PathItem
	doc  *oa2.Swagger
}

func (p *pathItem20) HasRef() bool {
	return p.item != nil && p.item.HasRef()
}

func (p *pathItem20) GetRef() string {
	if p.item == nil {
		return ""
	}
	return p.item.GetRef()
}

func (p *pathItem20) GetOperation(method string) Operation {
	if p.item == nil {
		return NilOperation{}
	}
	method = strings.ToLower(method)
	var op *oa2.Operation
	switch method {
	case "get":
		op = p.item.Get
	case "put":
		op = p.item.Put
	case "post":
		op = p.item.Post
	case "delete":
		op = p.item.Delete
	case "options":
		op = p.item.Options
	case "head":
		op = p.item.Head
	case "patch":
		op = p.item.Patch
	}
	if op == nil {
		return NilOperation{}
	}
	return &operation20{op: op, doc: p.doc}
}

func (p *pathItem20) GetAllOperations() map[string]Operation {
	if p.item == nil {
		return nil
	}
	result := make(map[string]Operation)
	if p.item.Get != nil {
		result["get"] = &operation20{op: p.item.Get, doc: p.doc}
	}
	if p.item.Put != nil {
		result["put"] = &operation20{op: p.item.Put, doc: p.doc}
	}
	if p.item.Post != nil {
		result["post"] = &operation20{op: p.item.Post, doc: p.doc}
	}
	if p.item.Delete != nil {
		result["delete"] = &operation20{op: p.item.Delete, doc: p.doc}
	}
	if p.item.Options != nil {
		result["options"] = &operation20{op: p.item.Options, doc: p.doc}
	}
	if p.item.Head != nil {
		result["head"] = &operation20{op: p.item.Head, doc: p.doc}
	}
	if p.item.Patch != nil {
		result["patch"] = &operation20{op: p.item.Patch, doc: p.doc}
	}
	return result
}

func (p *pathItem20) GetParameters() []Parameter {
	if p.item == nil || p.item.Parameters == nil {
		return nil
	}
	result := make([]Parameter, 0, len(p.item.Parameters))
	for _, param := range p.item.Parameters {
		if param != nil {
			result = append(result, &parameter20{param: param})
		}
	}
	return result
}

func (p *pathItem20) GetExtensions() map[string]any {
	if p.item == nil {
		return nil
	}
	return p.item.Extensions
}

// operation20 wraps OpenAPI 2.0 Operation
type operation20 struct {
	op  *oa2.Operation
	doc *oa2.Swagger
}

func (o *operation20) IsNil() bool {
	return o.op == nil
}

func (o *operation20) GetOperationID() string {
	if o.op == nil {
		return ""
	}
	return o.op.OperationID
}

func (o *operation20) GetSummary() string {
	if o.op == nil {
		return ""
	}
	return o.op.Summary
}

func (o *operation20) GetDescription() string {
	if o.op == nil {
		return ""
	}
	return o.op.Description
}

func (o *operation20) GetParameters() []Parameter {
	if o.op == nil || o.op.Parameters == nil {
		return nil
	}
	result := make([]Parameter, 0, len(o.op.Parameters))
	for _, param := range o.op.Parameters {
		if param != nil && !param.IsBodyParameter() {
			result = append(result, &parameter20{param: param})
		}
	}
	return result
}

func (o *operation20) GetRequestBody() RequestBody {
	if o.op == nil || o.op.Parameters == nil {
		return NilRequestBody{}
	}
	// In Swagger 2.0, request body is a parameter with in=body
	for _, param := range o.op.Parameters {
		if param != nil && param.IsBodyParameter() {
			return &requestBody20{param: param, doc: o.doc}
		}
	}
	return NilRequestBody{}
}

func (o *operation20) GetResponses() Responses {
	if o.op == nil || o.op.Responses == nil {
		return NilResponses{}
	}
	return &responses20{responses: o.op.Responses}
}

func (o *operation20) GetSecurity() []SecurityRequirement {
	if o.op == nil || o.op.Security == nil {
		return nil
	}
	result := make([]SecurityRequirement, len(o.op.Security))
	for i, sec := range o.op.Security {
		result[i] = SecurityRequirement(sec)
	}
	return result
}

func (o *operation20) GetTags() []string {
	if o.op == nil {
		return nil
	}
	return o.op.Tags
}

func (o *operation20) GetExtensions() map[string]any {
	if o.op == nil {
		return nil
	}
	return o.op.Extensions
}

func (o *operation20) GetExternalDocs() ExternalDocumentation {
	if o.op == nil || o.op.ExternalDocs == nil {
		return nil
	}
	return BaseExternalDocs{
		Description: o.op.ExternalDocs.Description,
		URL:         o.op.ExternalDocs.URL,
	}
}

func (o *operation20) GetDeprecated() bool {
	if o.op == nil {
		return false
	}
	return o.op.Deprecated
}

// parameter20 wraps OpenAPI 2.0 Parameter (non-body)
type parameter20 struct {
	param *oa2.Parameter
}

func (p *parameter20) GetName() string {
	if p.param == nil {
		return ""
	}
	return p.param.Name
}

func (p *parameter20) GetIn() string {
	if p.param == nil {
		return ""
	}
	return p.param.In
}

func (p *parameter20) GetRequired() bool {
	if p.param == nil {
		return false
	}
	return p.param.Required
}

func (p *parameter20) GetDescription() string {
	if p.param == nil {
		return ""
	}
	return p.param.Description
}

func (p *parameter20) GetSchema() Schema {
	if p.param == nil {
		return NilSchema{}
	}
	// In Swagger 2.0, non-body parameters have type/format directly, not schema
	// We create a synthetic schema from these fields
	if p.param.Type != "" {
		return &parameterSchema20{param: p.param}
	}
	return NilSchema{}
}

func (p *parameter20) IsBodyParameter() bool {
	return p.param != nil && p.param.In == "body"
}

func (p *parameter20) GetExtensions() map[string]any {
	if p.param == nil {
		return nil
	}
	return p.param.Extensions
}

func (p *parameter20) GetDeprecated() bool {
	// Swagger 2.0 parameters don't have deprecated field
	return false
}

func (p *parameter20) GetAllowEmptyValue() bool {
	if p.param == nil {
		return false
	}
	return p.param.AllowEmptyValue
}

func (p *parameter20) GetStyle() string {
	// Swagger 2.0 collectionFormat -> style adaptation could be done, but keeping simple for now
	if p.param == nil {
		return ""
	}
	// Basic mapping
	switch p.param.CollectionFormat {
	case "csv":
		return "form" // roughly equivalent
	case "ssv":
		return "spaceDelimited"
	case "tsv":
		return "tabDelimited"
	case "pipes":
		return "pipeDelimited"
	case "multi":
		return "form" // with explode=true
	}
	return ""
}

func (p *parameter20) GetExplode() bool {
	if p.param == nil {
		return false
	}
	// In Swagger 2.0, collectionFormat=multi implies explode=true
	return p.param.CollectionFormat == "multi"
}

func (p *parameter20) GetAllowReserved() bool {
	return false // Not supported in Swagger 2.0
}

// parameterSchema20 creates a schema from parameter fields
type parameterSchema20 struct {
	param *oa2.Parameter
}

func (s *parameterSchema20) IsNil() bool                      { return s.param == nil }
func (s *parameterSchema20) GetRef() string                   { return "" }
func (s *parameterSchema20) GetType() string                  { return s.param.Type }
func (s *parameterSchema20) GetFormat() string                { return s.param.Format }
func (s *parameterSchema20) GetDescription() string           { return s.param.Description }
func (s *parameterSchema20) GetProperties() map[string]Schema { return nil }
func (s *parameterSchema20) GetItems() Schema {
	if s.param.Items == nil {
		return NilSchema{}
	}
	return &itemsSchema20{items: s.param.Items}
}
func (s *parameterSchema20) GetRequired() []string  { return nil }
func (s *parameterSchema20) GetAllOf() []Schema     { return nil }
func (s *parameterSchema20) GetOneOf() []Schema     { return nil }
func (s *parameterSchema20) GetAnyOf() []Schema     { return nil }
func (s *parameterSchema20) IsBooleanSchema() bool  { return false }
func (s *parameterSchema20) GetBooleanValue() *bool { return nil }
func (s *parameterSchema20) GetExtensions() map[string]any {
	if s.param == nil {
		return nil
	}
	return s.param.Extensions
}
func (s *parameterSchema20) GetDiscriminator() Discriminator        { return nil }
func (s *parameterSchema20) GetXML() XML                            { return nil }
func (s *parameterSchema20) GetExternalDocs() ExternalDocumentation { return nil }
func (s *parameterSchema20) GetExample() any                        { return nil }
func (s *parameterSchema20) GetDefault() any {
	if s.param == nil {
		return nil
	}
	return s.param.Default
}
func (s *parameterSchema20) GetDeprecated() bool { return false }
func (s *parameterSchema20) GetReadOnly() bool   { return false }
func (s *parameterSchema20) GetWriteOnly() bool  { return false }

// itemsSchema20 wraps Items as a Schema
type itemsSchema20 struct {
	items *oa2.Items
}

func (s *itemsSchema20) IsNil() bool                      { return s.items == nil }
func (s *itemsSchema20) GetRef() string                   { return "" }
func (s *itemsSchema20) GetType() string                  { return s.items.Type }
func (s *itemsSchema20) GetFormat() string                { return s.items.Format }
func (s *itemsSchema20) GetDescription() string           { return "" }
func (s *itemsSchema20) GetProperties() map[string]Schema { return nil }
func (s *itemsSchema20) GetItems() Schema {
	if s.items.Items == nil {
		return NilSchema{}
	}
	return &itemsSchema20{items: s.items.Items}
}
func (s *itemsSchema20) GetRequired() []string  { return nil }
func (s *itemsSchema20) GetAllOf() []Schema     { return nil }
func (s *itemsSchema20) GetOneOf() []Schema     { return nil }
func (s *itemsSchema20) GetAnyOf() []Schema     { return nil }
func (s *itemsSchema20) IsBooleanSchema() bool  { return false }
func (s *itemsSchema20) GetBooleanValue() *bool { return nil }
func (s *itemsSchema20) GetExtensions() map[string]any {
	// Items in Swagger 2.0 don't formally support extensions but some parsers might add them
	return nil
}
func (s *itemsSchema20) GetDiscriminator() Discriminator        { return nil }
func (s *itemsSchema20) GetXML() XML                            { return nil }
func (s *itemsSchema20) GetExternalDocs() ExternalDocumentation { return nil }
func (s *itemsSchema20) GetExample() any                        { return nil }
func (s *itemsSchema20) GetDefault() any {
	if s.items == nil {
		return nil
	}
	return s.items.Default
}
func (s *itemsSchema20) GetDeprecated() bool { return false }
func (s *itemsSchema20) GetReadOnly() bool   { return false }
func (s *itemsSchema20) GetWriteOnly() bool  { return false }

// requestBody20 wraps a body parameter as RequestBody
type requestBody20 struct {
	param *oa2.Parameter
	doc   *oa2.Swagger
}

func (r *requestBody20) IsNil() bool {
	return r.param == nil
}

func (r *requestBody20) GetRequired() bool {
	if r.param == nil {
		return false
	}
	return r.param.Required
}

func (r *requestBody20) GetContent() map[string]MediaType {
	if r.param == nil || r.param.Schema == nil {
		return nil
	}
	// In Swagger 2.0, content types come from consumes
	// Default to application/json if not specified
	consumes := []string{"application/json"}
	if r.doc != nil && len(r.doc.Consumes) > 0 {
		consumes = r.doc.Consumes
	}
	result := make(map[string]MediaType)
	for _, ct := range consumes {
		result[ct] = &mediaType20{schema: r.param.Schema}
	}
	return result
}

func (r *requestBody20) GetDescription() string {
	if r.param == nil {
		return ""
	}
	return r.param.Description
}

func (r *requestBody20) GetExtensions() map[string]any {
	if r.param == nil {
		return nil
	}
	return r.param.Extensions
}

// mediaType20 wraps a schema as MediaType
type mediaType20 struct {
	schema *oa2.Schema
}

func (m *mediaType20) GetSchema() Schema {
	if m.schema == nil {
		return NilSchema{}
	}
	return &schema20{schema: m.schema}
}

func (m *mediaType20) GetExtensions() map[string]any {
	// MediaType doesn't exist in 2.0, so no extensions
	return nil
}

// responses20 wraps OpenAPI 2.0 Responses
type responses20 struct {
	responses *oa2.Responses
}

func (r *responses20) GetDefault() Response {
	if r.responses == nil || r.responses.Default == nil {
		return NilResponse{}
	}
	return &response20{resp: r.responses.Default}
}

func (r *responses20) GetStatusCodes() map[string]Response {
	if r.responses == nil || r.responses.StatusCode == nil {
		return nil
	}
	result := make(map[string]Response)
	for code, resp := range r.responses.StatusCode {
		if resp != nil {
			result[code] = &response20{resp: resp}
		}
	}
	return result
}

func (r *responses20) GetExtensions() map[string]any {
	if r.responses == nil {
		return nil
	}
	return r.responses.Extensions
}

// response20 wraps OpenAPI 2.0 Response
type response20 struct {
	resp *oa2.Response
}

func (r *response20) IsNil() bool {
	return r.resp == nil
}

func (r *response20) HasRef() bool {
	return r.resp != nil && r.resp.IsReference()
}

func (r *response20) GetRef() string {
	if r.resp == nil {
		return ""
	}
	return r.resp.Ref
}

func (r *response20) GetDescription() string {
	if r.resp == nil {
		return ""
	}
	return r.resp.Description
}

func (r *response20) GetHeaders() map[string]Header {
	if r.resp == nil || r.resp.Headers == nil {
		return nil
	}
	result := make(map[string]Header)
	for name, header := range r.resp.Headers {
		if header != nil {
			result[name] = &header20{header: header}
		}
	}
	return result
}

func (r *response20) GetContent() map[string]MediaType {
	// In Swagger 2.0, response has schema directly, not content
	// We still return nil since GetSchema handles this
	return nil
}

func (r *response20) GetSchema() Schema {
	if r.resp == nil || r.resp.Schema == nil {
		return NilSchema{}
	}
	return &schema20{schema: r.resp.Schema}
}

func (r *response20) GetExtensions() map[string]any {
	if r.resp == nil {
		return nil
	}
	return r.resp.Extensions
}

// header20 wraps OpenAPI 2.0 Header
type header20 struct {
	header *oa2.Header
}

func (h *header20) GetSchema() Schema {
	if h.header == nil {
		return NilSchema{}
	}
	// Headers in 2.0 have type/format directly, not schema
	return &headerSchema20{header: h.header}
}

func (h *header20) GetRequired() bool {
	// Swagger 2.0 headers don't have required field
	return false
}

func (h *header20) GetDescription() string {
	if h.header == nil {
		return ""
	}
	return h.header.Description
}

func (h *header20) GetExtensions() map[string]any {
	// Swagger 2.0 headers don't strictly support extensions like 3.0
	// But our struct might have them if parser supports it.
	// oa2.Header doesn't have Extensions field usually, checking...
	// Assuming oa2.Header might NOT have Extensions based on standard 2.0
	return nil
}

func (h *header20) GetDeprecated() bool {
	return false
}

// headerSchema20 creates a schema from header fields
type headerSchema20 struct {
	header *oa2.Header
}

func (s *headerSchema20) IsNil() bool                      { return s.header == nil }
func (s *headerSchema20) GetRef() string                   { return "" }
func (s *headerSchema20) GetType() string                  { return s.header.Type }
func (s *headerSchema20) GetFormat() string                { return s.header.Format }
func (s *headerSchema20) GetDescription() string           { return s.header.Description }
func (s *headerSchema20) GetProperties() map[string]Schema { return nil }
func (s *headerSchema20) GetItems() Schema {
	if s.header.Items == nil {
		return NilSchema{}
	}
	return &itemsSchema20{items: s.header.Items}
}
func (s *headerSchema20) GetRequired() []string                  { return nil }
func (s *headerSchema20) GetAllOf() []Schema                     { return nil }
func (s *headerSchema20) GetOneOf() []Schema                     { return nil }
func (s *headerSchema20) GetAnyOf() []Schema                     { return nil }
func (s *headerSchema20) IsBooleanSchema() bool                  { return false }
func (s *headerSchema20) GetBooleanValue() *bool                 { return nil }
func (s *headerSchema20) GetExtensions() map[string]any          { return nil }
func (s *headerSchema20) GetDiscriminator() Discriminator        { return nil }
func (s *headerSchema20) GetXML() XML                            { return nil }
func (s *headerSchema20) GetExternalDocs() ExternalDocumentation { return nil }
func (s *headerSchema20) GetExample() any                        { return nil }
func (s *headerSchema20) GetDefault() any {
	if s.header == nil {
		return nil
	}
	return s.header.Default
}
func (s *headerSchema20) GetDeprecated() bool { return false }
func (s *headerSchema20) GetReadOnly() bool   { return false }
func (s *headerSchema20) GetWriteOnly() bool  { return false }

// schema20 wraps OpenAPI 2.0 Schema
type schema20 struct {
	schema *oa2.Schema
}

func (s *schema20) IsNil() bool {
	return s.schema == nil
}

func (s *schema20) GetRef() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Ref
}

func (s *schema20) GetType() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Type
}

func (s *schema20) GetFormat() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Format
}

func (s *schema20) GetDescription() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Description
}

func (s *schema20) GetProperties() map[string]Schema {
	if s.schema == nil || s.schema.Properties == nil {
		return nil
	}
	result := make(map[string]Schema)
	for name, prop := range s.schema.Properties {
		if prop != nil {
			result[name] = &schema20{schema: prop}
		}
	}
	return result
}

func (s *schema20) GetItems() Schema {
	if s.schema == nil || s.schema.Items == nil {
		return NilSchema{}
	}
	return &schema20{schema: s.schema.Items}
}

func (s *schema20) GetRequired() []string {
	if s.schema == nil {
		return nil
	}
	return s.schema.Required
}

func (s *schema20) GetAllOf() []Schema {
	if s.schema == nil || s.schema.AllOf == nil {
		return nil
	}
	result := make([]Schema, 0, len(s.schema.AllOf))
	for _, sub := range s.schema.AllOf {
		if sub != nil {
			result = append(result, &schema20{schema: sub})
		}
	}
	return result
}

func (s *schema20) GetOneOf() []Schema {
	// Swagger 2.0 doesn't support oneOf
	return nil
}

func (s *schema20) GetAnyOf() []Schema {
	// Swagger 2.0 doesn't support anyOf
	return nil
}

func (s *schema20) IsBooleanSchema() bool {
	if s.schema == nil {
		return false
	}
	return s.schema.IsBooleanSchema()
}

func (s *schema20) GetBooleanValue() *bool {
	if s.schema == nil {
		return nil
	}
	return s.schema.BooleanValue()
}

func (s *schema20) GetExtensions() map[string]any {
	if s.schema == nil {
		return nil
	}
	return s.schema.Extensions
}

func (s *schema20) GetDiscriminator() Discriminator {
	if s.schema == nil {
		// Swagger 2.0 discriminator is just a string (property name)
		if s.schema.Discriminator == "" {
			return nil
		}
		return BaseDiscriminator{
			PropertyName: s.schema.Discriminator,
			// No mapping in 2.0
		}
	}
	return nil
}

func (s *schema20) GetXML() XML {
	if s.schema == nil || s.schema.XML == nil {
		return nil
	}
	return BaseXML{
		Name:      s.schema.XML.Name,
		Namespace: s.schema.XML.Namespace,
		Prefix:    s.schema.XML.Prefix,
		Attribute: s.schema.XML.Attribute,
		Wrapped:   s.schema.XML.Wrapped,
	}
}

func (s *schema20) GetExternalDocs() ExternalDocumentation {
	if s.schema == nil || s.schema.ExternalDocs == nil {
		return nil
	}
	return BaseExternalDocs{
		Description: s.schema.ExternalDocs.Description,
		URL:         s.schema.ExternalDocs.URL,
	}
}

func (s *schema20) GetExample() any {
	if s.schema == nil {
		return nil
	}
	return s.schema.Example
}

func (s *schema20) GetDefault() any {
	if s.schema == nil {
		return nil
	}
	return s.schema.Default
}

func (s *schema20) GetDeprecated() bool {
	// Not standard in Swagger 2.0 schema, but often supported
	// oa2.Schema might not have it. Checked: likely not.
	return false
}

func (s *schema20) GetReadOnly() bool {
	if s.schema == nil {
		return false
	}
	return s.schema.ReadOnly
}

func (s *schema20) GetWriteOnly() bool {
	// Not in Swagger 2.0
	return false
}

// securityScheme20 wraps OpenAPI 2.0 SecurityScheme
type securityScheme20 struct {
	scheme *oa2.SecurityScheme
}

func (s *securityScheme20) GetType() string {
	if s.scheme == nil {
		return ""
	}
	// Convert Swagger 2.0 "basic" to OpenAPI 3.0 "http"
	if s.scheme.Type == "basic" {
		return "http"
	}
	return s.scheme.Type
}

func (s *securityScheme20) GetName() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Name
}

func (s *securityScheme20) GetIn() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.In
}

func (s *securityScheme20) GetScheme() string {
	if s.scheme == nil {
		return ""
	}
	// For Swagger 2.0 basic auth, return "basic" as the scheme
	if s.scheme.Type == "basic" {
		return "basic"
	}
	return ""
}

func (s *securityScheme20) GetDescription() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Description
}

func (s *securityScheme20) GetFlow() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Flow
}

func (s *securityScheme20) GetAuthorizationURL() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.AuthorizationUrl
}

func (s *securityScheme20) GetTokenURL() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.TokenUrl
}

func (s *securityScheme20) GetScopes() map[string]string {
	if s.scheme == nil {
		return nil
	}
	return s.scheme.Scopes
}

func (s *securityScheme20) GetExtensions() map[string]any {
	if s.scheme == nil {
		return nil
	}
	return s.scheme.Extensions
}
