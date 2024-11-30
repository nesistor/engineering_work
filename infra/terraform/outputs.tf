output "gke_endpoint" {
  value = google_container_cluster.gke.endpoint
}

output "cluster_name" {
  value = google_container_cluster.gke.name
}
