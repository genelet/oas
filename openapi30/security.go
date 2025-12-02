// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi30

import "encoding/json"

// SecurityScheme defines a security scheme that can be used by the operations.
// It can also represent a Reference (when Ref is set).
type SecurityScheme struct {
	// Reference field
	Ref string `json:"$ref,omitempty"`

	// SecurityScheme fields
	Type             string         `json:"type,omitempty"` // apiKey, http, oauth2, openIdConnect
	Description      string         `json:"description,omitempty"`
	Name             string         `json:"name,omitempty"` // for apiKey
	In               string         `json:"in,omitempty"`   // for apiKey: query, header, cookie
	Scheme           string         `json:"scheme,omitempty"`
	BearerFormat     string         `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows    `json:"flows,omitempty"`
	OpenIdConnectUrl string         `json:"openIdConnectUrl,omitempty"`
	Extensions       map[string]any `json:"-"`
}

var securitySchemeKnownFields = []string{
	"$ref", "type", "description", "name", "in", "scheme", "bearerFormat", "flows", "openIdConnectUrl",
}

// IsReference checks if this security scheme is actually a reference ($ref)
func (ss *SecurityScheme) IsReference() bool {
	return ss != nil && ss.Ref != ""
}

// NewSecuritySchemeReference creates a security scheme that is actually a reference
func NewSecuritySchemeReference(ref string) *SecurityScheme {
	return &SecurityScheme{Ref: ref}
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

// OAuthFlows allows configuration of the supported OAuth Flows
type OAuthFlows struct {
	Implicit          *OAuthFlow     `json:"implicit,omitempty"`
	Password          *OAuthFlow     `json:"password,omitempty"`
	ClientCredentials *OAuthFlow     `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow     `json:"authorizationCode,omitempty"`
	Extensions        map[string]any `json:"-"`
}

var oauthFlowsKnownFields = []string{"implicit", "password", "clientCredentials", "authorizationCode"}

type oauthFlowsAlias OAuthFlows

func (of *OAuthFlows) UnmarshalJSON(data []byte) error {
	var alias oauthFlowsAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*of = OAuthFlows(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	of.Extensions = extractExtensions(raw, oauthFlowsKnownFields)
	return nil
}

func (of OAuthFlows) MarshalJSON() ([]byte, error) {
	alias := oauthFlowsAlias(of)
	return marshalWithExtensions(&alias, of.Extensions)
}

// OAuthFlow configuration details for a supported OAuth Flow
type OAuthFlow struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	RefreshUrl       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"`
	Extensions       map[string]any    `json:"-"`
}

var oauthFlowKnownFields = []string{"authorizationUrl", "tokenUrl", "refreshUrl", "scopes"}

type oauthFlowAlias OAuthFlow

func (of *OAuthFlow) UnmarshalJSON(data []byte) error {
	var alias oauthFlowAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*of = OAuthFlow(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	of.Extensions = extractExtensions(raw, oauthFlowKnownFields)
	return nil
}

func (of OAuthFlow) MarshalJSON() ([]byte, error) {
	alias := oauthFlowAlias(of)
	return marshalWithExtensions(&alias, of.Extensions)
}

// SecurityRequirement lists the required security schemes to execute an operation
type SecurityRequirement map[string][]string
