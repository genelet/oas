// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import "encoding/json"

// Info provides metadata about the API
type Info struct {
	Title          string   `json:"title"`
	Summary        string   `json:"summary,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version"`
	Extensions     map[string]any `json:"-"`
}

var infoKnownFields = []string{
	"title", "summary", "description", "termsOfService", "contact", "license", "version",
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

// Contact information for the exposed API
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

// License information for the exposed API
type License struct {
	Name       string         `json:"name"`
	Identifier string         `json:"identifier,omitempty"`
	URL        string         `json:"url,omitempty"`
	Extensions map[string]any `json:"-"`
}

var licenseKnownFields = []string{"name", "identifier", "url"}

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
