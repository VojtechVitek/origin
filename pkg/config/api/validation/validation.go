package validation

import (
	"strings"

	kubeapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	kubevalidation "github.com/GoogleCloudPlatform/kubernetes/pkg/api/validation"

	"github.com/openshift/origin/pkg/config/api"
	deployapi "github.com/openshift/origin/pkg/deploy/api"
	deployvalidation "github.com/openshift/origin/pkg/deploy/api/validation"
	routeapi "github.com/openshift/origin/pkg/route/api"
	routevalidation "github.com/openshift/origin/pkg/route/api/validation"
)

// ValidateConfig tests if required fields in the Config are set.
func ValidateConfig(config *api.Config) (errs errors.ErrorList) {
	if len(config.ID) == 0 {
		errs = append(errs, errors.NewFieldRequired("id", config.ID))
	}
	for i, item := range config.Items {
		err := errors.ErrorList{}
		switch obj := item.Object.(type) {
		case *kubeapi.ReplicationController:
			err = kubevalidation.ValidateReplicationController(obj)
		case *kubeapi.Pod:
			err = kubevalidation.ValidatePod(obj)
		case *kubeapi.Service:
			err = kubevalidation.ValidateService(obj)
		case *routeapi.Route:
			err = routevalidation.ValidateRoute(obj)
		case *deployapi.Deployment:
			err = deployvalidation.ValidateDeployment(obj)
		case *deployapi.DeploymentConfig:
			err = deployvalidation.ValidateDeploymentConfig(obj)
		default:
			// Pass-through unknown types.
		}
		// ignore namespace validation errors in Config
		err = filter(err, "namespace")
		errs = append(errs, err.PrefixIndex(i).Prefix("items")...)
	}
	return
}

func filter(errs errors.ErrorList, prefix string) errors.ErrorList {
	if errs == nil {
		return errs
	}
	next := errors.ErrorList{}
	for _, err := range errs {
		ve, ok := err.(errors.ValidationError)
		if ok && strings.HasPrefix(ve.Field, prefix) {
			continue
		}
		next = append(next, err)
	}
	return next
}
