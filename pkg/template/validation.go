package template

import (
	"regexp"

	errs "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"

	. "github.com/openshift/origin/pkg/template/api"
	"github.com/openshift/origin/pkg/template/generator"
)

var parameterNameExp = regexp.MustCompile(`^[a-zA-Z0-9\_]+$`)

// ValidateParameter tests if required fields in the Parameter are set.
func ValidateParameter(param *Parameter) error {
	allErrs := errs.ErrorList{}
	if !parameterNameExp.MatchString(param.Name) {
		return errs.NewInvalid("Parameter.Name", param.Name, allErrs)
	}
	if !generator.MatchesGeneratorExpression(param.Generate) {
		return errs.NewInvalid("Parameter.Generate", param.Name, allErrs)
	}
	return nil
}

// ValidateTemplateConfig tests if required fields in the TemplateConfig are set.
func ValidateTemplateConfig(config *TemplateConfig) error {
	allErrs := errs.ErrorList{}
	if config.ID == "" {
		return errs.NewInvalid("Config.ID", "", allErrs)
	}
	/*
			for i := range config.Services {
				allErrs = append(allErrs, api.ValidateService(&config.Services[i])...)
			}
			for i := range config.Pods {
				allErrs = append(allErrs, api.ValidatePod(&config.Pods[i])...)
			}
			for i := range config.ReplicationControllers {
				allErrs = append(allErrs, api.ValidateReplicationController(&config.ReplicationControllers[i])...)
			}

		for i := range config.Parameters {
			allErrs = append(allErrs, ValidateParameter(&config.Parameters[i]))
		}
		if len(allErrs) > 0 {
			return errs.NewInvalid()
		}
	*/
	return nil
}
