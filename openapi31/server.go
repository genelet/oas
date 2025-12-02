// Copyright (c) 2025 Greetingland LLC
// Created with the help of Claude Code
// MIT License - see LICENSE file for details

package openapi31

import "encoding/json"

// Server represents a server
type Server struct {
	URL         string                     `json:"url"`
	Description string                     `json:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty"`
	Extensions  map[string]any             `json:"-"`
}

var serverKnownFields = []string{"url", "description", "variables"}

type serverAlias Server

func (s *Server) UnmarshalJSON(data []byte) error {
	var alias serverAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*s = Server(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	s.Extensions = extractExtensions(raw, serverKnownFields)
	return nil
}

func (s Server) MarshalJSON() ([]byte, error) {
	alias := serverAlias(s)
	return marshalWithExtensions(&alias, s.Extensions)
}

// ServerVariable represents a server variable for server URL template substitution
type ServerVariable struct {
	Enum        []string       `json:"enum,omitempty"`
	Default     string         `json:"default"`
	Description string         `json:"description,omitempty"`
	Extensions  map[string]any `json:"-"`
}

var serverVariableKnownFields = []string{"enum", "default", "description"}

type serverVariableAlias ServerVariable

func (sv *ServerVariable) UnmarshalJSON(data []byte) error {
	var alias serverVariableAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*sv = ServerVariable(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	sv.Extensions = extractExtensions(raw, serverVariableKnownFields)
	return nil
}

func (sv ServerVariable) MarshalJSON() ([]byte, error) {
	alias := serverVariableAlias(sv)
	return marshalWithExtensions(&alias, sv.Extensions)
}
