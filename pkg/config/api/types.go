package api

import (
	baseapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
)

type Config struct {
	baseapi.JSONBase `json:",inline" yaml:",inline"`
	Name             string           `yaml:"name" json:"name"`
	Description      string           `yaml:"description" json:"description"`
	Items            []runtime.Object `yaml:"items" json:"items"`
}
