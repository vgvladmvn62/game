buildName=bullseye-rpi
version=develop
destination=destination

build:
	CGO_ENABLED=0 GO111MODULE=on go build -o $(buildName)-$(version)

rasp:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o $(buildName)-$(version) *.go

move:
	scp $(buildName)-$(version) $(destination)

clean: 
	rm $(buildName)-$(version)