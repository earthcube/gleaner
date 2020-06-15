#!/bin/bash 
export MINIO_ACCESS_KEY=MySecretAccessKey
export MINIO_SECRET_KEY=MySecretSecretKeyforMinio
export DATAVOL=/data/volumes

cmd/gleaner/gleaner -configfile configs/unidata_config.yaml

