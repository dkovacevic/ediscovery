variable "region" {
  description = "AWS region"
  type        = string
  default     = "eu-west-3"
}

variable "cluster_name" {
  description = "EKS Cluster name"
  type        = string
  default     = "barbara"
}

variable "app" {
  description = "Application name"
  type        = string
  default     = "lh-whatsapp"
}