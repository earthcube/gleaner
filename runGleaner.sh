#!/bin/bash

docker run  -v `pwd`/configs:/configs \
	-v /tmp:/tmp \
	--net=host \
	--privileged=true \
	--entrypoint=/gleaner \
	--env-file=docs/starterpack/demo.env \
	geocodes/gleaner:latest \
	-configfile /configs/eco.yaml $1
