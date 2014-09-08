package template

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/apiserver"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"

	"github.com/openshift/origin/pkg/template/api"
	"github.com/openshift/origin/pkg/template/api/validation"
	. "github.com/openshift/origin/pkg/template/generator"
)

type Storage struct{}

// NewRESTStorage returns a new apiserver.RESTStorage
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
	template, ok := obj.(*api.TemplateConfig)
	if !ok {
		return nil, fmt.Errorf("Not a template config.")
	}
	if errs := validation.ValidateTemplateConfig(template); len(errs) > 0 {
		return nil, fmt.Errorf("Invalid template config: %#v", errs)
	}
	return apiserver.MakeAsync(func() (interface{}, error) {
		generators := map[string]Generator{
			"generate": NewExpressionValueGenerator(rand.New(rand.NewSource(time.Now().UnixNano()))),
		}
		processor := NewTemplateProcessor(generators)
		config, err := processor.Process(template)
		return config, err
	}), nil
}
