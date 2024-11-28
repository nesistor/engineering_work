# Polityka dla aplikacji do odczytu tajemnic z Vault
path "jwt_keys/*" {
  capabilities = ["read"]
}
