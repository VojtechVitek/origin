package api

import (
	kubeapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
)

type TemplateConfig struct {
	kubeapi.JSONBase `json:",inline" yaml:",inline"`
	Name             string `yaml:"name" json:"name"`
	Description      string `yaml:"description" json:"description"`
	// TODO: This doesn't handle unregistered types. Define custom
	//       []interface{} type and it's unmarshaller instead.
	Items      []runtime.Object `yaml:"items" json:"items"`
	Parameters []Parameter      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Type        string `json:"type" yaml:"type"`
	Expression  string `json:"expression,omitempty" yaml:"expression,omitempty"`
	Value       string `json:"value,omitempty" yaml:"value,omitempty"`
}
