#!/bin/bash

REGISTRY_HOST=$1
if [ X$REGISTRY_HOST = "X" ];then
        echo "error, registry host is empty."
        exit 1
fi

# 镜像名称，不能有大写字母
DOCKER_IMAGE_NAME=$2
SERVER_IP=$3
HOST="root@$SERVER_IP"
APP_NAME=$4

REGISTRY_HOST="$REGISTRY_HOST:$DOCKER_IMAGE_NAME"
# shellcheck disable=SC2087

ssh -tt "$HOST" << EOF
	docker stop -t 60 \$(docker ps -qa --no-trunc --filter name=$APP_NAME*);
	docker rm -f \$(docker ps -qa --no-trunc --filter name=$APP_NAME*);
	docker rmi --force $(docker images | grep "$APP_NAME" | awk "{print \$3}");

  docker pull "$REGISTRY_HOST";
  docker run -d --restart=always \
  --log-opt tag="$DOCKER_IMAGE_NAME" \
  --name "$DOCKER_IMAGE_NAME" "$REGISTRY_HOST"
exit
EOF
