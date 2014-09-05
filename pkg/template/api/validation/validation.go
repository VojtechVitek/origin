package validation

import (
	"regexp"

	errs "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	//. "github.com/GoogleCloudPlatform/kubernetes/pkg/api/validation"

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
	/*for i := range config.Services {
		list = append(list, ValidateService(&config.Services[i])...)
	}
	for i := range config.Pods {
		list = append(list, ValidatePod(&config.Pods[i])...)
	}
	for i := range config.ReplicationControllers {
		list = append(list, ValidateReplicationController(&config.ReplicationControllers[i])...)
	}*/
	for i := range config.Parameters {
		list = append(list, ValidateParameter(&config.Parameters[i])...)
	}
	return
}
