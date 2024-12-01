#!/bin/bash

# Step 1: Find the pod name for Vault
vault_pod=$(kubectl get pods -n microservices-app -l app=vault -o jsonpath='{.items[0].metadata.name}')

# Check if Vault pod was found
if [ -z "$vault_pod" ]; then
    echo "Vault pod not found in the microservices-app namespace."
    exit 1
fi

echo "Found Vault pod: $vault_pod"

# Step 2: Execute commands in the Vault pod
kubectl exec -it $vault_pod -n microservices-app -- /bin/sh -c "
    # Step 3: Login to Vault using the root token
    vault login root_token

    # Step 4: Enable the secrets engine for JWT keys
    vault secrets enable -path=jwt_keys kv

    # Step 5: Write the JWT policy to Vault
    vault policy write jwt_policy /vault/secrets/jwt-policy/jwt_policy.hcl

    # Step 6: Store the public and private keys in Vault
    vault kv put jwt_keys/public_key public_key=@/vault/secrets/jwt-keys/public_key.pem
    vault kv put jwt_keys/private_key private_key=@/vault/secrets/jwt-keys/private_key.pem

    # Step 7: Get the private key from Vault
    vault kv get jwt_keys/private_key
"

echo "Vault setup complete."
