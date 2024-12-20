# Listener configuration for Vault API and cluster communication
listener "tcp" {
  address = "0.0.0.0:8200"        # Vault API listener (accessed externally)
  cluster_address = "0.0.0.0:8201" # Vault cluster communication listener (for HA)
  tls_disable = true                # Disable TLS for simplicity (set to 'false' in production)
}

# Storage backend configuration (Using Consul in this case)
storage "consul" {
  address = "microservices-consul-0.microservices-app.svc.cluster.local:8500"
  path    = "{{ .Values.vault.ha.backend.consul.path }}"
  scheme  = "{{ .Values.vault.ha.backend.consul.scheme }}"
}

# API and Cluster addresses for Vault
api_addr = "http://{{ .Release.Name }}-vault-service:8200"
cluster_addr = "http://{{ .Release.Name }}-vault:8201"

# Enable Vault UI
ui = true

# Kubernetes Auth configuration
auth "kubernetes" {
  kubelet_address    = "https://kubernetes.default.svc"
  kubelet_ca_cert    = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
  token_reviewer_jwt = "/var/run/secrets/kubernetes.io/serviceaccount/token"
  kubernetes_host    = "https://kubernetes.default.svc"
}

# Enable transit backend for key management
secrets "transit" {
  # Here, you can define any additional options if needed
  # Example: path = "transit/"
}

# Enable any additional secrets engines here (e.g., kv, transit)
# Example: kv for key-value store
# secrets "kv" {
#   path = "secret/"
# }
