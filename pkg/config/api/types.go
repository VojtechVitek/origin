package api

import (
	kubeapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
)

type Config struct {
	kubeapi.JSONBase `json:",inline" yaml:",inline"`
	Name             string `yaml:"name" json:"name"`
	Description      string `yaml:"description" json:"description"`
	// TODO: This doesn't handle unregistered types. Define custom
	//       []interface{} type and its unmarshaller instead.
	Items []runtime.Object `yaml:"items" json:"items"`
}
