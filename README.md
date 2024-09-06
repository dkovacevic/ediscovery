
## Add HELM Repository
```
$ helm repo add nginx-stable https://helm.nginx.com/stable
$ helm repo update
```

## Install the Nginx Controller
```
helm upgrade --install ingress-nginx ingress-nginx \
             --repo https://kubernetes.github.io/ingress-nginx \
             --namespace ingress-nginx --create-namespace
```

https://blog.saeloun.com/2023/03/21/setup-nginx-ingress-aws-eks/
