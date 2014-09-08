package template

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	_ "github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta1"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	"github.com/openshift/origin/pkg/template/api"
	. "github.com/openshift/origin/pkg/template/generator"
)

func TestNewTemplate(t *testing.T) {
	var template api.TemplateConfig

	jsonData, _ := ioutil.ReadFile("example/project.json")
	if err := json.Unmarshal(jsonData, &template); err != nil {
		t.Errorf("Unable to process the JSON template file: %v", err)
	}
}

func TestCustomParameter(t *testing.T) {
	var template api.TemplateConfig

	jsonData, _ := ioutil.ReadFile("example/project.json")
	json.Unmarshal(jsonData, &template)

	processor := NewTemplateProcessor(nil)
	processor.AddCustomTemplateParameter(api.Parameter{Name: "CUSTOM_PARAM", Value: "1"}, &template)
	processor.AddCustomTemplateParameter(api.Parameter{Name: "CUSTOM_PARAM", Value: "2"}, &template)

	if p := processor.GetTemplateParameterByName("CUSTOM_PARAM", &template); p == nil {
		t.Errorf("Unable to add a custom parameter to the template")
	} else {
		if p.Value != "2" {
			t.Errorf("Unable to replace the custom parameter value in template")
		}
	}
}

func ExampleProcessTemplateParameters() {
	var template api.TemplateConfig
	jsonData, _ := ioutil.ReadFile("example/project.json")
	json.Unmarshal(jsonData, &template)

	generators := map[string]Generator{
		"expression": NewExpressionValueGenerator(rand.New(rand.NewSource(1337))),
	}
	processor := NewTemplateProcessor(generators)

	// Define custom parameter for transformation:
	customParam := api.Parameter{Name: "CUSTOM_PARAM1", Value: "1"}
	processor.AddCustomTemplateParameter(customParam, &template)

	// Process the template config into the result config
	config, _ := processor.Process(&template)
	// Reset the timestamp
	config.CreationTimestamp = util.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)

	result, _ := json.Marshal(config)
	fmt.Println(string(result))
	// Output:
	// {"kind":"Config","id":"guestbook","creationTimestamp":"1980-01-01T00:00:00Z","name":"guestbook-example","description":"Example shows how to build a simple multi-tier application using Kubernetes and Docker","items":[{"kind":"Service","id":"frontend","creationTimestamp":null,"apiVersion":"v1beta1","port":5432,"selector":{"name":"frontend"},"containerPort":0},{"kind":"Service","id":"redismaster","creationTimestamp":null,"apiVersion":"v1beta1","port":10000,"selector":{"name":"redis-master"},"containerPort":0},{"kind":"Service","id":"redisslave","creationTimestamp":null,"apiVersion":"v1beta1","port":10001,"labels":{"name":"redisslave"},"selector":{"name":"redisslave"},"containerPort":0},{"kind":"Pod","id":"redis-master-2","creationTimestamp":null,"apiVersion":"v1beta1","labels":{"name":"redis-master"},"desiredState":{"manifest":{"version":"v1beta1","id":"redis-master-2","volumes":null,"containers":[{"name":"master","image":"dockerfile/redis","ports":[{"containerPort":6379}],"env":[{"name":"REDIS_PASSWORD","key":"REDIS_PASSWORD","value":"P8vxbV4C"}]}]},"restartpolicy":{}},"currentState":{"manifest":{"version":"","id":"","volumes":null,"containers":null},"restartpolicy":{}}},{"kind":"ReplicationController","id":"frontendController","creationTimestamp":null,"apiVersion":"v1beta1","desiredState":{"replicas":3,"replicaSelector":{"name":"frontend"},"podTemplate":{"desiredState":{"manifest":{"version":"v1beta1","id":"frontendController","volumes":null,"containers":[{"name":"php-redis","image":"brendanburns/php-redis","ports":[{"hostPort":8000,"containerPort":80}],"env":[{"name":"ADMIN_USERNAME","key":"ADMIN_USERNAME","value":"adminQ3H"},{"name":"ADMIN_PASSWORD","key":"ADMIN_PASSWORD","value":"dwNJiJwW"},{"name":"REDIS_PASSWORD","key":"REDIS_PASSWORD","value":"P8vxbV4C"}]}]},"restartpolicy":{}},"labels":{"name":"frontend"}}},"labels":{"name":"frontend"}},{"kind":"ReplicationController","id":"redisSlaveController","creationTimestamp":null,"apiVersion":"v1beta1","desiredState":{"replicas":2,"replicaSelector":{"name":"redisslave"},"podTemplate":{"desiredState":{"manifest":{"version":"v1beta1","id":"redisSlaveController","volumes":null,"containers":[{"name":"slave","image":"brendanburns/redis-slave","ports":[{"hostPort":6380,"containerPort":6379}],"env":[{"name":"REDIS_PASSWORD","key":"REDIS_PASSWORD","value":"P8vxbV4C"}]}]},"restartpolicy":{}},"labels":{"name":"redisslave"}}},"labels":{"name":"redisslave"}}]}
}
