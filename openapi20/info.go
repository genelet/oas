// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi20

import "encoding/json"

// Info provides metadata about the API
type Info struct {
	Title          string         `json:"title"`
	Version        string         `json:"version"`
	Description    string         `json:"description,omitempty"`
	TermsOfService string         `json:"termsOfService,omitempty"`
	Contact        *Contact       `json:"contact,omitempty"`
	License        *License       `json:"license,omitempty"`
	Extensions     map[string]any `json:"-"`
}

var infoKnownFields = []string{
	"title", "version", "description", "termsOfService", "contact", "license",
}

type infoAlias Info

func (i *Info) UnmarshalJSON(data []byte) error {
	var alias infoAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*i = Info(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	i.Extensions = extractExtensions(raw, infoKnownFields)
	return nil
}

func (i Info) MarshalJSON() ([]byte, error) {
	alias := infoAlias(i)
	return marshalWithExtensions(&alias, i.Extensions)
}

// Contact contains contact information for the API
type Contact struct {
	Name       string         `json:"name,omitempty"`
	URL        string         `json:"url,omitempty"`
	Email      string         `json:"email,omitempty"`
	Extensions map[string]any `json:"-"`
}

var contactKnownFields = []string{"name", "url", "email"}

type contactAlias Contact

func (c *Contact) UnmarshalJSON(data []byte) error {
	var alias contactAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*c = Contact(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	c.Extensions = extractExtensions(raw, contactKnownFields)
	return nil
}

func (c Contact) MarshalJSON() ([]byte, error) {
	alias := contactAlias(c)
	return marshalWithExtensions(&alias, c.Extensions)
}

// License contains license information for the API
type License struct {
	Name       string         `json:"name"`
	URL        string         `json:"url,omitempty"`
	Extensions map[string]any `json:"-"`
}

var licenseKnownFields = []string{"name", "url"}

type licenseAlias License

func (l *License) UnmarshalJSON(data []byte) error {
	var alias licenseAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*l = License(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	l.Extensions = extractExtensions(raw, licenseKnownFields)
	return nil
}

func (l License) MarshalJSON() ([]byte, error) {
	alias := licenseAlias(l)
	return marshalWithExtensions(&alias, l.Extensions)
}

// ExternalDocumentation allows referencing an external resource for extended documentation
type ExternalDocumentation struct {
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url"`
	Extensions  map[string]any `json:"-"`
}

var externalDocsKnownFields = []string{"description", "url"}

type externalDocsAlias ExternalDocumentation

func (e *ExternalDocumentation) UnmarshalJSON(data []byte) error {
	var alias externalDocsAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*e = ExternalDocumentation(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Extensions = extractExtensions(raw, externalDocsKnownFields)
	return nil
}

func (e ExternalDocumentation) MarshalJSON() ([]byte, error) {
	alias := externalDocsAlias(e)
	return marshalWithExtensions(&alias, e.Extensions)
}
