BINARY := gleaner
VERSION :=`cat VERSION`
.DEFAULT_GOAL := gleaner

gleaner:
	cd cmd/$(BINARY) ; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o $(BINARY)

glcon:
	cd cmd/glcon ; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o glcon

docker:
	docker build  --tag="nsfearthcube/gleaner:$(VERSION)"  --file=./build/Dockerfile . ; \
	docker tag nsfearthcube/gleaner:$(VERSION) nsfearthcube/gleaner:latest

removeimage:
	docker rmi --force nsfearthcube/gleaner:$(VERSION)
	docker rmi --force nsfearthcube/gleaner:latest

publish: docker
	docker push nsfearthcube/gleaner:$(VERSION) ; \
	docker push nsfearthcube/gleaner:latest