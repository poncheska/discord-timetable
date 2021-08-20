.PHONY: build
build:
	docker build -t poncheska/dtt-bot -f docker/Dockerfile .
	docker push poncheska/dtt-bot