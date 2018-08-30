/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package crd

import (
	"context"
	"fmt"
	"net/http"

	crewv1alpha1 "github.com/mengqiy/webhook-example/pkg/apis/crew/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Validator validates Pods
type Validator struct {
	Client  client.Client
	Decoder admission.Decoder
}

// Implement admission.Handler so the controller can handle admission request.
var _ admission.Handler = &Validator{}

// Validator admits a pod iff a specific annotation exists.
func (v *Validator) Handle(ctx context.Context, req admission.Request) admission.Response {
	fm := &crewv1alpha1.Firstmate{}

	err := v.Decoder.Decode(req, fm)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	allowed, reason, err := validateFirstMateFn(ctx, fm)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

func validateFirstMateFn(ctx context.Context, fm *crewv1alpha1.Firstmate) (bool, string, error) {
	v, ok := ctx.Value(admission.StringKey("foo")).(string)
	if !ok {
		return false, "",
			fmt.Errorf("the value associated with key %q is expected to be a string", v)
	}
	if fm.Spec.Foo != v {
		return false, "can't find desired value with spec.foo", nil
	}
	return true, "", nil
}
