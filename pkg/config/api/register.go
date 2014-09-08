package api

import "github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

func init() {
	runtime.AddKnownTypes("v1beta1", Config{})
	runtime.AddKnownTypes("", Config{})
}