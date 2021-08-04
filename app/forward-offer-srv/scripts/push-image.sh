#!/bin/bash
APP_NAME=$3
# 私人仓库地址，通过第一个参数传进来，可以灵活构建成目标仓库地址的镜像
#REGISTRY_HOST="901739638239.dkr.ecr.ap-northeast-1.amazonaws.com"
REGISTRY_HOST=$1
if [ X$REGISTRY_HOST = "X" ];then
    echo "error, registry host is empty."
    exit 1
fi
# 镜像名称，不能有大写字母
DOCKER_IMAGE_NAME="$APP_NAME"

# 版本tag，通过第二参数传进来，如果为空，默认为latest
TAG=$2
if [ X$TAG = "X" ];then
    TAG="latest"
fi

# 通过TAG来判断部署环境，如果TAG是v1.0.0形式，表示正式环境；如果TAG是release-1.0.0形式，表示预生产环境；其他为开发环境。
# 使用前必须在文件~/.aws/credentials添加拥有ecr读写的aws账号信息，内容如下：
# 上传镜像
docker push $REGISTRY_HOST:$DOCKER_IMAGE_NAME-$TAG
if [ "$?" != "0" ]; then
    echo "docker push $REGISTRY_HOST:$DOCKER_IMAGE_NAME-$TAG failed."
    exit 1
fi
echo "docker push $REGISTRY_HOST:$DOCKER_IMAGE_NAME-$TAG success."
