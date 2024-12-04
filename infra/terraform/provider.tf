provider "google" {
  credentials = file(var.gcp_credentials_file) # Path to the JSON service account key file
  project     = var.gcp_project               # GCP project ID
  region      = var.gcp_region                # Region used in the project
  zone        = var.gcp_zone                  # Zone used in the project
}
