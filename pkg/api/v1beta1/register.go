package v1beta1

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/openshift/origin/pkg/config"
	templateapi "github.com/openshift/origin/pkg/template/api"
)

func init() {
	runtime.AddKnownTypes("v1beta1",
		templateapi.TemplateConfig{},
		config.Config{},
	)
}
