resource "google_kms_key_ring" "vault_keyring" {
  name     = var.kms_keyring
  location = "global"
}

resource "google_kms_crypto_key" "vault_init_key" {
  name     = var.kms_key
  key_ring = google_kms_key_ring.vault_keyring.id
  purpose  = "ENCRYPT_DECRYPT"
}
