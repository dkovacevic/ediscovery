resource "kubernetes_storage_class" "ediscovery_storage" {
  metadata {
    name = "${var.app}-storage-class"
  }
  storage_provisioner = "kubernetes.io/aws-ebs"
  parameters = {
    type = "gp2"
  }
  reclaim_policy = "Retain"
}

resource "kubernetes_persistent_volume_claim" "ediscovery_pvc" {
  metadata {
    name = "${var.app}-pvc"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    storage_class_name = kubernetes_storage_class.ediscovery_storage.metadata[0].name
    resources {
      requests = {
        storage = "1Gi"
      }
    }
  }
}
