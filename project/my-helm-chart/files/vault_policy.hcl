path "sys/mounts" {
  capabilities = ["create", "read", "update", "delete"]
}

path "secret/*" {
  capabilities = ["read"]
}
