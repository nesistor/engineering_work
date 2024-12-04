variable "gcp_credentials_file" {
  description = "Path to the JSON service account credentials file for GCP"
  type        = string
  default     = "/home/karol/terraform/cred/jenkins-deployer-key.json"
}

variable "gcp_project" {
  description = "The Google Cloud project ID"
  type        = string
  default     = "my-microservices-app" # Replace with your project ID
}

variable "gcp_region" {
  description = "The Google Cloud region"
  type        = string
  default     = "us-central1"
}

variable "gcp_zone" {
  description = "The Google Cloud zone"
  type        = string
  default     = "us-central1-a"
}

variable "kms_keyring" {
  description = "The name of the KMS KeyRing"
  type        = string
  default     = "vault-keyring"
}

variable "kms_key" {
  description = "The name of the KMS Key"
  type        = string
  default     = "vault-init-key"
}
