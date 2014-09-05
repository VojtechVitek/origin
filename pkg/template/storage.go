package template

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/apiserver"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	configapi "github.com/openshift/origin/pkg/config/api"
	"github.com/openshift/origin/pkg/template/api"
	"github.com/openshift/origin/pkg/template/api/validation"
)

type Storage struct{}

// NewStorage returns a new apiserver.RESTStorage
// for the given TemplateConfig.
func NewStorage() apiserver.RESTStorage {
	return &Storage{}
}

func (s *Storage) List(selector labels.Selector) (interface{}, error) {
	return nil, fmt.Errorf("TemplateConfig can only be created.")
}

func (s *Storage) Get(id string) (interface{}, error) {
	return nil, fmt.Errorf("TemplateConfig can only be created.")
}

func (s *Storage) New() interface{} {
	return &api.TemplateConfig{}
}

func (s *Storage) Delete(id string) (<-chan interface{}, error) {
	return apiserver.MakeAsync(func() (interface{}, error) {
		return nil, fmt.Errorf("TemplateConfig can only be created.")
	}), nil
}

func (s *Storage) Update(minion interface{}) (<-chan interface{}, error) {
	return nil, fmt.Errorf("TemplateConfig can only be created.")
}

func (s *Storage) Create(obj interface{}) (<-chan interface{}, error) {
	t, ok := obj.(*api.TemplateConfig)
	if !ok {
		return nil, fmt.Errorf("Not a template config.")
	}
	if errs := validation.ValidateTemplateConfig(t); len(errs) > 0 {
		return nil, fmt.Errorf("Invalid template config: %#v", errs)
	}
	return apiserver.MakeAsync(func() (interface{}, error) {
		GenerateParameterValues(t, rand.New(rand.NewSource(time.Now().UnixNano())))
		err := ProcessEnvParameters(t)
		return s.toApiConfig(*t), err
	}), nil
}

func (s *Storage) toApiConfig(t api.TemplateConfig) configapi.Config {
	config := configapi.Config{
		Name:        t.Name,
		Description: t.Description,
		Items:       t.Items,
	}
	config.CreationTimestamp = util.Now()
	return config
}
