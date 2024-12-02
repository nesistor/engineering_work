provider "google" {
  credentials = file("<path_to_your_service_account_key_file>")
  project     = var.gcp_project
  region      = var.gcp_region
  zone        = var.gcp_zone
}
