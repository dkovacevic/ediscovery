resource "kubernetes_deployment" "ediscovery" {
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
          image = "536697232357.dkr.ecr.us-east-1.amazonaws.com/ediscovery:latest"
          name = var.app

          port {
            container_port = 8080
          }
          volume_mount {
            name       = "ediscovery-storage"
            mount_path = "/opt/ediscovery/data"
          }
        }
        volume {
          name = "ediscovery-storage"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.ediscovery_pvc.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "ediscovery" {
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
