PWD := ${CURDIR}
NAME=yakha
TAG?= dev
REGISTRY=localhost:5000

setup:
	kubectl apply -k https://github.com/metacontroller/metacontroller/manifests/production 
install:
	cd hook/; go generate ./...; cd ..; 
	kubectl apply -f ./crds
	kubectl apply -f ./controller.yaml
example:
	kubectl apply -f ./examples/network.yaml

image:
	docker build -t $(NAME):$(TAG) -f ./hook/Dockerfile ./hook
push: image
	docker tag $(NAME):$(TAG) mmrxx/$(NAME):$(TAG)
	docker push mmrxx/$(NAME):$(TAG)
clean:
	kubectl delete -f ./examples/network.yaml
	kubectl delete -f ./controller.yaml
	kubectl -n metacontroller delete pod/metacontroller-0