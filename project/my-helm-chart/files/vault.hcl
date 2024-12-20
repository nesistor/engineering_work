# Listener configuration for Vault API and cluster communication
listener "tcp" {
  address       = "0.0.0.0:8200"  # Vault API listener (accessed externally)
  cluster_address = "0.0.0.0:8201" # Vault cluster communication listener (for HA)
  tls_disable   = true  # Disable TLS for simplicity (set to 'false' in production)
}

# Storage backend configuration (Using Consul in this case)
storage "consul" {
  address = "{{ .Values.vault.ha.backend.consul.address }}"
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
  # Kubernetes service URL
  kubelet_address = "https://kubernetes.default.svc"
  
  # CA certificate used to verify the Kubernetes API server
  kubelet_ca_cert = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
  
  # JWT token used to authenticate the service account to Vault
  token_reviewer_jwt = "/var/run/secrets/kubernetes.io/serviceaccount/token"
  
  # Kubernetes API server URL
  kubernetes_host = "https://kubernetes.default.svc"
}

# Enable any additional secrets engines here (e.g., kv, transit)
# Example: kv for key-value store
# secrets "kv" {
#   path = "secret/"
# }

