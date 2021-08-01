#!/bin/bash

APP_NAME=$3
configPath=$4

# 私人仓库地址，通过第一个参数传进来，可以灵活构建成目标仓库地址的镜像
#REGISTRY_HOST="901739638239.dkr.ecr.ap-northeast-1.amazonaws.com"
REGISTRY_HOST=$1
if [ X$REGISTRY_HOST = "X" ];then
        echo "error, registry host is empty."
        exit 1
fi

# 镜像名称，不能有大写字母
DOCKER_IMAGE_NAME="$APP_NAME"

# 版本tag，通过第二个参数传进来，如果为空，默认为latest
TAG=$2
if [ X$TAG = "X" ];then
        TAG="latest"
fi

# Dockerfile文件位置
SCRIPT_PATH="./scripts"

# 配置文件位置
cp ./wallet-srv $SCRIPT_PATH/wallet-srv
cp ./config/conf.yaml $SCRIPT_PATH/conf.yaml
cp ./config/confPro.yaml $SCRIPT_PATH/confPro.yaml

# 基于Dockerfile所在的目录构建镜像
echo "docker build -t $REGISTRY_HOST:$DOCKER_IMAGE_NAME-$TAG $SCRIPT_PATH"
docker build --build-arg configPath=$configPath -t  $REGISTRY_HOST:$DOCKER_IMAGE_NAME-$TAG $SCRIPT_PATH

rm $SCRIPT_PATH/wallet-srv
rm $SCRIPT_PATH/conf.yaml
rm $SCRIPT_PATH/confPro.yaml