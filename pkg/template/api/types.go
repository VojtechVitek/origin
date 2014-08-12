package api

import (
	baseapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/openshift/origin/pkg/config"
)

type TemplateConfig struct {
	baseapi.JSONBase `json:",inline" yaml:",inline"`
	Name             string      `yaml:"name" json:"name"`
	Description      string      `yaml:"description" json:"description"`
	Parameters       []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	config.Config    `json:",inline" yaml:",inline"`
}

type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Type        string `json:"type" yaml:"type"`
	Generate    string `json:"generate,omitempty" yaml:"generate,omitempty"`
	Value       string `json:"value,omitempty" yaml:"value,omitempty"`
}

func init() {
	baseapi.AddKnownTypes("",
		TemplateConfig{},
	)

	baseapi.AddKnownTypes("v1beta1",
		TemplateConfig{},
	)
}
