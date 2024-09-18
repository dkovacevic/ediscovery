resource "kubernetes_manifest" "ediscovery" {
  manifest = {
    "apiVersion" = "networking.k8s.io/v1"
    "kind" = "Ingress"
    "metadata" = {
      "name" = "ediscovery"
      "namespace" = "default"
      "annotations" = {
        "nginx.ingress.kubernetes.io/rewrite-target" = "/"
        "cert-manager.io/cluster-issuer" = "letsencrypt"
      }
    }
    "spec" = {
      "ingressClassName" = "nginx"
      "tls" = [
        {
          "secretName" = "letsencrypt"
          "hosts" = [
            "ediscovery.cz",
            "www.ediscovery.cz"
          ]
        }
      ]
      "defaultBackend" = {
        "service" = {
          "name" = var.app
          "port" = {
            "number" = 80
          }
        }
      }
      "rules" = [
        {
          "host" = "ediscovery.cz"
          "http" = {
            "paths" = [
              {
                "path" = "/"
                "pathType" = "Prefix"
                "backend" = {
                  "service" = {
                    "name" = var.app
                    "port" = {
                      "number" = 80
                    }
                  }
                }
              }
            ]
          }
        },
        {
          "host" = "www.ediscovery.cz"
          "http" = {
            "paths" = [
              {
                "path" = "/"
                "pathType" = "Prefix"
                "backend" = {
                  "service" = {
                    "name" = var.app
                    "port" = {
                      "number" = 80
                    }
                  }
                }
              }
            ]
          }
        }
      ]
    }
  }
}
