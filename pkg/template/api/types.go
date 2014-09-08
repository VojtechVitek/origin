package api

import configapi "github.com/openshift/origin/pkg/config/api"

type TemplateConfig struct {
	configapi.Config `json:",inline" yaml:",inline"`
	Parameters       []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Type        string `json:"type" yaml:"type"`
	Expression  string `json:"expression,omitempty" yaml:"expression,omitempty"`
	Value       string `json:"value,omitempty" yaml:"value,omitempty"`
}
