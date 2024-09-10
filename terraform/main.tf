module "aws_eks" {
  source = "./modules/eks"

  cluster_name = var.cluster_name
  region = var.region
}

module "ingress_helm" {
  source = "./modules/helm"
}

module "lh_kubernetes" {
  source = "./modules/kubernetes"

  app = var.app
  cluster_name = var.cluster_name
}
