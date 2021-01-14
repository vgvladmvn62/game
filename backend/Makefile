.PHONY: dockerize clean build version publish

version=develop
imagename="empty"

build: bullseye

publish: dockerize
	docker login
	docker push $(imagename):$(version)

version:
	@echo $(imagename) $(version)

dockerize:
	docker build -t $(imagename):$(version) . -f deployments/Dockerfile

clean: 
	rm bullseye

bullseye: 
	CGO_ENABLED=0 GO111MODULE=on go build -o bullseye -mod=vendor

test:
	go test -v -race ./...
