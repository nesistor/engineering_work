provider "google" {
  credentials = file(var.gcp_credentials_file) # Path to the JSON service account key file
  project     = var.gcp_project               # GCP project ID
  region      = var.gcp_region                # Region used in the project
  zone        = var.gcp_zone                  # Zone used in the project
}

variable "gcp_credentials_file" {
  description = "Path to the JSON service account credentials file for GCP"
  type        = string
  default     = "<path_to_your_service_account_key_file>" # Replace with the actual path
}

variable "gcp_project" {
  description = "Google Cloud project ID"
  type        = string
  default     = "my-microservices-app" # Your project
}

variable "gcp_region" {
  description = "Region used in the Google Cloud project"
  type        = string
  default     = "us-central1" # Default region
}

variable "gcp_zone" {
  description = "Zone used in the Google Cloud project"
  type        = string
  default     = "us-central1-a" # Default zone
}

output "gcp_project_id" {
  description = "Google Cloud project ID"
  value       = var.gcp_project
}

output "gcp_region" {
  description = "Google Cloud project region"
  value       = var.gcp_region
}

output "gcp_zone" {
  description = "Google Cloud project zone"
  value       = var.gcp_zone
}
