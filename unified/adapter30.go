// Package unified provides an adapter for OpenAPI 3.0 documents
// Copyright (c) Greetingland LLC

package unified

import (
	"strings"

	oa3 "github.com/genelet/oas/openapi30"
)

// Document30 wraps an OpenAPI 3.0 document and implements Document
type Document30 struct {
	doc *oa3.OpenAPI
}

// NewDocument30 creates a new Document adapter for OpenAPI 3.0
func NewDocument30(doc *oa3.OpenAPI) Document {
	return &Document30{doc: doc}
}

// GetRaw returns the underlying OpenAPI 3.0 document
func (d *Document30) GetRaw() *oa3.OpenAPI {
	return d.doc
}

func (d *Document30) Version() string {
	return d.doc.OpenAPI
}

func (d *Document30) GetServerURL() string {
	if len(d.doc.Servers) > 0 && d.doc.Servers[0] != nil {
		return d.doc.Servers[0].URL
	}
	return ""
}

func (d *Document30) GetInfo() DocumentInfo {
	if d.doc.Info == nil {
		return &documentInfo30{}
	}
	return &documentInfo30{info: d.doc.Info}
}

func (d *Document30) GetPaths() map[string]PathItem {
	if d.doc.Paths == nil || d.doc.Paths.Paths == nil {
		return nil
	}
	result := make(map[string]PathItem)
	for path, item := range d.doc.Paths.Paths {
		if item != nil {
			result[path] = &pathItem30{item: item}
		}
	}
	return result
}

func (d *Document30) GetSecuritySchemes() map[string]SecurityScheme {
	if d.doc.Components == nil || d.doc.Components.SecuritySchemes == nil {
		return nil
	}
	result := make(map[string]SecurityScheme)
	for name, scheme := range d.doc.Components.SecuritySchemes {
		if scheme != nil {
			result[name] = &securityScheme30{scheme: scheme}
		}
	}
	return result
}

func (d *Document30) GetGlobalSecurity() []SecurityRequirement {
	if d.doc.Security == nil {
		return nil
	}
	result := make([]SecurityRequirement, len(d.doc.Security))
	for i, sec := range d.doc.Security {
		result[i] = SecurityRequirement(sec)
	}
	return result
}

func (d *Document30) GetExtensions() map[string]any {
	return d.doc.Extensions
}

// documentInfo30 wraps OpenAPI 3.0 Info
type documentInfo30 struct {
	info *oa3.Info
}

func (i *documentInfo30) GetTitle() string {
	if i.info == nil {
		return ""
	}
	return i.info.Title
}

func (i *documentInfo30) GetVersion() string {
	if i.info == nil {
		return ""
	}
	return i.info.Version
}

func (i *documentInfo30) GetDescription() string {
	if i.info == nil {
		return ""
	}
	return i.info.Description
}

func (i *documentInfo30) GetExtensions() map[string]any {
	if i.info == nil {
		return nil
	}
	return i.info.Extensions
}

// pathItem30 wraps OpenAPI 3.0 PathItem
type pathItem30 struct {
	item *oa3.PathItem
}

func (p *pathItem30) HasRef() bool {
	return p.item != nil && p.item.HasRef()
}

func (p *pathItem30) GetRef() string {
	if p.item == nil {
		return ""
	}
	return p.item.GetRef()
}

func (p *pathItem30) GetOperation(method string) Operation {
	if p.item == nil {
		return NilOperation{}
	}
	method = strings.ToLower(method)
	var op *oa3.Operation
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
	case "trace":
		op = p.item.Trace
	}
	if op == nil {
		return NilOperation{}
	}
	return &operation30{op: op}
}

func (p *pathItem30) GetAllOperations() map[string]Operation {
	if p.item == nil {
		return nil
	}
	result := make(map[string]Operation)
	if p.item.Get != nil {
		result["get"] = &operation30{op: p.item.Get}
	}
	if p.item.Put != nil {
		result["put"] = &operation30{op: p.item.Put}
	}
	if p.item.Post != nil {
		result["post"] = &operation30{op: p.item.Post}
	}
	if p.item.Delete != nil {
		result["delete"] = &operation30{op: p.item.Delete}
	}
	if p.item.Options != nil {
		result["options"] = &operation30{op: p.item.Options}
	}
	if p.item.Head != nil {
		result["head"] = &operation30{op: p.item.Head}
	}
	if p.item.Patch != nil {
		result["patch"] = &operation30{op: p.item.Patch}
	}
	if p.item.Trace != nil {
		result["trace"] = &operation30{op: p.item.Trace}
	}
	return result
}

func (p *pathItem30) GetParameters() []Parameter {
	if p.item == nil || p.item.Parameters == nil {
		return nil
	}
	result := make([]Parameter, 0, len(p.item.Parameters))
	for _, param := range p.item.Parameters {
		if param != nil {
			result = append(result, &parameter30{param: param})
		}
	}
	return result
}

func (p *pathItem30) GetExtensions() map[string]any {
	if p.item == nil {
		return nil
	}
	return p.item.Extensions
}

// operation30 wraps OpenAPI 3.0 Operation
type operation30 struct {
	op *oa3.Operation
}

func (o *operation30) IsNil() bool {
	return o.op == nil
}

func (o *operation30) GetOperationID() string {
	if o.op == nil {
		return ""
	}
	return o.op.OperationID
}

func (o *operation30) GetSummary() string {
	if o.op == nil {
		return ""
	}
	return o.op.Summary
}

func (o *operation30) GetDescription() string {
	if o.op == nil {
		return ""
	}
	return o.op.Description
}

func (o *operation30) GetParameters() []Parameter {
	if o.op == nil || o.op.Parameters == nil {
		return nil
	}
	result := make([]Parameter, 0, len(o.op.Parameters))
	for _, param := range o.op.Parameters {
		if param != nil {
			result = append(result, &parameter30{param: param})
		}
	}
	return result
}

func (o *operation30) GetRequestBody() RequestBody {
	if o.op == nil || o.op.RequestBody == nil {
		return NilRequestBody{}
	}
	return &requestBody30{rb: o.op.RequestBody}
}

func (o *operation30) GetResponses() Responses {
	if o.op == nil || o.op.Responses == nil {
		return NilResponses{}
	}
	return &responses30{responses: o.op.Responses}
}

func (o *operation30) GetSecurity() []SecurityRequirement {
	if o.op == nil || o.op.Security == nil {
		return nil
	}
	result := make([]SecurityRequirement, len(o.op.Security))
	for i, sec := range o.op.Security {
		result[i] = SecurityRequirement(sec)
	}
	return result
}

func (o *operation30) GetTags() []string {
	if o.op == nil {
		return nil
	}
	return o.op.Tags
}

func (o *operation30) GetExtensions() map[string]any {
	if o.op == nil {
		return nil
	}
	return o.op.Extensions
}

func (o *operation30) GetExternalDocs() ExternalDocumentation {
	if o.op == nil || o.op.ExternalDocs == nil {
		return nil
	}
	return BaseExternalDocs{
		Description: o.op.ExternalDocs.Description,
		URL:         o.op.ExternalDocs.URL,
	}
}

func (o *operation30) GetDeprecated() bool {
	if o.op == nil {
		return false
	}
	return o.op.Deprecated
}

// parameter30 wraps OpenAPI 3.0 Parameter
type parameter30 struct {
	param *oa3.Parameter
}

func (p *parameter30) GetName() string {
	if p.param == nil {
		return ""
	}
	return p.param.Name
}

func (p *parameter30) GetIn() string {
	if p.param == nil {
		return ""
	}
	return p.param.In
}

func (p *parameter30) GetRequired() bool {
	if p.param == nil {
		return false
	}
	return p.param.Required
}

func (p *parameter30) GetDescription() string {
	if p.param == nil {
		return ""
	}
	return p.param.Description
}

func (p *parameter30) GetSchema() Schema {
	if p.param == nil || p.param.Schema == nil {
		return NilSchema{}
	}
	return &schema30{schema: p.param.Schema}
}

func (p *parameter30) IsBodyParameter() bool {
	return false // OpenAPI 3.0 doesn't have body parameters
}

func (p *parameter30) GetExtensions() map[string]any {
	if p.param == nil {
		return nil
	}
	return p.param.Extensions
}

func (p *parameter30) GetDeprecated() bool {
	if p.param == nil {
		return false
	}
	return p.param.Deprecated
}

func (p *parameter30) GetAllowEmptyValue() bool {
	if p.param == nil {
		return false
	}
	return p.param.AllowEmptyValue
}

func (p *parameter30) GetStyle() string {
	if p.param == nil {
		return ""
	}
	return p.param.Style
}

func (p *parameter30) GetExplode() bool {
	if p.param == nil || p.param.Explode == nil {
		return false
	}
	return *p.param.Explode
}

func (p *parameter30) GetAllowReserved() bool {
	if p.param == nil {
		return false
	}
	return p.param.AllowReserved
}

// requestBody30 wraps OpenAPI 3.0 RequestBody
type requestBody30 struct {
	rb *oa3.RequestBody
}

func (r *requestBody30) IsNil() bool {
	return r.rb == nil
}

func (r *requestBody30) GetRequired() bool {
	if r.rb == nil {
		return false
	}
	return r.rb.Required
}

func (r *requestBody30) GetContent() map[string]MediaType {
	if r.rb == nil || r.rb.Content == nil {
		return nil
	}
	result := make(map[string]MediaType)
	for mt, content := range r.rb.Content {
		if content != nil {
			result[mt] = &mediaType30{mt: content}
		}
	}
	return result
}

func (r *requestBody30) GetDescription() string {
	if r.rb == nil {
		return ""
	}
	return r.rb.Description
}

func (r *requestBody30) GetExtensions() map[string]any {
	if r.rb == nil {
		return nil
	}
	return r.rb.Extensions
}

// mediaType30 wraps OpenAPI 3.0 MediaType
type mediaType30 struct {
	mt *oa3.MediaType
}

func (m *mediaType30) GetSchema() Schema {
	if m.mt == nil || m.mt.Schema == nil {
		return NilSchema{}
	}
	return &schema30{schema: m.mt.Schema}
}

func (m *mediaType30) GetExtensions() map[string]any {
	if m.mt == nil {
		return nil
	}
	return m.mt.Extensions
}

// responses30 wraps OpenAPI 3.0 Responses
type responses30 struct {
	responses *oa3.Responses
}

func (r *responses30) GetDefault() Response {
	if r.responses == nil || r.responses.Default == nil {
		return NilResponse{}
	}
	return &response30{resp: r.responses.Default}
}

func (r *responses30) GetStatusCodes() map[string]Response {
	if r.responses == nil || r.responses.StatusCode == nil {
		return nil
	}
	result := make(map[string]Response)
	for code, resp := range r.responses.StatusCode {
		if resp != nil {
			result[code] = &response30{resp: resp}
		}
	}
	return result
}

func (r *responses30) GetExtensions() map[string]any {
	if r.responses == nil {
		return nil
	}
	return r.responses.Extensions
}

// response30 wraps OpenAPI 3.0 Response
type response30 struct {
	resp *oa3.Response
}

func (r *response30) IsNil() bool {
	return r.resp == nil
}

func (r *response30) HasRef() bool {
	return r.resp != nil && r.resp.IsReference()
}

func (r *response30) GetRef() string {
	if r.resp == nil {
		return ""
	}
	return r.resp.Ref
}

func (r *response30) GetDescription() string {
	if r.resp == nil {
		return ""
	}
	return r.resp.Description
}

func (r *response30) GetHeaders() map[string]Header {
	if r.resp == nil || r.resp.Headers == nil {
		return nil
	}
	result := make(map[string]Header)
	for name, header := range r.resp.Headers {
		if header != nil {
			result[name] = &header30{header: header}
		}
	}
	return result
}

func (r *response30) GetContent() map[string]MediaType {
	if r.resp == nil || r.resp.Content == nil {
		return nil
	}
	result := make(map[string]MediaType)
	for mt, content := range r.resp.Content {
		if content != nil {
			result[mt] = &mediaType30{mt: content}
		}
	}
	return result
}

func (r *response30) GetSchema() Schema {
	// OpenAPI 3.0 doesn't have schema directly on response
	return NilSchema{}
}

func (r *response30) GetExtensions() map[string]any {
	if r.resp == nil {
		return nil
	}
	return r.resp.Extensions
}

// header30 wraps OpenAPI 3.0 Header
type header30 struct {
	header *oa3.Header
}

func (h *header30) GetSchema() Schema {
	if h.header == nil || h.header.Schema == nil {
		return NilSchema{}
	}
	return &schema30{schema: h.header.Schema}
}

func (h *header30) GetRequired() bool {
	if h.header == nil {
		return false
	}
	return h.header.Required
}

func (h *header30) GetDescription() string {
	if h.header == nil {
		return ""
	}
	return h.header.Description
}

func (h *header30) GetExtensions() map[string]any {
	if h.header == nil {
		return nil
	}
	return h.header.Extensions
}

func (h *header30) GetDeprecated() bool {
	if h.header == nil {
		return false
	}
	return h.header.Deprecated
}

// schema30 wraps OpenAPI 3.0 Schema
type schema30 struct {
	schema *oa3.Schema
}

func (s *schema30) IsNil() bool {
	return s.schema == nil
}

func (s *schema30) GetRef() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Ref
}

func (s *schema30) GetType() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Type
}

func (s *schema30) GetFormat() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Format
}

func (s *schema30) GetDescription() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Description
}

func (s *schema30) GetProperties() map[string]Schema {
	if s.schema == nil || s.schema.Properties == nil {
		return nil
	}
	result := make(map[string]Schema)
	for name, prop := range s.schema.Properties {
		if prop != nil {
			result[name] = &schema30{schema: prop}
		}
	}
	return result
}

func (s *schema30) GetItems() Schema {
	if s.schema == nil || s.schema.Items == nil {
		return NilSchema{}
	}
	return &schema30{schema: s.schema.Items}
}

func (s *schema30) GetRequired() []string {
	if s.schema == nil {
		return nil
	}
	return s.schema.Required
}

func (s *schema30) GetAllOf() []Schema {
	if s.schema == nil || s.schema.AllOf == nil {
		return nil
	}
	result := make([]Schema, 0, len(s.schema.AllOf))
	for _, sub := range s.schema.AllOf {
		if sub != nil {
			result = append(result, &schema30{schema: sub})
		}
	}
	return result
}

func (s *schema30) GetOneOf() []Schema {
	if s.schema == nil || s.schema.OneOf == nil {
		return nil
	}
	result := make([]Schema, 0, len(s.schema.OneOf))
	for _, sub := range s.schema.OneOf {
		if sub != nil {
			result = append(result, &schema30{schema: sub})
		}
	}
	return result
}

func (s *schema30) GetAnyOf() []Schema {
	if s.schema == nil || s.schema.AnyOf == nil {
		return nil
	}
	result := make([]Schema, 0, len(s.schema.AnyOf))
	for _, sub := range s.schema.AnyOf {
		if sub != nil {
			result = append(result, &schema30{schema: sub})
		}
	}
	return result
}

func (s *schema30) IsBooleanSchema() bool {
	if s.schema == nil {
		return false
	}
	return s.schema.IsBooleanSchema()
}

func (s *schema30) GetBooleanValue() *bool {
	if s.schema == nil {
		return nil
	}
	return s.schema.BooleanValue()
}

func (s *schema30) GetExtensions() map[string]any {
	if s.schema == nil {
		return nil
	}
	return s.schema.Extensions
}

func (s *schema30) GetDiscriminator() Discriminator {
	if s.schema == nil || s.schema.Discriminator == nil {
		return nil
	}
	return BaseDiscriminator{
		PropertyName: s.schema.Discriminator.PropertyName,
		Mapping:      s.schema.Discriminator.Mapping,
	}
}

func (s *schema30) GetXML() XML {
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

func (s *schema30) GetExternalDocs() ExternalDocumentation {
	if s.schema == nil || s.schema.ExternalDocs == nil {
		return nil
	}
	return BaseExternalDocs{
		Description: s.schema.ExternalDocs.Description,
		URL:         s.schema.ExternalDocs.URL,
	}
}

func (s *schema30) GetExample() any {
	if s.schema == nil {
		return nil
	}
	return s.schema.Example
}

func (s *schema30) GetDefault() any {
	if s.schema == nil {
		return nil
	}
	return s.schema.Default
}

func (s *schema30) GetDeprecated() bool {
	if s.schema == nil {
		return false
	}
	return s.schema.Deprecated
}

func (s *schema30) GetReadOnly() bool {
	if s.schema == nil {
		return false
	}
	return s.schema.ReadOnly
}

func (s *schema30) GetWriteOnly() bool {
	if s.schema == nil {
		return false
	}
	return s.schema.WriteOnly
}

// securityScheme30 wraps OpenAPI 3.0 SecurityScheme
type securityScheme30 struct {
	scheme *oa3.SecurityScheme
}

func (s *securityScheme30) GetType() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Type
}

func (s *securityScheme30) GetName() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Name
}

func (s *securityScheme30) GetIn() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.In
}

func (s *securityScheme30) GetScheme() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Scheme
}

func (s *securityScheme30) GetDescription() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Description
}

func (s *securityScheme30) GetFlow() string {
	// OpenAPI 3.0 uses flows object, not a single flow string
	// Return the first available flow type
	if s.scheme == nil || s.scheme.Flows == nil {
		return ""
	}
	if s.scheme.Flows.Implicit != nil {
		return "implicit"
	}
	if s.scheme.Flows.Password != nil {
		return "password"
	}
	if s.scheme.Flows.ClientCredentials != nil {
		return "clientCredentials"
	}
	if s.scheme.Flows.AuthorizationCode != nil {
		return "authorizationCode"
	}
	return ""
}

func (s *securityScheme30) GetAuthorizationURL() string {
	if s.scheme == nil || s.scheme.Flows == nil {
		return ""
	}
	if s.scheme.Flows.Implicit != nil {
		return s.scheme.Flows.Implicit.AuthorizationUrl
	}
	if s.scheme.Flows.AuthorizationCode != nil {
		return s.scheme.Flows.AuthorizationCode.AuthorizationUrl
	}
	return ""
}

func (s *securityScheme30) GetTokenURL() string {
	if s.scheme == nil || s.scheme.Flows == nil {
		return ""
	}
	if s.scheme.Flows.Password != nil {
		return s.scheme.Flows.Password.TokenUrl
	}
	if s.scheme.Flows.ClientCredentials != nil {
		return s.scheme.Flows.ClientCredentials.TokenUrl
	}
	if s.scheme.Flows.AuthorizationCode != nil {
		return s.scheme.Flows.AuthorizationCode.TokenUrl
	}
	return ""
}

func (s *securityScheme30) GetScopes() map[string]string {
	if s.scheme == nil || s.scheme.Flows == nil {
		return nil
	}
	if s.scheme.Flows.Implicit != nil && s.scheme.Flows.Implicit.Scopes != nil {
		return s.scheme.Flows.Implicit.Scopes
	}
	if s.scheme.Flows.Password != nil && s.scheme.Flows.Password.Scopes != nil {
		return s.scheme.Flows.Password.Scopes
	}
	if s.scheme.Flows.ClientCredentials != nil && s.scheme.Flows.ClientCredentials.Scopes != nil {
		return s.scheme.Flows.ClientCredentials.Scopes
	}
	if s.scheme.Flows.AuthorizationCode != nil && s.scheme.Flows.AuthorizationCode.Scopes != nil {
		return s.scheme.Flows.AuthorizationCode.Scopes
	}
	return nil
}

func (s *securityScheme30) GetExtensions() map[string]any {
	if s.scheme == nil {
		return nil
	}
	return s.scheme.Extensions
}
