package validation

import (
	"testing"

	kubeapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	"github.com/openshift/origin/pkg/config/api"
)

func TestValidateConfig(t *testing.T) {
	var tests = []struct {
		config          *api.Config
		isValidExpected bool
	}{
		{ // Config with empty ID, should fail
			&api.Config{},
			false,
		},
		{ // Config with ID, should pass
			&api.Config{
				JSONBase: kubeapi.JSONBase{ID: "ConfigId"},
			},
			true,
		},
		{ // Config with an item of unknown Kind, should pass
			&api.Config{
				JSONBase: kubeapi.JSONBase{ID: "templateId"},
				Items:    []runtime.EmbeddedObject{{}},
			},
			true,
		},
	}

	for _, test := range tests {
		errs := ValidateConfig(test.config)
		if len(errs) != 0 && test.isValidExpected {
			t.Errorf("Unexpected non-empty error list: %#v", errs)
		}
		if len(errs) == 0 && !test.isValidExpected {
			t.Errorf("Unexpected empty error list: %#v", errs)
		}
	}
}
