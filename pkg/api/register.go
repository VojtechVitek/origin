package api

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/openshift/origin/pkg/config"
	templateapi "github.com/openshift/origin/pkg/template/api"
)

func init() {
	runtime.AddKnownTypes("",
		templateapi.TemplateConfig{},
		config.Config{},
	)
}
