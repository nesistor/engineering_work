provider "google" {
  project = "your-gcp-project-id"
  region  = "us-central1"
}

provider "kubernetes" {
  host                   = google_container_cluster.gke.endpoint
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(google_container_cluster.gke.master_auth[0].cluster_ca_certificate)
}

data "google_client_config" "default" {}

resource "google_container_cluster" "gke" {
  name     = "my-gke-cluster"
  location = "us-central1"
  initial_node_count = 3

  node_config {
    machine_type = "e2-medium"
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

resource "google_secret_manager_secret" "api_key" {
  secret_id = "auth-service-api-key"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "api_key_version" {
  secret = google_secret_manager_secret.api_key.id
  secret_data = "your-secret-api-key"
}

resource "google_project_iam_binding" "secret_access" {
  role    = "roles/secretmanager.secretAccessor"
  members = ["serviceAccount:my-service-account@your-gcp-project-id.iam.gserviceaccount.com"]
}
