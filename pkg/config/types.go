package config

import "github.com/GoogleCloudPlatform/kubernetes/pkg/api"

type Config struct {
	api.JSONBase `json:",inline" yaml:",inline"`
	Name         string          `yaml:"name" json:"name"`
	Description  string          `yaml:"description" json:"description"`
	Items        []api.APIObject `yaml:"items" json:"items"`
}

func init() {
	api.AddKnownTypes("",
		Config{},
	)

	api.AddKnownTypes("v1beta1",
		Config{},
	)
}
