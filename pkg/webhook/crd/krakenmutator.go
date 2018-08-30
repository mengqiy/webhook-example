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

	creaturesv1alpha1 "github.com/mengqiy/example-crd-apis/pkg/apis/creatures/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// KrakenMutator annotates Pods
type KrakenMutator struct {
	Client  client.Client
	Decoder admission.Decoder
}

// Implement admission.Handler so the controller can handle admission request.
var _ admission.Handler = &Mutator{}

// Mutator changes a field in a CR.
func (a *KrakenMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	k := &creaturesv1alpha1.Kraken{}

	err := a.Decoder.Decode(req, k)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	copy := k.DeepCopy()

	err = mutateKrakenFn(ctx, copy)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.PatchResponse(k, copy)
}

// mutateKrakenFn add an annotation to the given pod
func mutateKrakenFn(ctx context.Context, k *creaturesv1alpha1.Kraken) error {
	v, ok := ctx.Value(admission.StringKey("foo")).(string)
	if !ok {
		return fmt.Errorf("the value associated with %v is expected to be a string", "foo")
	}
	anno := k.GetAnnotations()
	if anno == nil {
		anno = map[string]string{}
	}
	anno["foo"] = v
	k.SetAnnotations(anno)
	return nil
}
