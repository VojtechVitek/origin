package api

import "github.com/openshift/origin/pkg/config"

type TemplateConfig struct {
	Parameters    []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	config.Config `json:",inline" yaml:",inline"`
}

type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Type        string `json:"type" yaml:"type"`
	Generate    string `json:"generate,omitempty" yaml:"generate,omitempty"`
	Value       string `json:"value,omitempty" yaml:"value,omitempty"`
}
