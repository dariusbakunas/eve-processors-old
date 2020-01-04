package esi

import "testing"

func TestCrypt(t *testing.T) {
	crypt := Crypt{key: "owifu092uroiusodifjoisfoisdufiosu"}
	encrypted, err := crypt.Encrypt("This is plain text")

	if err != nil {
		t.Errorf("Encrypt error: %q", err)
	}

	t.Logf("Encrypted: %s", encrypted)

	decrypted, err := crypt.Decrypt(encrypted)

	if err != nil {
		t.Errorf("Decrypt error: %q", err)
	}

	if decrypted != "This is plain text" {
		t.Errorf("Decrypt() = %q, want %s", decrypted, "This is plain text")
	}
}