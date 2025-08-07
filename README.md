## Build
```agsl
docker buildx build --platform linux/amd64 -t dejankovacevic/ediscovery:0.2.4 .
```

## Add HELM Repository
```shell
helm repo add nginx-stable https://helm.nginx.com/stable
helm repo update
```

## Install the Nginx Controller
```shell
helm upgrade --install ingress-nginx ingress-nginx \
             --repo https://kubernetes.github.io/ingress-nginx \
             --namespace ingress-nginx --create-namespace
```

https://blog.saeloun.com/2023/03/21/setup-nginx-ingress-aws-eks/

## Add the Helm repository
```shell
helm repo add jetstack https://charts.jetstack.io --force-update
```

## Install cert-manager
```shell
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.15.3 \
  --set crds.enabled=true
```