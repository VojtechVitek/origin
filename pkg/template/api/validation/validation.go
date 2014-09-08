package validation

import (
	"fmt"
	"regexp"

	kubeapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	. "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	. "github.com/GoogleCloudPlatform/kubernetes/pkg/api/validation"

	. "github.com/openshift/origin/pkg/template/api"
)

var parameterNameExp = regexp.MustCompile(`^[a-zA-Z0-9\_]+$`)

// ValidateParameter tests if required fields in the Parameter are set.
func ValidateParameter(param *Parameter) (errs ErrorList) {
	if param.Name == "" {
		errs = append(errs, NewFieldRequired("name", param.Name))
		return
	}
	if !parameterNameExp.MatchString(param.Name) {
		errs = append(errs, NewFieldInvalid("name", param.Name))
	}
	return
}

// ValidateTemplateConfig tests if required fields in the TemplateConfig are set.
func ValidateTemplateConfig(config *TemplateConfig) (errs ErrorList) {
	if config.ID == "" {
		errs = append(errs, NewFieldRequired("id", config.ID))
	}
	for i, item := range config.Items {
		err := ErrorList{}
		switch obj := item.Object.(type) {
		case *kubeapi.ReplicationController:
			err = ValidateReplicationController(obj)
		case *kubeapi.Pod:
			err = ValidatePod(obj)
		case *kubeapi.Service:
			err = ValidateService(obj)
		default:
			err = append(err, NewFieldInvalid("kind", fmt.Sprintf("%T", item)))
		}
		errs = append(errs, err.Prefix(fmt.Sprintf("items[%d]", i))...)
	}
	for i := range config.Parameters {
		paramErr := ValidateParameter(&config.Parameters[i])
		errs = append(errs, paramErr.Prefix(fmt.Sprintf("parameters[%d]", i))...)
	}
	return
}
