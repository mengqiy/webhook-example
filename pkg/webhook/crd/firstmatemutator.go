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

// Mutator annotates Pods
type Mutator struct {
	Client  client.Client
	Decoder admission.Decoder
}

// Implement admission.Handler so the controller can handle admission request.
var _ admission.Handler = &Mutator{}

// Mutator changes a field in a CR.
func (a *Mutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	fm := &crewv1alpha1.Firstmate{}

	err := a.Decoder.Decode(req, fm)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	copy := fm.DeepCopy()

	err = mutateFirstMateFn(ctx, copy)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.PatchResponse(fm, copy)
}

// mutateFirstMateFn add an annotation to the given pod
func mutateFirstMateFn(ctx context.Context, fm *crewv1alpha1.Firstmate) error {
	v, ok := ctx.Value(admission.StringKey("foo")).(string)
	if !ok {
		return fmt.Errorf("the value associated with %v is expected to be a string", "foo")
	}
	fm.Spec.Foo = v
	return nil
}
