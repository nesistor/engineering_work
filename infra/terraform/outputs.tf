output "gke_cluster_name" {
  description = "The name of the GKE cluster"
  value       = google_container_cluster.gke_cluster.name
}

output "gke_cluster_endpoint" {
  description = "The endpoint of the GKE cluster"
  value       = google_container_cluster.gke_cluster.endpoint
}

output "gke_cluster_zone" {
  description = "The zone of the GKE cluster"
  value       = google_container_cluster.gke_cluster.location
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
