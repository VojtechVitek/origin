package validation

import (
	"regexp"

	baseapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	errs "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	. "github.com/GoogleCloudPlatform/kubernetes/pkg/api/validation"

	. "github.com/openshift/origin/pkg/template/api"
	"github.com/openshift/origin/pkg/template/generator"
)

var parameterNameExp = regexp.MustCompile(`^[a-zA-Z0-9\_]+$`)

// ValidateParameter tests if required fields in the Parameter are set.
func ValidateParameter(param *Parameter) (list errs.ErrorList) {
	if !parameterNameExp.MatchString(param.Name) {
		list = append(list, errs.NewInvalid("Parameter.Name", param.Name, nil))
	}
	if !generator.MatchesGeneratorExpression(param.Generate) {
		list = append(list, errs.NewInvalid("Parameter.Generate", param.Name, nil))
	}
	return
}

// ValidateTemplateConfig tests if required fields in the TemplateConfig are set.
func ValidateTemplateConfig(config *TemplateConfig) (list errs.ErrorList) {
	if config.ID == "" {
		list = append(list, errs.NewInvalid("Config.ID", "", nil))
	}
	for _, item := range config.Items {
		switch obj := item.Object.(type) {
		case *baseapi.ReplicationController:
			list = append(list, ValidateReplicationController(obj)...)
		case *baseapi.Pod:
			list = append(list, ValidatePod(obj)...)
		case *baseapi.Service:
			list = append(list, ValidateService(obj)...)
		default:
			//TODO print the Kind
			list = append(list, errs.NewInvalid("Config.Items", "", nil))
		}
	}
	for i := range config.Parameters {
		list = append(list, ValidateParameter(&config.Parameters[i])...)
	}
	return
}
