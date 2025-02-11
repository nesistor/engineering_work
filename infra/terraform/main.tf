resource "google_container_cluster" "gke_cluster" {
  name     = "my-gke-cluster"
  location = var.gcp_zone

  initial_node_count = 2

  node_config {
    machine_type  = "e2-medium"
    disk_size_gb  = 50          # Set disk size to 50GB
    disk_type     = "pd-standard" # Use HDD instead of SSD
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
    machine_type  = "e2-medium"
    disk_size_gb  = 50          # Set disk size to 50GB
    disk_type     = "pd-standard" # Use HDD instead of SSD
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
}
