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

  # Optional: Enabling GKE node pool auto-scaling
  # enable_autoscaling = true
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
