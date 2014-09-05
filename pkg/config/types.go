package config

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
)

type Config struct {
	api.JSONBase `json:",inline" yaml:",inline"`
	Name         string           `yaml:"name" json:"name"`
	Description  string           `yaml:"description" json:"description"`
	Items        []runtime.Object `yaml:"items" json:"items"`
}

func init() {
	runtime.AddKnownTypes("v1beta1", Config{})
	runtime.AddKnownTypes("", Config{})
}
