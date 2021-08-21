.PHONY: build
build:
	docker build -t poncheska/dtt-bot -f docker/Dockerfile .
	docker push poncheska/dtt-bot

.PHONY: apply
apply:
	kubectl apply -f deployment/namespace.yaml
	kubectl apply -f deployment/config.yaml
	kubectl apply -f deployment/deployment.yaml

.PHONY: delete
delete:
	kubectl delete -f deployment/deployment.yaml
	kubectl delete -f deployment/config.yaml
	kubectl delete -f deployment/namespace.yaml
