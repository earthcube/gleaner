# Running Gleaner 


Note: This README was added as part of the ECO effort to get the 2.0.x 
version of Gleaner running.


## Requirements

Assumes Ubuntu LTS

* Docker/swarm
* Docker compose
* Go 

## Install Go

Install `go`, required to build `gleaner`:
```
sudo apt install golang-go
```

Modify `.bash_profile`, add GOPATH/bin
```
export PATH=$PATH:$(go env GOPATH)/bin
```

Install `dep` for dependency management:
```
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

## Setup Gleaner

```
mkdir -p go/src/earthcube.org/Project418
cd go/src/earthcube.org/Project418
git clone https://github.com/craig-willis/gleaner -b v2.0.4
```

Build gleaner binary
```
cd gleaner
dep ensure
make gleaner
```

Results:
```
cd cmd/gleaner ; \
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o gleaner
```

Binary is now in `cmd/gleaner`


## Run the starterpack

```
cd docs/starterpack
```

For this example, we'll run a local nginx server with a subset of the IEDA sitemap:
```
docker run --rm -d -p 9090:80 -v `pwd`/test:/usr/share/nginx/html/ nginx
```

Confirm that you can request the sitemap:
```
curl localhost:9090/ieda-sitemap.xml
```

Run the Gleaner stack:
```
. demo-exp.env
docker-compose -f gleaner-compose-ECAM.yml up
```

## Run Gleaner


```
cd ~/go/src/earthcube.org/Project418/gleaner
. docs/starterpack/demo-exp.env
cmd/gleaner/gleaner -setup
cmd/gleaner/gleaner -configfile docs/starterpack/config.yaml
```

Confirm spatial entries:
```
docker exec -it <tile38-container> sh
tile38-cli scan p418
```

