APP_VERSION=v1.0.0
HARBOR_HOST=192.168.3.30:8086
HARBOR_ADDR=${HARBOR_HOST}/mateforce
DOCKER_IMAGE=api-gateway

echo ${APP_VERSION}
docker login -u admin -p QQabc123++ HARBOR_HOST
docker build -t ${HARBOR_ADDR}/${DOCKER_IMAGE}:${APP_VERSION} -f ./deployment/Dockerfile .
docker push ${HARBOR_ADDR}/${DOCKER_IMAGE}:${APP_VERSION}
docker rmi ${HARBOR_ADDR}/${DOCKER_IMAGE}:${APP_VERSION} -f

sed -i "s/VERSION_NUMBER/${APP_VERSION}/g" deployment/k8s-deployment.yml
kubectl apply -f deployment/k8s-deployment.yml --namespace=develop
sed -i "s/${APP_VERSION}/VERSION_NUMBER/g" deployment/k8s-deployment.yml

# docker login -u admin -p QQabc123++  harbor.mateforce.vip:10087
# echo '{ "insecure-registries":["harbor.win,.com"]}' > /etc/docker/daemon.json

#sed -i "s/VERSION_NUMBER/${APP_VERSION}/g" k8s-d.yaml
#kubectl apply -f k8s-d.yml --namespace=develop
#sed -i "s/${APP_VERSION}/VERSION_NUMBER/g" k8s-d.yaml