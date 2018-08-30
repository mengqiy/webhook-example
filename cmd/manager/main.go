/*

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

package main

import (
	"log"
	"os"

	"github.com/mengqiy/webhook-example/pkg/webhook/crd"

	externalapis "github.com/mengqiy/example-crd-apis/pkg/apis"
	creaturesv1alpha1 "github.com/mengqiy/example-crd-apis/pkg/apis/creatures/v1alpha1"
	"github.com/mengqiy/webhook-example/pkg/apis"
	crewv1alpha1 "github.com/mengqiy/webhook-example/pkg/apis/crew/v1alpha1"
	"github.com/mengqiy/webhook-example/pkg/webhook/pod"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
	"sigs.k8s.io/controller-runtime/pkg/webhook/types"
)

func main() {

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Fatal(err)
	}
	if err := externalapis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Fatal(err)
	}

	//// Setup all Controllers
	//if err := controller.AddToManager(mgr); err != nil {
	//	log.Fatal(err)
	//}

	// Setup webhooks
	mutatingPodWebhook, err := builder.NewWebhookBuilder().
		Name("mutatingpods.k8s.io").
		Type(types.WebhookTypeMutating).
		Path("/mutating-pods").
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(mgr).
		ForType(&corev1.Pod{}).
		Build(&pod.Mutator{Client: mgr.GetClient(), Decoder: mgr.GetAdmissionDecoder()})
	if err != nil {
		log.Fatalf("unable to setup mutating webhook: %v", err)
	}

	validatingPodWebhook, err := builder.NewWebhookBuilder().
		Name("validatingpods.k8s.io").
		Type(types.WebhookTypeValidating).
		Path("/validating-pods").
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(mgr).
		ForType(&corev1.Pod{}).
		Build(&pod.Validator{Client: mgr.GetClient(), Decoder: mgr.GetAdmissionDecoder()})
	if err != nil {
		log.Fatalf("unable to setup validating webhook: %v", err)
	}

	mutatingFirstmateWebhook, err := builder.NewWebhookBuilder().
		Name("mutatingfirstmates.k8s.io").
		Type(types.WebhookTypeMutating).
		Path("/mutating-firstmates").
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(mgr).
		ForType(&crewv1alpha1.Firstmate{}).
		Build(&crd.Mutator{Client: mgr.GetClient(), Decoder: mgr.GetAdmissionDecoder()})
	if err != nil {
		log.Fatalf("unable to setup mutating webhook: %v", err)
	}

	validatingFirstmateWebhook, err := builder.NewWebhookBuilder().
		Name("validatingfirstmates.k8s.io").
		Type(types.WebhookTypeValidating).
		Path("/validating-firstmates").
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(mgr).
		ForType(&crewv1alpha1.Firstmate{}).
		Build(&crd.Validator{Client: mgr.GetClient(), Decoder: mgr.GetAdmissionDecoder()})
	if err != nil {
		log.Fatalf("unable to setup validating webhook: %v", err)
	}

	mutatingKrakenWebhook, err := builder.NewWebhookBuilder().
		Name("mutatingkraken.k8s.io").
		Type(types.WebhookTypeMutating).
		Path("/mutating-kraken").
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(mgr).
		ForType(&creaturesv1alpha1.Kraken{}).
		Build(&crd.KrakenMutator{Client: mgr.GetClient(), Decoder: mgr.GetAdmissionDecoder()})
	if err != nil {
		log.Fatalf("unable to setup validating webhook: %v", err)
	}

	as, err := webhook.NewServer("foo-admission-server", mgr, webhook.ServerOptions{
		Port:    9876,
		CertDir: "/tmp/cert",
		KVMap:   map[string]interface{}{"foo": "bar"},
		BootstrapOptions: &webhook.BootstrapOptions{
			Secret: &apitypes.NamespacedName{
				Namespace: "default",
				Name:      "foo-admission-server-secret",
			},

			Service: &webhook.Service{
				Namespace: "default",
				Name:      "foo-admission-server-service",
				// Selectors should select the pods that runs this webhook server.
				Selectors: map[string]string{
					"app": "foo-admission-server",
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("unable to create a new webhook server: %v", err)
		os.Exit(1)
	}
	err = as.Register(
		mutatingPodWebhook,
		validatingPodWebhook,
		mutatingFirstmateWebhook,
		validatingFirstmateWebhook,
		mutatingKrakenWebhook)
	if err != nil {
		log.Fatalf("unable to register webhooks in the admission server: %v", err)
	}

	log.Printf("Starting the Cmd.")

	// Start the Cmd
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
