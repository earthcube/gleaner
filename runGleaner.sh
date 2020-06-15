#!/bin/bash

docker run  -v `pwd`/configs:/configs \
	-v /tmp:/tmp \
	-v /bleve:/bleve \
	-v /data/volumes:/DV \
	--net=host \
	--privileged=true \
	--entrypoint=/gleaner \
	--env-file=docs/starterpack/demo.env \
	geocodes/gleaner:latest \
	-configfile /configs/eco.yaml $1
