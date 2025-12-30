// Package unified provides an adapter for OpenAPI 3.1 documents
// Copyright (c) Greetingland LLC

package unified

import (
	"strings"

	oa31 "github.com/genelet/oas/openapi31"
)

// Document31 wraps an OpenAPI 3.1 document and implements Document
type Document31 struct {
	doc *oa31.OpenAPI
}

// NewDocument31 creates a new Document adapter for OpenAPI 3.1
func NewDocument31(doc *oa31.OpenAPI) Document {
	return &Document31{doc: doc}
}

// GetRaw returns the underlying OpenAPI 3.1 document
func (d *Document31) GetRaw() *oa31.OpenAPI {
	return d.doc
}

func (d *Document31) Version() string {
	return d.doc.OpenAPI
}

func (d *Document31) GetServerURL() string {
	if len(d.doc.Servers) > 0 && d.doc.Servers[0] != nil {
		return d.doc.Servers[0].URL
	}
	return ""
}

func (d *Document31) GetInfo() DocumentInfo {
	if d.doc.Info == nil {
		return &documentInfo31{}
	}
	return &documentInfo31{info: d.doc.Info}
}

func (d *Document31) GetPaths() map[string]PathItem {
	if d.doc.Paths == nil || d.doc.Paths.Paths == nil {
		return nil
	}
	result := make(map[string]PathItem)
	for path, item := range d.doc.Paths.Paths {
		if item != nil {
			result[path] = &pathItem31{item: item}
		}
	}
	return result
}

func (d *Document31) GetSecuritySchemes() map[string]SecurityScheme {
	if d.doc.Components == nil || d.doc.Components.SecuritySchemes == nil {
		return nil
	}
	result := make(map[string]SecurityScheme)
	for name, scheme := range d.doc.Components.SecuritySchemes {
		if scheme != nil {
			result[name] = &securityScheme31{scheme: scheme}
		}
	}
	return result
}

func (d *Document31) GetGlobalSecurity() []SecurityRequirement {
	if d.doc.Security == nil {
		return nil
	}
	result := make([]SecurityRequirement, len(d.doc.Security))
	for i, sec := range d.doc.Security {
		result[i] = SecurityRequirement(sec)
	}
	return result
}

func (d *Document31) GetExtensions() map[string]any {
	return d.doc.Extensions
}

// documentInfo31 wraps OpenAPI 3.1 Info
type documentInfo31 struct {
	info *oa31.Info
}

func (i *documentInfo31) GetTitle() string {
	if i.info == nil {
		return ""
	}
	return i.info.Title
}

func (i *documentInfo31) GetVersion() string {
	if i.info == nil {
		return ""
	}
	return i.info.Version
}

func (i *documentInfo31) GetDescription() string {
	if i.info == nil {
		return ""
	}
	return i.info.Description
}

func (i *documentInfo31) GetExtensions() map[string]any {
	if i.info == nil {
		return nil
	}
	return i.info.Extensions
}

// pathItem31 wraps OpenAPI 3.1 PathItem
type pathItem31 struct {
	item *oa31.PathItem
}

func (p *pathItem31) HasRef() bool {
	return p.item != nil && p.item.HasRef()
}

func (p *pathItem31) GetRef() string {
	if p.item == nil {
		return ""
	}
	return p.item.GetRef()
}

func (p *pathItem31) GetOperation(method string) Operation {
	if p.item == nil {
		return NilOperation{}
	}
	method = strings.ToLower(method)
	var op *oa31.Operation
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
	return &operation31{op: op}
}

func (p *pathItem31) GetAllOperations() map[string]Operation {
	if p.item == nil {
		return nil
	}
	result := make(map[string]Operation)
	if p.item.Get != nil {
		result["get"] = &operation31{op: p.item.Get}
	}
	if p.item.Put != nil {
		result["put"] = &operation31{op: p.item.Put}
	}
	if p.item.Post != nil {
		result["post"] = &operation31{op: p.item.Post}
	}
	if p.item.Delete != nil {
		result["delete"] = &operation31{op: p.item.Delete}
	}
	if p.item.Options != nil {
		result["options"] = &operation31{op: p.item.Options}
	}
	if p.item.Head != nil {
		result["head"] = &operation31{op: p.item.Head}
	}
	if p.item.Patch != nil {
		result["patch"] = &operation31{op: p.item.Patch}
	}
	if p.item.Trace != nil {
		result["trace"] = &operation31{op: p.item.Trace}
	}
	return result
}

func (p *pathItem31) GetParameters() []Parameter {
	if p.item == nil || p.item.Parameters == nil {
		return nil
	}
	result := make([]Parameter, 0, len(p.item.Parameters))
	for _, param := range p.item.Parameters {
		if param != nil {
			result = append(result, &parameter31{param: param})
		}
	}
	return result
}

func (p *pathItem31) GetExtensions() map[string]any {
	if p.item == nil {
		return nil
	}
	return p.item.Extensions
}

// operation31 wraps OpenAPI 3.1 Operation
type operation31 struct {
	op *oa31.Operation
}

func (o *operation31) IsNil() bool {
	return o.op == nil
}

func (o *operation31) GetOperationID() string {
	if o.op == nil {
		return ""
	}
	return o.op.OperationID
}

func (o *operation31) GetSummary() string {
	if o.op == nil {
		return ""
	}
	return o.op.Summary
}

func (o *operation31) GetDescription() string {
	if o.op == nil {
		return ""
	}
	return o.op.Description
}

func (o *operation31) GetParameters() []Parameter {
	if o.op == nil || o.op.Parameters == nil {
		return nil
	}
	result := make([]Parameter, 0, len(o.op.Parameters))
	for _, param := range o.op.Parameters {
		if param != nil {
			result = append(result, &parameter31{param: param})
		}
	}
	return result
}

func (o *operation31) GetRequestBody() RequestBody {
	if o.op == nil || o.op.RequestBody == nil {
		return NilRequestBody{}
	}
	return &requestBody31{rb: o.op.RequestBody}
}

func (o *operation31) GetResponses() Responses {
	if o.op == nil || o.op.Responses == nil {
		return NilResponses{}
	}
	return &responses31{responses: o.op.Responses}
}

func (o *operation31) GetSecurity() []SecurityRequirement {
	if o.op == nil || o.op.Security == nil {
		return nil
	}
	result := make([]SecurityRequirement, len(o.op.Security))
	for i, sec := range o.op.Security {
		result[i] = SecurityRequirement(sec)
	}
	return result
}

func (o *operation31) GetTags() []string {
	if o.op == nil {
		return nil
	}
	return o.op.Tags
}

func (o *operation31) GetExtensions() map[string]any {
	if o.op == nil {
		return nil
	}
	return o.op.Extensions
}

func (o *operation31) GetExternalDocs() ExternalDocumentation {
	if o.op == nil || o.op.ExternalDocs == nil {
		return nil
	}
	return BaseExternalDocs{
		Description: o.op.ExternalDocs.Description,
		URL:         o.op.ExternalDocs.URL,
	}
}

func (o *operation31) GetDeprecated() bool {
	if o.op == nil {
		return false
	}
	// Missing in 3.1 schema struct, returning false as default
	return false
}

// parameter31 wraps OpenAPI 3.1 Parameter
type parameter31 struct {
	param *oa31.Parameter
}

func (p *parameter31) GetName() string {
	if p.param == nil {
		return ""
	}
	return p.param.Name
}

func (p *parameter31) GetIn() string {
	if p.param == nil {
		return ""
	}
	return p.param.In
}

func (p *parameter31) GetRequired() bool {
	if p.param == nil {
		return false
	}
	return p.param.Required
}

func (p *parameter31) GetDescription() string {
	if p.param == nil {
		return ""
	}
	return p.param.Description
}

func (p *parameter31) GetSchema() Schema {
	if p.param == nil || p.param.Schema == nil {
		return NilSchema{}
	}
	return &schema31{schema: p.param.Schema}
}

func (p *parameter31) IsBodyParameter() bool {
	return false // OpenAPI 3.1 doesn't have body parameters
}

func (p *parameter31) GetExtensions() map[string]any {
	if p.param == nil {
		return nil
	}
	return p.param.Extensions
}

func (p *parameter31) GetDeprecated() bool {
	if p.param == nil {
		return false
	}
	return p.param.Deprecated
}

func (p *parameter31) GetAllowEmptyValue() bool {
	if p.param == nil {
		return false
	}
	return p.param.AllowEmptyValue
}

func (p *parameter31) GetStyle() string {
	if p.param == nil {
		return ""
	}
	return p.param.Style
}

func (p *parameter31) GetExplode() bool {
	if p.param == nil || p.param.Explode == nil {
		// Default values for explode depend on style, but for now defaulting to false if nil
		// OpenApi spec says: "When style is form, the default value is true. For all other styles, the default value is false."
		// However, adhering to simplifcation here:
		return false
	}
	return *p.param.Explode
}

func (p *parameter31) GetAllowReserved() bool {
	if p.param == nil {
		return false
	}
	return p.param.AllowReserved
}

// requestBody31 wraps OpenAPI 3.1 RequestBody
type requestBody31 struct {
	rb *oa31.RequestBody
}

func (r *requestBody31) IsNil() bool {
	return r.rb == nil
}

func (r *requestBody31) GetRequired() bool {
	if r.rb == nil {
		return false
	}
	return r.rb.Required
}

func (r *requestBody31) GetContent() map[string]MediaType {
	if r.rb == nil || r.rb.Content == nil {
		return nil
	}
	result := make(map[string]MediaType)
	for mt, content := range r.rb.Content {
		if content != nil {
			result[mt] = &mediaType31{mt: content}
		}
	}
	return result
}

func (r *requestBody31) GetDescription() string {
	if r.rb == nil {
		return ""
	}
	return r.rb.Description
}

func (r *requestBody31) GetExtensions() map[string]any {
	if r.rb == nil {
		return nil
	}
	return r.rb.Extensions
}

// mediaType31 wraps OpenAPI 3.1 MediaType
type mediaType31 struct {
	mt *oa31.MediaType
}

func (m *mediaType31) GetSchema() Schema {
	if m.mt == nil || m.mt.Schema == nil {
		return NilSchema{}
	}
	return &schema31{schema: m.mt.Schema}
}

func (m *mediaType31) GetExtensions() map[string]any {
	if m.mt == nil {
		return nil
	}
	return m.mt.Extensions
}

// responses31 wraps OpenAPI 3.1 Responses
type responses31 struct {
	responses *oa31.Responses
}

func (r *responses31) GetDefault() Response {
	if r.responses == nil || r.responses.Default == nil {
		return NilResponse{}
	}
	return &response31{resp: r.responses.Default}
}

func (r *responses31) GetStatusCodes() map[string]Response {
	if r.responses == nil || r.responses.StatusCode == nil {
		return nil
	}
	result := make(map[string]Response)
	for code, resp := range r.responses.StatusCode {
		if resp != nil {
			result[code] = &response31{resp: resp}
		}
	}
	return result
}

func (r *responses31) GetExtensions() map[string]any {
	if r.responses == nil {
		return nil
	}
	return r.responses.Extensions
}

// response31 wraps OpenAPI 3.1 Response
type response31 struct {
	resp *oa31.Response
}

func (r *response31) IsNil() bool {
	return r.resp == nil
}

func (r *response31) HasRef() bool {
	return r.resp != nil && r.resp.IsReference()
}

func (r *response31) GetRef() string {
	if r.resp == nil {
		return ""
	}
	return r.resp.Ref
}

func (r *response31) GetDescription() string {
	if r.resp == nil {
		return ""
	}
	return r.resp.Description
}

func (r *response31) GetHeaders() map[string]Header {
	if r.resp == nil || r.resp.Headers == nil {
		return nil
	}
	result := make(map[string]Header)
	for name, header := range r.resp.Headers {
		if header != nil {
			result[name] = &header31{header: header}
		}
	}
	return result
}

func (r *response31) GetContent() map[string]MediaType {
	if r.resp == nil || r.resp.Content == nil {
		return nil
	}
	result := make(map[string]MediaType)
	for mt, content := range r.resp.Content {
		if content != nil {
			result[mt] = &mediaType31{mt: content}
		}
	}
	return result
}

func (r *response31) GetSchema() Schema {
	// OpenAPI 3.1 doesn't have schema directly on response
	return NilSchema{}
}

func (r *response31) GetExtensions() map[string]any {
	if r.resp == nil {
		return nil
	}
	return r.resp.Extensions
}

// header31 wraps OpenAPI 3.1 Header
type header31 struct {
	header *oa31.Header
}

func (h *header31) GetSchema() Schema {
	if h.header == nil || h.header.Schema == nil {
		return NilSchema{}
	}
	return &schema31{schema: h.header.Schema}
}

func (h *header31) GetRequired() bool {
	if h.header == nil {
		return false
	}
	return h.header.Required
}

func (h *header31) GetDescription() string {
	if h.header == nil {
		return ""
	}
	return h.header.Description
}

func (h *header31) GetExtensions() map[string]any {
	if h.header == nil {
		return nil
	}
	return h.header.Extensions
}

func (h *header31) GetDeprecated() bool {
	if h.header == nil {
		return false
	}
	return h.header.Deprecated
}

// schema31 wraps OpenAPI 3.1 Schema
type schema31 struct {
	schema *oa31.Schema
}

func (s *schema31) IsNil() bool {
	return s.schema == nil
}

func (s *schema31) GetRef() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Ref
}

func (s *schema31) GetType() string {
	if s.schema == nil || s.schema.Type == nil {
		return ""
	}
	// In 3.1, type can be a string or array of strings
	// Return the first non-null type for compatibility
	if s.schema.Type.String != "" {
		return s.schema.Type.String
	}
	for _, t := range s.schema.Type.Array {
		if t != "null" {
			return t
		}
	}
	if len(s.schema.Type.Array) > 0 {
		return s.schema.Type.Array[0]
	}
	return ""
}

func (s *schema31) GetFormat() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Format
}

func (s *schema31) GetDescription() string {
	if s.schema == nil {
		return ""
	}
	return s.schema.Description
}

func (s *schema31) GetProperties() map[string]Schema {
	if s.schema == nil || s.schema.Properties == nil {
		return nil
	}
	result := make(map[string]Schema)
	for name, prop := range s.schema.Properties {
		if prop != nil {
			result[name] = &schema31{schema: prop}
		}
	}
	return result
}

func (s *schema31) GetItems() Schema {
	if s.schema == nil || s.schema.Items == nil {
		return NilSchema{}
	}
	return &schema31{schema: s.schema.Items}
}

func (s *schema31) GetRequired() []string {
	if s.schema == nil {
		return nil
	}
	return s.schema.Required
}

func (s *schema31) GetAllOf() []Schema {
	if s.schema == nil || s.schema.AllOf == nil {
		return nil
	}
	result := make([]Schema, 0, len(s.schema.AllOf))
	for _, sub := range s.schema.AllOf {
		if sub != nil {
			result = append(result, &schema31{schema: sub})
		}
	}
	return result
}

func (s *schema31) GetOneOf() []Schema {
	if s.schema == nil || s.schema.OneOf == nil {
		return nil
	}
	result := make([]Schema, 0, len(s.schema.OneOf))
	for _, sub := range s.schema.OneOf {
		if sub != nil {
			result = append(result, &schema31{schema: sub})
		}
	}
	return result
}

func (s *schema31) GetAnyOf() []Schema {
	if s.schema == nil || s.schema.AnyOf == nil {
		return nil
	}
	result := make([]Schema, 0, len(s.schema.AnyOf))
	for _, sub := range s.schema.AnyOf {
		if sub != nil {
			result = append(result, &schema31{schema: sub})
		}
	}
	return result
}

func (s *schema31) IsBooleanSchema() bool {
	if s.schema == nil {
		return false
	}
	return s.schema.IsBooleanSchema()
}

func (s *schema31) GetBooleanValue() *bool {
	if s.schema == nil {
		return nil
	}
	return s.schema.BooleanValue()
}

func (s *schema31) GetExtensions() map[string]any {
	if s.schema == nil {
		return nil
	}
	return s.schema.Extensions
}

func (s *schema31) GetDiscriminator() Discriminator {
	if s.schema == nil || s.schema.Discriminator == nil {
		return nil
	}
	return BaseDiscriminator{
		PropertyName: s.schema.Discriminator.PropertyName,
		Mapping:      s.schema.Discriminator.Mapping,
	}
}

func (s *schema31) GetXML() XML {
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

func (s *schema31) GetExternalDocs() ExternalDocumentation {
	if s.schema == nil || s.schema.ExternalDocs == nil {
		return nil
	}
	return BaseExternalDocs{
		Description: s.schema.ExternalDocs.Description,
		URL:         s.schema.ExternalDocs.URL,
	}
}

func (s *schema31) GetExample() any {
	if s.schema == nil {
		return nil
	}
	return s.schema.Example
}

func (s *schema31) GetDefault() any {
	if s.schema == nil {
		return nil
	}
	return s.schema.Default
}

func (s *schema31) GetDeprecated() bool {
	if s.schema == nil {
		return false
	}
	return s.schema.Deprecated
}

func (s *schema31) GetReadOnly() bool {
	if s.schema == nil {
		return false
	}
	return s.schema.ReadOnly
}

func (s *schema31) GetWriteOnly() bool {
	if s.schema == nil {
		return false
	}
	return s.schema.WriteOnly
}

// securityScheme31 wraps OpenAPI 3.1 SecurityScheme
type securityScheme31 struct {
	scheme *oa31.SecurityScheme
}

func (s *securityScheme31) GetType() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Type
}

func (s *securityScheme31) GetName() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Name
}

func (s *securityScheme31) GetIn() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.In
}

func (s *securityScheme31) GetScheme() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Scheme
}

func (s *securityScheme31) GetDescription() string {
	if s.scheme == nil {
		return ""
	}
	return s.scheme.Description
}

func (s *securityScheme31) GetFlow() string {
	// OpenAPI 3.1 uses flows object, not a single flow string
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

func (s *securityScheme31) GetAuthorizationURL() string {
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

func (s *securityScheme31) GetTokenURL() string {
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

func (s *securityScheme31) GetScopes() map[string]string {
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

func (s *securityScheme31) GetExtensions() map[string]any {
	if s.scheme == nil {
		return nil
	}
	return s.scheme.Extensions
}
