buildName=bullseye-rpi
version=develop

build:
	CGO_ENABLED=0 GO111MODULE=on go build -o $(buildName)-$(version)

clean: 
	rm $(buildName)-$(version)