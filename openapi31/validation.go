// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a validation error with path context
type ValidationError struct {
	Path    string
	Message string
}

func (e ValidationError) Error() string {
	if e.Path == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationResult contains all validation errors
type ValidationResult struct {
	Errors []ValidationError
}

// Valid returns true if there are no validation errors
func (r *ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

// Error returns a combined error message
func (r *ValidationResult) Error() string {
	if r.Valid() {
		return ""
	}
	var msgs []string
	for _, e := range r.Errors {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "; ")
}

func (r *ValidationResult) addError(path, message string) {
	r.Errors = append(r.Errors, ValidationError{Path: path, Message: message})
}

// Validate validates the OpenAPI document against the OpenAPI 3.1 specification
func (o *OpenAPI) Validate() *ValidationResult {
	result := &ValidationResult{}

	if o == nil {
		result.addError("", "OpenAPI document is nil")
		return result
	}

	// Required: openapi
	if o.OpenAPI == "" {
		result.addError("openapi", "required field is missing")
	} else if !strings.HasPrefix(o.OpenAPI, "3.1") {
		result.addError("openapi", fmt.Sprintf("expected 3.1.x version, got %s", o.OpenAPI))
	}

	// Required: info
	if o.Info == nil {
		result.addError("info", "required field is missing")
	} else {
		o.Info.validate("info", result)
	}

	// At least one of: paths, webhooks, or components
	hasPaths := o.Paths != nil && len(o.Paths.Paths) > 0
	hasWebhooks := len(o.Webhooks) > 0
	hasComponents := o.Components != nil
	if !hasPaths && !hasWebhooks && !hasComponents {
		result.addError("", "must have at least one of: paths, webhooks, or components")
	}

	// Optional: paths
	if o.Paths != nil {
		o.Paths.validate("paths", result)
	}

	// Optional: webhooks
	for name, pathItem := range o.Webhooks {
		if pathItem != nil {
			pathItem.validate(fmt.Sprintf("webhooks[%s]", name), result)
		}
	}

	// Optional: servers
	for i, server := range o.Servers {
		if server != nil {
			server.validate(fmt.Sprintf("servers[%d]", i), result)
		}
	}

	// Optional: components
	if o.Components != nil {
		o.Components.validate("components", result)
	}

	// Optional: tags
	for i, tag := range o.Tags {
		if tag != nil {
			tag.validate(fmt.Sprintf("tags[%d]", i), result)
		}
	}

	return result
}

func (i *Info) validate(path string, result *ValidationResult) {
	// Required: title
	if i.Title == "" {
		result.addError(path+".title", "required field is missing")
	}
	// Required: version
	if i.Version == "" {
		result.addError(path+".version", "required field is missing")
	}
	// Optional: license
	if i.License != nil {
		i.License.validate(path+".license", result)
	}
}

func (l *License) validate(path string, result *ValidationResult) {
	// Required: name
	if l.Name == "" {
		result.addError(path+".name", "required field is missing")
	}
	// Mutual exclusion: identifier and url cannot both be present
	if l.Identifier != "" && l.URL != "" {
		result.addError(path, "identifier and url are mutually exclusive")
	}
}

func (s *Server) validate(path string, result *ValidationResult) {
	// Required: url
	if s.URL == "" {
		result.addError(path+".url", "required field is missing")
	}
	// Validate variables
	for name, v := range s.Variables {
		if v != nil {
			v.validate(fmt.Sprintf("%s.variables[%s]", path, name), result)
		}
	}
}

func (v *ServerVariable) validate(path string, result *ValidationResult) {
	// Required: default
	if v.Default == "" {
		result.addError(path+".default", "required field is missing")
	}
	// If enum is provided, default must be in enum
	if len(v.Enum) > 0 {
		found := false
		for _, e := range v.Enum {
			if e == v.Default {
				found = true
				break
			}
		}
		if !found {
			result.addError(path+".default", "default value must be one of the enum values")
		}
	}
}

func (p *Paths) validate(path string, result *ValidationResult) {
	for pathPattern, pathItem := range p.Paths {
		// Path must start with /
		if !strings.HasPrefix(pathPattern, "/") {
			result.addError(path+"."+pathPattern, "path must start with /")
		}
		if pathItem != nil {
			pathItem.validate(fmt.Sprintf("%s[%s]", path, pathPattern), result)
		}
	}
}

func (p *PathItem) validate(path string, result *ValidationResult) {
	// Skip validation for references
	if p.Ref != "" {
		return
	}

	// Validate parameters at path level
	for i, param := range p.Parameters {
		if param != nil {
			param.validate(fmt.Sprintf("%s.parameters[%d]", path, i), result)
		}
	}

	// Validate operations
	if p.Get != nil {
		p.Get.validate(path+".get", result)
	}
	if p.Put != nil {
		p.Put.validate(path+".put", result)
	}
	if p.Post != nil {
		p.Post.validate(path+".post", result)
	}
	if p.Delete != nil {
		p.Delete.validate(path+".delete", result)
	}
	if p.Options != nil {
		p.Options.validate(path+".options", result)
	}
	if p.Head != nil {
		p.Head.validate(path+".head", result)
	}
	if p.Patch != nil {
		p.Patch.validate(path+".patch", result)
	}
	if p.Trace != nil {
		p.Trace.validate(path+".trace", result)
	}
}

func (o *Operation) validate(path string, result *ValidationResult) {
	// Required: responses (unless it's a webhook)
	if o.Responses == nil {
		result.addError(path+".responses", "required field is missing")
	} else {
		o.Responses.validate(path+".responses", result)
	}

	// Validate parameters
	for i, param := range o.Parameters {
		if param != nil {
			param.validate(fmt.Sprintf("%s.parameters[%d]", path, i), result)
		}
	}

	// Validate requestBody
	if o.RequestBody != nil {
		o.RequestBody.validate(path+".requestBody", result)
	}

	// Validate callbacks
	for name, callback := range o.Callbacks {
		if callback != nil {
			callback.validate(fmt.Sprintf("%s.callbacks[%s]", path, name), result)
		}
	}
}

func (r *Responses) validate(path string, result *ValidationResult) {
	// minProperties: 1 - must have at least one response
	hasResponse := r.Default != nil || len(r.StatusCode) > 0
	if !hasResponse {
		result.addError(path, "must contain at least one response")
	}

	// Validate status code pattern
	statusCodePattern := regexp.MustCompile(`^[1-5][0-9][0-9]$|^[1-5]XX$`)
	for code, resp := range r.StatusCode {
		if !statusCodePattern.MatchString(code) {
			result.addError(path+"."+code, "invalid status code pattern, must be 3-digit code or pattern like 2XX")
		}
		if resp != nil {
			resp.validate(path+"."+code, result)
		}
	}

	if r.Default != nil {
		r.Default.validate(path+".default", result)
	}
}

func (r *Response) validate(path string, result *ValidationResult) {
	// Skip validation for references
	if r.IsReference() {
		return
	}

	// Required: description
	if r.Description == "" {
		result.addError(path+".description", "required field is missing")
	}

	// Validate headers
	for name, header := range r.Headers {
		if header != nil {
			header.validate(fmt.Sprintf("%s.headers[%s]", path, name), result)
		}
	}

	// Validate content
	for mediaType, mt := range r.Content {
		if mt != nil {
			mt.validate(fmt.Sprintf("%s.content[%s]", path, mediaType), result)
		}
	}

	// Validate links
	for name, link := range r.Links {
		if link != nil {
			link.validate(fmt.Sprintf("%s.links[%s]", path, name), result)
		}
	}
}

func (p *Parameter) validate(path string, result *ValidationResult) {
	// Skip validation for references
	if p.IsReference() {
		return
	}

	// Required: name
	if p.Name == "" {
		result.addError(path+".name", "required field is missing")
	}

	// Required: in
	if p.In == "" {
		result.addError(path+".in", "required field is missing")
	} else {
		validIn := map[string]bool{"query": true, "header": true, "path": true, "cookie": true}
		if !validIn[p.In] {
			result.addError(path+".in", fmt.Sprintf("must be one of: query, header, path, cookie; got %s", p.In))
		}
	}

	// Path parameters must have required: true
	if p.In == "path" && !p.Required {
		result.addError(path+".required", "path parameters must have required: true")
	}

	// Validate style based on 'in' value
	if p.Style != "" {
		validStyles := map[string][]string{
			"path":   {"matrix", "label", "simple"},
			"query":  {"form", "spaceDelimited", "pipeDelimited", "deepObject"},
			"header": {"simple"},
			"cookie": {"form"},
		}
		if styles, ok := validStyles[p.In]; ok {
			valid := false
			for _, s := range styles {
				if s == p.Style {
					valid = true
					break
				}
			}
			if !valid {
				result.addError(path+".style", fmt.Sprintf("invalid style '%s' for parameter in '%s'", p.Style, p.In))
			}
		}
	}

	// Schema XOR Content - must have one but not both
	hasSchema := p.Schema != nil
	hasContent := len(p.Content) > 0
	if !hasSchema && !hasContent {
		result.addError(path, "must have either 'schema' or 'content'")
	}
	if hasSchema && hasContent {
		result.addError(path, "cannot have both 'schema' and 'content'")
	}

	// Example XOR Examples - cannot have both
	if p.Example != nil && len(p.Examples) > 0 {
		result.addError(path, "cannot have both 'example' and 'examples'")
	}

	// Validate schema
	if p.Schema != nil {
		p.Schema.validate(path+".schema", result)
	}
}

func (h *Header) validate(path string, result *ValidationResult) {
	// Skip validation for references
	if h.IsReference() {
		return
	}

	// Schema XOR Content
	hasSchema := h.Schema != nil
	hasContent := len(h.Content) > 0
	if !hasSchema && !hasContent {
		result.addError(path, "must have either 'schema' or 'content'")
	}
	if hasSchema && hasContent {
		result.addError(path, "cannot have both 'schema' and 'content'")
	}

	// Example XOR Examples
	if h.Example != nil && len(h.Examples) > 0 {
		result.addError(path, "cannot have both 'example' and 'examples'")
	}

	// Style must be 'simple' for headers
	if h.Style != "" && h.Style != "simple" {
		result.addError(path+".style", "header style must be 'simple'")
	}
}

func (r *RequestBody) validate(path string, result *ValidationResult) {
	// Skip validation for references
	if r.IsReference() {
		return
	}

	// Required: content
	if len(r.Content) == 0 {
		result.addError(path+".content", "required field is missing")
	}

	// Validate content
	for mediaType, mt := range r.Content {
		if mt != nil {
			mt.validate(fmt.Sprintf("%s.content[%s]", path, mediaType), result)
		}
	}
}

func (m *MediaType) validate(path string, result *ValidationResult) {
	// Example XOR Examples
	if m.Example != nil && len(m.Examples) > 0 {
		result.addError(path, "cannot have both 'example' and 'examples'")
	}

	// Validate schema
	if m.Schema != nil {
		m.Schema.validate(path+".schema", result)
	}

	// Validate encoding
	for name, enc := range m.Encoding {
		if enc != nil {
			enc.validate(fmt.Sprintf("%s.encoding[%s]", path, name), result)
		}
	}
}

func (e *Encoding) validate(path string, result *ValidationResult) {
	// Validate style
	if e.Style != "" {
		validStyles := []string{"form", "spaceDelimited", "pipeDelimited", "deepObject"}
		valid := false
		for _, s := range validStyles {
			if s == e.Style {
				valid = true
				break
			}
		}
		if !valid {
			result.addError(path+".style", fmt.Sprintf("invalid encoding style '%s'", e.Style))
		}
	}

	// Validate headers
	for name, header := range e.Headers {
		if header != nil {
			header.validate(fmt.Sprintf("%s.headers[%s]", path, name), result)
		}
	}
}

func (s *Schema) validate(path string, result *ValidationResult) {
	// Boolean schemas are always valid
	if s.IsBooleanSchema() {
		return
	}

	// References are valid (resolution is separate concern)
	if s.Ref != "" {
		return
	}

	// Validate type
	if s.Type != nil && !s.Type.IsEmpty() {
		validTypes := []string{"string", "number", "integer", "boolean", "array", "object", "null"}
		// Get all types from the StringOrStringArray
		var types []string
		if s.Type.String != "" {
			types = []string{s.Type.String}
		} else {
			types = s.Type.Array
		}
		for _, t := range types {
			valid := false
			for _, vt := range validTypes {
				if t == vt {
					valid = true
					break
				}
			}
			if !valid {
				result.addError(path+".type", fmt.Sprintf("invalid type '%s'", t))
			}
		}
	}

	// Array type must have items or prefixItems
	if s.Type != nil && s.Type.Contains("array") && s.Items == nil && len(s.PrefixItems) == 0 {
		result.addError(path, "array type should have items or prefixItems defined")
	}

	// Validate numeric constraints
	if s.Minimum != nil && s.Maximum != nil {
		if *s.Minimum > *s.Maximum {
			result.addError(path, "minimum cannot be greater than maximum")
		}
	}
	if s.ExclusiveMinimum != nil && s.ExclusiveMaximum != nil {
		if *s.ExclusiveMinimum >= *s.ExclusiveMaximum {
			result.addError(path, "exclusiveMinimum must be less than exclusiveMaximum")
		}
	}
	if s.MinLength != nil && s.MaxLength != nil {
		if *s.MinLength > *s.MaxLength {
			result.addError(path, "minLength cannot be greater than maxLength")
		}
	}
	if s.MinItems != nil && s.MaxItems != nil {
		if *s.MinItems > *s.MaxItems {
			result.addError(path, "minItems cannot be greater than maxItems")
		}
	}
	if s.MinProperties != nil && s.MaxProperties != nil {
		if *s.MinProperties > *s.MaxProperties {
			result.addError(path, "minProperties cannot be greater than maxProperties")
		}
	}
	if s.MinContains != nil && s.MaxContains != nil {
		if *s.MinContains > *s.MaxContains {
			result.addError(path, "minContains cannot be greater than maxContains")
		}
	}

	// Validate pattern is valid regex
	if s.Pattern != "" {
		if _, err := regexp.Compile(s.Pattern); err != nil {
			result.addError(path+".pattern", fmt.Sprintf("invalid regex pattern: %v", err))
		}
	}

	// Validate required fields exist in properties
	if len(s.Required) > 0 && len(s.Properties) > 0 {
		for _, req := range s.Required {
			if _, exists := s.Properties[req]; !exists {
				result.addError(path+".required", fmt.Sprintf("required property '%s' not defined in properties", req))
			}
		}
	}

	// Validate nested schemas
	if s.Items != nil {
		s.Items.validate(path+".items", result)
	}
	for i, schema := range s.PrefixItems {
		if schema != nil {
			schema.validate(fmt.Sprintf("%s.prefixItems[%d]", path, i), result)
		}
	}
	for name, prop := range s.Properties {
		if prop != nil {
			prop.validate(fmt.Sprintf("%s.properties[%s]", path, name), result)
		}
	}
	if s.AdditionalProperties != nil {
		s.AdditionalProperties.validate(path+".additionalProperties", result)
	}
	for i, schema := range s.AllOf {
		if schema != nil {
			schema.validate(fmt.Sprintf("%s.allOf[%d]", path, i), result)
		}
	}
	for i, schema := range s.AnyOf {
		if schema != nil {
			schema.validate(fmt.Sprintf("%s.anyOf[%d]", path, i), result)
		}
	}
	for i, schema := range s.OneOf {
		if schema != nil {
			schema.validate(fmt.Sprintf("%s.oneOf[%d]", path, i), result)
		}
	}
	if s.Not != nil {
		s.Not.validate(path+".not", result)
	}
	if s.If != nil {
		s.If.validate(path+".if", result)
	}
	if s.Then != nil {
		s.Then.validate(path+".then", result)
	}
	if s.Else != nil {
		s.Else.validate(path+".else", result)
	}
	if s.Contains != nil {
		s.Contains.validate(path+".contains", result)
	}
	if s.UnevaluatedItems != nil {
		s.UnevaluatedItems.validate(path+".unevaluatedItems", result)
	}
	if s.UnevaluatedProperties != nil {
		s.UnevaluatedProperties.validate(path+".unevaluatedProperties", result)
	}
}

func (l *Link) validate(path string, result *ValidationResult) {
	// Skip validation for references
	if l.IsReference() {
		return
	}

	// operationId XOR operationRef - cannot have both
	if l.OperationId != "" && l.OperationRef != "" {
		result.addError(path, "cannot have both 'operationId' and 'operationRef'")
	}
}

func (c *Callback) validate(path string, result *ValidationResult) {
	// Skip validation for references
	if c.IsReference() {
		return
	}

	for expr, pathItem := range c.Paths {
		if pathItem != nil {
			pathItem.validate(fmt.Sprintf("%s[%s]", path, expr), result)
		}
	}
}

func (t *Tag) validate(path string, result *ValidationResult) {
	// Required: name
	if t.Name == "" {
		result.addError(path+".name", "required field is missing")
	}
}

func (c *Components) validate(path string, result *ValidationResult) {
	// Validate component name pattern
	namePattern := regexp.MustCompile(`^[a-zA-Z0-9\.\-_]+$`)

	// Validate schemas
	for name, schema := range c.Schemas {
		if !namePattern.MatchString(name) {
			result.addError(fmt.Sprintf("%s.schemas[%s]", path, name), "component name contains invalid characters")
		}
		if schema != nil {
			schema.validate(fmt.Sprintf("%s.schemas[%s]", path, name), result)
		}
	}

	// Validate responses
	for name, resp := range c.Responses {
		if !namePattern.MatchString(name) {
			result.addError(fmt.Sprintf("%s.responses[%s]", path, name), "component name contains invalid characters")
		}
		if resp != nil {
			resp.validate(fmt.Sprintf("%s.responses[%s]", path, name), result)
		}
	}

	// Validate parameters
	for name, param := range c.Parameters {
		if !namePattern.MatchString(name) {
			result.addError(fmt.Sprintf("%s.parameters[%s]", path, name), "component name contains invalid characters")
		}
		if param != nil {
			param.validate(fmt.Sprintf("%s.parameters[%s]", path, name), result)
		}
	}

	// Validate requestBodies
	for name, rb := range c.RequestBodies {
		if !namePattern.MatchString(name) {
			result.addError(fmt.Sprintf("%s.requestBodies[%s]", path, name), "component name contains invalid characters")
		}
		if rb != nil {
			rb.validate(fmt.Sprintf("%s.requestBodies[%s]", path, name), result)
		}
	}

	// Validate headers
	for name, header := range c.Headers {
		if !namePattern.MatchString(name) {
			result.addError(fmt.Sprintf("%s.headers[%s]", path, name), "component name contains invalid characters")
		}
		if header != nil {
			header.validate(fmt.Sprintf("%s.headers[%s]", path, name), result)
		}
	}

	// Validate securitySchemes
	for name, ss := range c.SecuritySchemes {
		if !namePattern.MatchString(name) {
			result.addError(fmt.Sprintf("%s.securitySchemes[%s]", path, name), "component name contains invalid characters")
		}
		if ss != nil {
			ss.validate(fmt.Sprintf("%s.securitySchemes[%s]", path, name), result)
		}
	}

	// Validate links
	for name, link := range c.Links {
		if !namePattern.MatchString(name) {
			result.addError(fmt.Sprintf("%s.links[%s]", path, name), "component name contains invalid characters")
		}
		if link != nil {
			link.validate(fmt.Sprintf("%s.links[%s]", path, name), result)
		}
	}

	// Validate callbacks
	for name, cb := range c.Callbacks {
		if !namePattern.MatchString(name) {
			result.addError(fmt.Sprintf("%s.callbacks[%s]", path, name), "component name contains invalid characters")
		}
		if cb != nil {
			cb.validate(fmt.Sprintf("%s.callbacks[%s]", path, name), result)
		}
	}

	// Validate pathItems (OpenAPI 3.1 specific)
	for name, pathItem := range c.PathItems {
		if !namePattern.MatchString(name) {
			result.addError(fmt.Sprintf("%s.pathItems[%s]", path, name), "component name contains invalid characters")
		}
		if pathItem != nil {
			pathItem.validate(fmt.Sprintf("%s.pathItems[%s]", path, name), result)
		}
	}
}

func (ss *SecurityScheme) validate(path string, result *ValidationResult) {
	// Skip validation for references
	if ss.IsReference() {
		return
	}

	// Required: type
	if ss.Type == "" {
		result.addError(path+".type", "required field is missing")
	} else {
		validTypes := map[string]bool{"apiKey": true, "http": true, "mutualTLS": true, "oauth2": true, "openIdConnect": true}
		if !validTypes[ss.Type] {
			result.addError(path+".type", fmt.Sprintf("must be one of: apiKey, http, mutualTLS, oauth2, openIdConnect; got %s", ss.Type))
		}
	}

	// Type-specific requirements
	switch ss.Type {
	case "apiKey":
		if ss.Name == "" {
			result.addError(path+".name", "required for apiKey type")
		}
		if ss.In == "" {
			result.addError(path+".in", "required for apiKey type")
		} else {
			validIn := map[string]bool{"query": true, "header": true, "cookie": true}
			if !validIn[ss.In] {
				result.addError(path+".in", "must be one of: query, header, cookie")
			}
		}
	case "http":
		if ss.Scheme == "" {
			result.addError(path+".scheme", "required for http type")
		}
	case "oauth2":
		if ss.Flows == nil {
			result.addError(path+".flows", "required for oauth2 type")
		} else {
			ss.Flows.validate(path+".flows", result)
		}
	case "openIdConnect":
		if ss.OpenIdConnectUrl == "" {
			result.addError(path+".openIdConnectUrl", "required for openIdConnect type")
		}
	}
}

func (f *OAuthFlows) validate(path string, result *ValidationResult) {
	// At least one flow must be defined
	hasFlow := f.Implicit != nil || f.Password != nil || f.ClientCredentials != nil || f.AuthorizationCode != nil
	if !hasFlow {
		result.addError(path, "at least one OAuth flow must be defined")
	}

	if f.Implicit != nil {
		// Implicit requires authorizationUrl
		if f.Implicit.AuthorizationUrl == "" {
			result.addError(path+".implicit.authorizationUrl", "required for implicit flow")
		}
		if f.Implicit.Scopes == nil {
			result.addError(path+".implicit.scopes", "required field is missing")
		}
	}

	if f.Password != nil {
		// Password requires tokenUrl
		if f.Password.TokenUrl == "" {
			result.addError(path+".password.tokenUrl", "required for password flow")
		}
		if f.Password.Scopes == nil {
			result.addError(path+".password.scopes", "required field is missing")
		}
	}

	if f.ClientCredentials != nil {
		// ClientCredentials requires tokenUrl
		if f.ClientCredentials.TokenUrl == "" {
			result.addError(path+".clientCredentials.tokenUrl", "required for clientCredentials flow")
		}
		if f.ClientCredentials.Scopes == nil {
			result.addError(path+".clientCredentials.scopes", "required field is missing")
		}
	}

	if f.AuthorizationCode != nil {
		// AuthorizationCode requires both authorizationUrl and tokenUrl
		if f.AuthorizationCode.AuthorizationUrl == "" {
			result.addError(path+".authorizationCode.authorizationUrl", "required for authorizationCode flow")
		}
		if f.AuthorizationCode.TokenUrl == "" {
			result.addError(path+".authorizationCode.tokenUrl", "required for authorizationCode flow")
		}
		if f.AuthorizationCode.Scopes == nil {
			result.addError(path+".authorizationCode.scopes", "required field is missing")
		}
	}
}
