resource "kubernetes_persistent_volume" "ediscovery_pv" {
  metadata {
    name = "ediscovery-pv"
  }

  spec {
    capacity = {
      storage = "1Gi"  # This should match the size of your EBS volume
    }

    access_modes = ["ReadWriteOnce"]

    persistent_volume_source {
      csi {
        driver = "ebs.csi.aws.com"
        volume_handle = "vol-063ac3727a0ca64e2"  # Referencing the EBS volume from data source
        fs_type   = "ext4"                          # Filesystem type on the EBS volume
      }
    }

    storage_class_name = "manual"
  }
}

resource "kubernetes_persistent_volume_claim" "ediscovery_pvc" {
  metadata {
    name = "ediscovery-pvc"
  }

  spec {
    access_modes = ["ReadWriteOnce"]

    resources {
      requests = {
        storage = "1Gi"  # Must match the PV size
      }
    }

    storage_class_name = "manual"
    volume_name        = "ediscovery-pv"
  }
}

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
        affinity {
          node_affinity {
            required_during_scheduling_ignored_during_execution {
              node_selector_term {
                match_expressions {
                  key = "topology.kubernetes.io/zone"
                  operator = "In"
                  values = ["us-east-1b"]
                }
              }
            }
          }
        }
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
            claim_name = "ediscovery-pvc"
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
