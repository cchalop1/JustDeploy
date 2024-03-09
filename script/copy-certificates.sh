#!/bin/bash

PATH_CERT_DOCKER=/root/cert-docker

HOST=

scp root@$HOST:$PATH_CERT_DOCKER/ca.pem .
scp root@$HOST:$PATH_CERT_DOCKER/cert.pem .
scp root@$HOST:$PATH_CERT_DOCKER/key.pem .