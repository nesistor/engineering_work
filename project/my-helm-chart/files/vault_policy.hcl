path "sys/mounts" {
  capabilities = ["create", "read", "update", "delete"]
}

path "secret/*" {
  capabilities = ["create", "read", "update", "delete"]
}

path "jwt_keys/*" {
  capabilities = ["read"]
}
