package template

import (
	"math/rand"
	"regexp"
	"strings"

	baseapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/golang/glog"
	"github.com/openshift/origin/pkg/template/api"
	"github.com/openshift/origin/pkg/template/generator"
)

var parameterExp = regexp.MustCompile(`\$\{([a-zA-Z0-9\_]+)\}`)

// AddCustomTemplateParameter allow to pass the custom parameter to the
// template. It will replace the existing parameter, when it is already
// defined in the template.
func AddCustomTemplateParameter(p api.Parameter, t *api.TemplateConfig) {
	if param := GetTemplateParameterByName(p.Name, t); param != nil {
		*param = p
	} else {
		t.Parameters = append(t.Parameters, p)
	}
}

// GetTemplateParameterByName will return the pointer to the Template
// parameter based on the Parameter name.
func GetTemplateParameterByName(name string, t *api.TemplateConfig) *api.Parameter {
	for i, param := range t.Parameters {
		if param.Name == name {
			return &(t.Parameters[i])
		}
	}
	return nil
}

// ProcessParameters searches the ReplicationControllers and Pods defined in the
// TemplateConfig and substitutes the references to parameters with the parameter
// values in the 'env' blocks defined for containers.
//
// Parameter expression example:
//   - ${PARAMETER_NAME}
func ProcessEnvParameters(t *api.TemplateConfig) error {
	// Make searching for given parameter name/value more effective
	paramMap := make(map[string]string, len(t.Parameters))
	for _, param := range t.Parameters {
		paramMap[param.Name] = param.Value
	}

	for i, item := range t.Items {
		switch obj := item.Object.(type) {
		case *baseapi.ReplicationController:
			subManifestParams(
				&obj.DesiredState.PodTemplate.DesiredState.Manifest,
				paramMap,
			)
			t.Items[i] = runtime.Object{Object: *obj}
		case *baseapi.Pod:
			subManifestParams(
				&obj.DesiredState.Manifest,
				paramMap,
			)
			t.Items[i] = runtime.Object{Object: *obj}
		default:
			glog.V(1).Infof("Unable to process parameters for resource '%T'.", obj)
		}
	}

	return nil
}

// GenerateParameterValue generates Value for each Parameter of the given
// Template that has Generate field specified and doesn't have any Value yet.
//
// Examples of what certain Generate field values generate:
//   - "test[0-9]{1}x" => "test7x"
//   - "[0-1]{8}" => "01001100"
//   - "0x[A-F0-9]{4}" => "0xB3AF"
//   - "[a-zA-Z0-9]{8}" => "hW4yQU5i"
//   - "password" => "hW4yQU5i"
//   - "[GET:http://api.example.com/generateRandomValue]" => remote string
func GenerateParameterValues(t *api.TemplateConfig, seed *rand.Rand) error {
	for i, _ := range t.Parameters {
		p := &t.Parameters[i]
		if p.Generate != "" && p.Value == "" {
			// Inherit the seed from parameter
			generator, err := generator.NewExpressionValueGenerator(seed)
			if err != nil {
				return err
			}
			value, err := generator.GenerateValue(p.Generate)
			if err != nil {
				return err
			}
			p.Value = value
		}
	}
	return nil
}

// subManifestParams is a helper method that iterates over any ContainerManifest
// object and search for the Env arrays.
// Then it will do the substitution of parameters in the Env values.
func subManifestParams(manifest *baseapi.ContainerManifest, params map[string]string) error {
	for i, _ := range manifest.Containers {
		for e, _ := range manifest.Containers[i].Env {
			envValue := &manifest.Containers[i].Env[e].Value
			// Match all parameter expressions found in the given env var
			for _, match := range parameterExp.FindAllStringSubmatch(*envValue, -1) {
				// Substitute expression with its value, if corresponding parameter found
				if len(match) > 1 {
					if paramValue, found := params[match[1]]; found {
						*envValue = strings.Replace(*envValue, match[0], paramValue, 1)
					}
				}
			}
		}
	}
	return nil
}
