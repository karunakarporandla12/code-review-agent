package security

import "testing"

func TestEncryptDecryptRoundTrip(t *testing.T) {
	plain := "const app = { secret: \"super-secret\" };"
	cipherText, err := EncryptString("test-encryption-key", plain)
	if err != nil {
		t.Fatalf("encrypt returned error: %v", err)
	}
	if cipherText == plain {
		t.Fatal("ciphertext should not match plaintext")
	}
	out, err := DecryptString("test-encryption-key", cipherText)
	if err != nil {
		t.Fatalf("decrypt returned error: %v", err)
	}
	if out != plain {
		t.Fatalf("expected %q, got %q", plain, out)
	}
}

func TestDecryptLegacyPlainTextReturnsOriginalValue(t *testing.T) {
	value := "plain-text-value"
	out, err := DecryptString("test-encryption-key", value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != value {
		t.Fatalf("expected %q, got %q", value, out)
	}
}
