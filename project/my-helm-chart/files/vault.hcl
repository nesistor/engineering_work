# Listener Configuration
listener "tcp" {
  address       = "0.0.0.0:8200"
  cluster_address = "0.0.0.0:8201"
  tls_disable   = true
}

# Storage Backend Configuration: Consul
storage "consul" {
  address = "localhost:8500" # Zmień na właściwy adres Consul, np. consul.service.consul:8500
  path    = "vault/"         # Ścieżka przechowywania w Consul
}

# High Availability Configuration
ha_storage "consul" {
  address = "localhost:8500" # Zmień na właściwy adres Consul
  path    = "vault-ha/"      # Ścieżka dla konfiguracji HA
}

# API Configuration
api_addr = "http://0.0.0.0:8200"       # Adres API Vault
cluster_addr = "http://0.0.0.0:8201"   # Adres klastra Vault

# UI Configuration
ui = true
