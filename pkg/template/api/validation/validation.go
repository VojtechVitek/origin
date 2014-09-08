package validation

import (
	"fmt"
	"regexp"

	kubeapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	errs "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	. "github.com/GoogleCloudPlatform/kubernetes/pkg/api/validation"

	. "github.com/openshift/origin/pkg/template/api"
)

var parameterNameExp = regexp.MustCompile(`^[a-zA-Z0-9\_]+$`)

// ValidateParameter tests if required fields in the Parameter are set.
func ValidateParameter(param *Parameter) (allErrs errs.ErrorList) {
	if param.Name == "" {
		allErrs = append(allErrs, errs.NewFieldRequired("name", param.Name))
		return
	}
	if !parameterNameExp.MatchString(param.Name) {
		allErrs = append(allErrs, errs.NewFieldInvalid("name", param.Name))
	}
	return
}

// ValidateTemplateConfig tests if required fields in the TemplateConfig are set.
func ValidateTemplateConfig(config *TemplateConfig) (allErrs errs.ErrorList) {
	if config.ID == "" {
		allErrs = append(allErrs, errs.NewFieldRequired("id", config.ID))
	}
	for i, item := range config.Items {
		itemErr := errs.ErrorList{}
		switch obj := item.Object.(type) {
		case *kubeapi.ReplicationController:
			itemErr = ValidateReplicationController(obj)
		case *kubeapi.Pod:
			itemErr = ValidatePod(obj)
		case *kubeapi.Service:
			itemErr = ValidateService(obj)
		default:
			itemErr = append(itemErr, errs.NewFieldInvalid("kind", fmt.Sprintf("%T", item)))
		}
		allErrs = append(allErrs, itemErr.Prefix(fmt.Sprintf("items[%d]", i))...)
	}
	for i := range config.Parameters {
		paramErr := ValidateParameter(&config.Parameters[i])
		allErrs = append(allErrs, paramErr.Prefix(fmt.Sprintf("parameters[%d]", i))...)
	}
	return
}
