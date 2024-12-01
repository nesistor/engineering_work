path "sys/mounts" {
  capabilities = ["create", "read", "update", "delete"]
}

path "secret/*" {
  capabilities = ["read"]
}

path "jwt_keys/*" {
  capabilities = ["read"]
}
