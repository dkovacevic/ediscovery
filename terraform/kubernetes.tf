data "aws_eks_cluster" "cluster" {
  name = var.cluster_name
  depends_on = [
    module.eks]
}

provider "kubernetes" {
  host = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command = "aws"
    args = [
      "eks",
      "get-token",
      "--cluster-name",
      var.cluster_name
    ]
  }
}

resource "kubernetes_deployment" "lh-whatsapp" {
  metadata {
    name = var.app
    labels = {
      app = var.app
    }
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        app = var.app
      }
    }
    template {
      metadata {
        labels = {
          app = var.app
        }
      }
      spec {
        container {
          image = "536697232357.dkr.ecr.eu-west-3.amazonaws.com/orion:latest"
          name = var.app

          port {
            container_port = 8080
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "lh_whatsapp" {
  metadata {
    name = var.app
  }

  spec {
    selector = {
      app = var.app
    }

    type = "NodePort"

    port {
      name = "http"
      port = 80
      target_port = 8080
      protocol = "TCP"
    }
  }
}

resource "kubernetes_manifest" "letsencrypt_cluster_issuer" {
  manifest = {
    "apiVersion" = "cert-manager.io/v1"
    "kind" = "ClusterIssuer"
    "metadata" = {
      "name" = "letsencrypt"
    }
    "spec" = {
      "acme" = {
        "server" = "https://acme-v02.api.letsencrypt.org/directory"
        "email" = "dejankov@gmail.com"
        "privateKeySecretRef" = {
          "name" = "letsencrypt"
        }
        "solvers" = [
          {
            "http01" = {
              "ingress" = {
                "class" = "nginx"
              }
            }
          }]
      }
    }
  }
  depends_on = [module.eks]
}

resource "kubernetes_manifest" "orion_ingress" {
  manifest = {
    "apiVersion" = "networking.k8s.io/v1"
    "kind" = "Ingress"
    "metadata" = {
      "name" = "orion-ingress"
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
