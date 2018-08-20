# Webhook Example

This repo includes an example about how to implement an admission
webhook using the library in the `sigs.k8s.io/controller-runtime`.
This can be an prototype of how this can be integrated with `kubebuilder`

## Repo structure

`example` contains some example code of how to implement an admission
webhook using the library in the `controller-runtime` repo.

`manifests` contains some k8s resources for deploying the webhook
in your cluster. You can run `kubectl apply -f mainifests/` to
deploy them.

`Dockerfile` knows how to build an docker image for the webhook server.

`Makefile` can let you build and push your image; deploy the webhook.
Set the `IMG` environment variable before building or pushing your image.

## Workflow

- make docker-build
- make docker-push
- make install
- try to create a deployment `kubectl create deployment nginx --image=nginx`
- make clean
