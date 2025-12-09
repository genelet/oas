// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi20

import "encoding/json"

// SecurityScheme defines a security scheme in Swagger 2.0.
// Types can be: "basic", "apiKey", or "oauth2"
type SecurityScheme struct {
	// Type of security scheme: basic, apiKey, or oauth2
	Type string `json:"type"`

	// Description of the security scheme
	Description string `json:"description,omitempty"`

	// apiKey specific fields
	Name string `json:"name,omitempty"` // Name of the header or query parameter
	In   string `json:"in,omitempty"`   // Location: "header" or "query"

	// oauth2 specific fields
	Flow             string            `json:"flow,omitempty"`             // implicit, password, application, accessCode
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"` // Required for implicit and accessCode flows
	TokenUrl         string            `json:"tokenUrl,omitempty"`         // Required for password, application, and accessCode flows
	Scopes           map[string]string `json:"scopes,omitempty"`           // Available scopes for OAuth2

	Extensions map[string]any `json:"-"`
}

var securitySchemeKnownFields = []string{
	"type", "description", "name", "in", "flow", "authorizationUrl", "tokenUrl", "scopes",
}

type securitySchemeAlias SecurityScheme

func (ss *SecurityScheme) UnmarshalJSON(data []byte) error {
	var alias securitySchemeAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*ss = SecurityScheme(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	ss.Extensions = extractExtensions(raw, securitySchemeKnownFields)
	return nil
}

func (ss SecurityScheme) MarshalJSON() ([]byte, error) {
	alias := securitySchemeAlias(ss)
	return marshalWithExtensions(&alias, ss.Extensions)
}

// IsBasic returns true if this is a basic authentication scheme
func (ss *SecurityScheme) IsBasic() bool {
	return ss != nil && ss.Type == "basic"
}

// IsAPIKey returns true if this is an API key authentication scheme
func (ss *SecurityScheme) IsAPIKey() bool {
	return ss != nil && ss.Type == "apiKey"
}

// IsOAuth2 returns true if this is an OAuth2 authentication scheme
func (ss *SecurityScheme) IsOAuth2() bool {
	return ss != nil && ss.Type == "oauth2"
}

// SecurityRequirement lists the required security schemes to execute an operation.
// Each key is a security scheme name, and the value is a list of required scopes.
type SecurityRequirement map[string][]string
