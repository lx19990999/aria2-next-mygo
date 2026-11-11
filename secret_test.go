package main

import "testing"

func TestEncryptDecryptSecret(t *testing.T) {
	plain := "DiRoKMa7b1wT"
	enc, err := encryptSecret(plain)
	if err != nil {
		t.Fatal(err)
	}
	if enc == "" || enc == plain {
		t.Fatalf("expected opaque ciphertext, got %q", enc)
	}
	got, err := decryptSecret(enc)
	if err != nil {
		t.Fatal(err)
	}
	if got != plain {
		t.Fatalf("got %q want %q", got, plain)
	}
}
