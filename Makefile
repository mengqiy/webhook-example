# Image URL to use all building/pushing image targets
IMG ?= gcr.io/yourprojectname/webhook:v1

all: webhook

# Build manager binary
webhook: fmt vet
	go build -o bin/webhook github.com/mengqiy/webhook-example/example/

# Run go fmt against code
fmt:
	go fmt ./example/...

# Run go vet against code
vet:
	go vet ./example/...

# Install the webhook server and grants permissions
install:
	kubectl apply -f manifests/

# Build the docker image
docker-build:
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

clean:
	./hack/cleanup.sh
