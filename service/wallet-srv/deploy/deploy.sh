APP_VERSION=v1.0.0
export KUBECONFIG=/home/admin.conf
kubectl get no
sed -i "s/VERSION_NUMBER/${APP_VERSION}/g" service/wallet-srv/deploy/k8s-deployment.yml
kubectl apply -f service/wallet-srv/deploy/k8s-deployment.yml --namespace=develop
#sed -i "s/${APP_VERSION}/VERSION_NUMBER/g" k8s-deployment.yml