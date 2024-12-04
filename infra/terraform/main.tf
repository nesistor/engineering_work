resource "google_container_cluster" "gke_cluster" {
  name     = "my-gke-cluster"
  location = var.gcp_zone

  initial_node_count = 2

  node_config {
    machine_type = "e2-medium"
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
}

resource "google_container_node_pool" "primary_node_pool" {
  name       = "primary-node-pool"
  cluster    = google_container_cluster.gke_cluster.name
  location   = var.gcp_zone
  node_count = 2

  node_config {
    machine_type = "e2-medium"
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
}

resource "google_service_account" "sa" {
  account_id   = "my-service-account"
  display_name = "My Service Account"
  project      = var.gcp_project
}

# KMS Module
module "kms" {
  source      = "./kms"        # Ścieżka do modułu KMS
  kms_keyring = var.kms_keyring
  kms_key     = var.kms_key
}
