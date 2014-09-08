package template

import (
	"fmt"
	"regexp"
	"strings"

	kubeapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/golang/glog"

	config "github.com/openshift/origin/pkg/config/api"
	"github.com/openshift/origin/pkg/template/api"
	. "github.com/openshift/origin/pkg/template/generator"
)

var parameterExp = regexp.MustCompile(`\$\{([a-zA-Z0-9\_]+)\}`)

// TemplateProcessor transforms TemplateConfig objects into Config objects.
type TemplateProcessor struct {
	Generators map[string]Generator
}

// NewTemplateProcessor creates new TemplateProcessor and initializes
// it's set of generators.
func NewTemplateProcessor(generators map[string]Generator) *TemplateProcessor {
	return &TemplateProcessor{Generators: generators}
}

// Process transforms TemplateConfig object into Config object. It generates
// Parameter values using the defined set of generators first, and then it
// substitutes all Parameter expression occurances with their corresponding
// values (currently in the containers' Environment variables only).
func (p *TemplateProcessor) Process(template *api.TemplateConfig) (*config.Config, error) {
	if err := p.GenerateParameterValues(template); err != nil {
		return nil, err
	}
	if err := p.SubstituteParameters(template); err != nil {
		return nil, err
	}

	config := &config.Config{
		Name:        template.Name,
		Description: template.Description,
		Items:       template.Items,
	}
	config.ID = template.ID
	config.Kind = "Config"
	config.CreationTimestamp = util.Now()
	return config, nil
}

// AddParameter adds new custom parameter to the TemplateConfig. It overrides
// the existing parameter, if already defined.
func (p *TemplateProcessor) AddParameter(t *api.TemplateConfig, param api.Parameter) {
	if existing := p.GetParameterByName(t, param.Name); existing != nil {
		*existing = param
	} else {
		t.Parameters = append(t.Parameters, param)
	}
}

// GetParameterByName searches for a Parameter in the TemplateConfig
// based on it's name.
func (p *TemplateProcessor) GetParameterByName(t *api.TemplateConfig, name string) *api.Parameter {
	for i, param := range t.Parameters {
		if param.Name == name {
			return &(t.Parameters[i])
		}
	}
	return nil
}

// SubstituteParameters loops over all Environment variables defined for
// all ReplicationController and Pod containers and substitutes all
// Parameter expression occurances with their corresponding values.
//
// Example of Parameter expression:
//   - ${PARAMETER_NAME}
func (p *TemplateProcessor) SubstituteParameters(t *api.TemplateConfig) error {
	// Make searching for given parameter name/value more effective
	paramMap := make(map[string]string, len(t.Parameters))
	for _, param := range t.Parameters {
		paramMap[param.Name] = param.Value
	}

	// manifestSubstituteParameters is a helper function that iterates
	// over the given manifest.
	substituteParametersInManifest := func(manifest *kubeapi.ContainerManifest) {
		for i, _ := range manifest.Containers {
			for e, _ := range manifest.Containers[i].Env {
				envValue := &manifest.Containers[i].Env[e].Value
				// Match all parameter expressions found in the given env var
				for _, match := range parameterExp.FindAllStringSubmatch(*envValue, -1) {
					// Substitute expression with its value, if corresponding parameter found
					if len(match) > 1 {
						if paramValue, found := paramMap[match[1]]; found {
							*envValue = strings.Replace(*envValue, match[0], paramValue, 1)
						}
					}
				}
			}
		}
	}

	for i, item := range t.Items {
		switch obj := item.Object.(type) {
		case *kubeapi.ReplicationController:
			substituteParametersInManifest(&obj.DesiredState.PodTemplate.DesiredState.Manifest)
			t.Items[i] = runtime.Object{Object: *obj}
		case *kubeapi.Pod:
			substituteParametersInManifest(&obj.DesiredState.Manifest)
			t.Items[i] = runtime.Object{Object: *obj}
		default:
			glog.V(1).Infof("Unable to process parameters for resource '%T'.", obj)
		}
	}

	return nil
}

// GenerateParameterValues generates Value for each Parameter of the given
// TemplateConfig that has Expression field specified and doesn't have any
// Value yet.
//
// Examples (Expression => Value):
//   - "test[0-9]{1}x" => "test7x"
//   - "[0-1]{8}" => "01001100"
//   - "0x[A-F0-9]{4}" => "0xB3AF"
//   - "[a-zA-Z0-9]{8}" => "hW4yQU5i"
func (p *TemplateProcessor) GenerateParameterValues(t *api.TemplateConfig) error {
	for i, _ := range t.Parameters {
		param := &t.Parameters[i]
		if param.Expression != "" && param.Value == "" {
			generator, ok := p.Generators["expression"]
			if !ok {
				return fmt.Errorf("Can't find expression generator.")
			}
			value, err := generator.GenerateValue(param.Expression)
			if err != nil {
				return err
			}
			param.Value, ok = value.(string)
			if !ok {
				return fmt.Errorf("Can't convert the generated value %v to string.", value)
			}
		}
	}
	return nil
}
