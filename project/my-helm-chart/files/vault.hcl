# vault.hcl
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
