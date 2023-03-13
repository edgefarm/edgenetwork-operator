PWD := ${CURDIR}
NAME=edgenetwork-operator/controller
REGISTRY=ghcr.io/edgefarm
TAG?= latest

generate:
	go generate ./...

setup:
	kubectl apply -k https://github.com/metacontroller/metacontroller/manifests/production 
	
crd: generate
	kubectl apply -f ./manifests/crds/

controller:
	kubectl apply -f ./manifests/controller.yaml

install: crd controller

example:
	kubectl apply -f ./examples/network.yaml

build:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH="$(GOARCH)" go build -o main  -gcflags "all=-N -l" -ldflags '-extldflags "-static"' cmd/controller/main.go

image:
	docker build -t $(REGISTRY)/$(NAME):$(TAG) -f ./Dockerfile .

push: image
	docker push $(REGISTRY)/$(NAME):$(TAG)

clean:
	kubectl delete -f ./examples/network.yaml
	kubectl delete -f ./manifests/controller.yaml
	kubectl -n metacontroller delete pod/metacontroller-0

.PHONY: generate setup install example image push clean