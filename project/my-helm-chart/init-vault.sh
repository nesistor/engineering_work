#!/bin/bash

# Ustawienie adresu Vault
export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=$(kubectl get secret vault-root-token-secret -o jsonpath="{.data.token}" | base64 -d)

# Dodanie kluczy JWT
vault kv put jwt_keys/private key=@/home/karol/code/microservices-new/done/project/my-helm-chart/files/private_key.pem
vault kv put jwt_keys/public key=@/home/karol/code/microservices-new/done/project/my-helm-chart/files/public_key.pem

# Tworzenie polityki JWT
vault policy write jwt_policy - <<EOF
path "jwt_keys/*" {
  capabilities = ["read", "list"]
}
EOF

# Tworzenie roli Kubernetes
vault write auth/kubernetes/role/my-role \
  bound_service_account_names=vault \
  bound_service_account_namespaces=default \
  policies=jwt_policy \
  ttl=24h
