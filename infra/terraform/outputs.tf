output "gke_cluster_name" {
  description = "The name of the GKE cluster"
  value       = google_container_cluster.primary.name
}

output "gke_cluster_endpoint" {
  description = "The endpoint of the GKE cluster"
  value       = google_container_cluster.primary.endpoint
}

output "gke_cluster_zone" {
  description = "The zone of the GKE cluster"
  value       = google_container_cluster.primary.zone
}

output "gcr_docker_registry" {
  description = "The Docker Registry URL for the GCP project"
  value       = "us-central1-docker.pkg.dev/${var.gcp_project}/my-microservices-repo"
}

output "kubernetes_namespace" {
  description = "Namespace for microservices deployment"
  value       = "microservices-app"
}

output "helm_release_name" {
  description = "Name of the Helm release"
  value       = "microservices"
}

output "service_account_email" {
  description = "Email of the Google service account used"
  value       = google_service_account.sa.email
}
