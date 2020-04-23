BINARY := gleaner
VERSION :=`cat VERSION`
.DEFAULT_GOAL := gleaner
ORG := geocodes

gleaner:
	cd cmd/$(BINARY) ; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o $(BINARY)

gleaner.exe:
	cd cmd/$(BINARY) ; \
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 env go build -o $(BINARY).exe

gleaner.darwin:
	cd cmd/$(BINARY) ; \
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 env go build -o $(BINARY)_darwin

releases: gleaner gleaner.exe gleaner.darwin

glcon:
	cd cmd/glcon ; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o glcon

docker:
	docker build  --tag="$(ORG)/gleaner:$(VERSION)"  --file=./build/Dockerfile . ; \
	docker tag $(ORG)/gleaner:$(VERSION) $(ORG)/gleaner:latest

removeimage:
	docker rmi --force $(ORG)/gleaner:$(VERSION)
	docker rmi --force $(ORG)/gleaner:latest

publish: docker
	docker push $(ORG)/gleaner:$(VERSION) ; \
	docker push $(ORG)/gleaner:latest
