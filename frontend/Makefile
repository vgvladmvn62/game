.PHONY: version dockerize publish

version=develop
imagename=bullseye-frontend

version:
	@echo $(imagename) $(version)

dockerize:
	docker build -t $(imagename):$(version) . -f Dockerfile

publish: dockerize
	docker login
	docker push $(imagename):$(version)
