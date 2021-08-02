APP_VERSION=v1.0.0
export KUBECONFIG=/var/jenkins_home/admin.conf
sed -i "s/VERSION_NUMBER/${APP_VERSION}/g" service/exchange-srv/deploy/k8s-deployment.yml
kubectl apply -f service/exchange-srv/deploy/k8s-deployment.yml --namespace=develop
#sed -i "s/${APP_VERSION}/VERSION_NUMBER/g" k8s-deployment.yml