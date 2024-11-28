storage "file" {
  path = "/vault/data"  # Lokalizacja przechowywania danych w kontenerze
}

# HTTP API
api_addr = "http://0.0.0.0:8200"
cluster_addr = "https://0.0.0.0:8201"  # Wewnętrzny adres klastra Vault (dla HA)

# Adresy i certyfikaty TLS
listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = 1  # Zablokowanie TLS, należy włączyć TLS w produkcji
}

# Autentykacja
disable_mlock = true  # To jest opcjonalne, zależy od Twojego środowiska

# Konfiguracja Kubernetes Auth Method
auth "kubernetes" {
  # Adres do API Kubernetes
  kubernetes_host = "https://kubernetes.default.svc"

  # Tokeny autentykacji z Kubernetes (domyślnie Vault użyje service account token w podach)
  kubernetes_ca_cert = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"  # Certyfikat CA Kubernetes
  kubernetes_token_reviewer_jwt = "/var/run/secrets/kubernetes.io/serviceaccount/token"  # Token serwisu
}
