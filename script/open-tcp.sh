#!/bin/bash

PATH_CERT_DOCKER=/root/cert-docker

mkdir -p /etc/systemd/system/docker.service.d 

echo "[Service]
ExecStart=
ExecStart=/usr/bin/dockerd --tlsverify --tlscacert=$PATH_CERT_DOCKER/ca.pem --tlscert=$PATH_CERT_DOCKER/server-cert.pem --tlskey=$PATH_CERT_DOCKER/server-key.pem -H fd:// -H=0.0.0.0:2376" >> /etc/systemd/system/docker.service.d/override.conf

systemctl daemon-reload

systemctl restart docker.service