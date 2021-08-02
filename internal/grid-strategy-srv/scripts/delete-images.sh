#!/bin/bash

APP_NAME=$1
# shellcheck disable=SC2046
ssh -tt root@172.28.119.42 << EOF
docker rmi -f $(docker images | grep "$APP_NAME" | awk '{print $3}')
exit
EOF

