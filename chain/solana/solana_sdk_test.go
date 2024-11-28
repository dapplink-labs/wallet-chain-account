package solana

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	// Test basic functionality
	// Generated Private Key (Hex): 88d592a872efde04ec0623db61f7d0eeccfc74b4f44b0c6dd0f4755829be413a7e376c64c64e88054b7a2d25dc716f45551d2f796ddc9e7be405e49c522b887c
	// Generated Public Key (Hex): 7e376c64c64e88054b7a2d25dc716f45551d2f796ddc9e7be405e49c522b887c
	// Generated Solana Address (Base58): 9VhPRjzizPY95TyBrve7heeJTZnofgkQYJpLxRSZGZ3H
	t.Run("Basic Generation", func(t *testing.T) {
		privateKey, publicKey, address, err := generateKeyPair()

		// Check for errors
		if err != nil {
			t.Fatalf("Failed to generate key pair: %v", err)
		}

		t.Logf("Generated Private Key (Hex): %s", hex.EncodeToString(*privateKey))
		t.Logf("Generated Public Key (Hex): %s", hex.EncodeToString(*publicKey))
		t.Logf("Generated Solana Address (Base58): %s", address)
		t.Logf("Private Key Length: %d bytes", len(*privateKey))
		t.Logf("Public Key Length: %d bytes", len(*publicKey))
		t.Logf("Address Length: %d characters", len(address))

		// Check private key length
		if len(*privateKey) != ed25519.PrivateKeySize {
			t.Errorf("Invalid private key length: expected %d, got %d", ed25519.PrivateKeySize, len(*privateKey))
		}

		// Check public key length
		if len(*publicKey) != ed25519.PublicKeySize {
			t.Errorf("Invalid public key length: expected %d, got %d", ed25519.PublicKeySize, len(*publicKey))
		}

		// Check address is not empty
		if len(address) == 0 {
			t.Error("Generated address is empty")
		}

		// Verify relationship between public and private keys
		derivedPubKey := privateKey.Public().(ed25519.PublicKey)
		if !bytes.Equal(derivedPubKey, *publicKey) {
			t.Error("Derived public key does not match generated public key")
		}
	})

	t.Run("Signature Verification", func(t *testing.T) {
		t.Log("Starting signature verification test...")

		privateKey, publicKey, address, err := generateKeyPair()
		if err != nil {
			t.Fatalf("Failed to generate key pair: %v", err)
		}

		t.Logf("Test keypair address: %s", address)

		// Test message
		message := []byte("Hello, Solana!")
		t.Logf("Test message: %s", string(message))

		// Sign
		signature := ed25519.Sign(*privateKey, message)
		t.Logf("Generated signature (Hex): %s", hex.EncodeToString(signature))

		// Verify signature
		if !ed25519.Verify(*publicKey, message, signature) {
			t.Error("Signature verification failed")
		} else {
			t.Log("Signature verification successful")
		}
	})
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, err := generateKeyPair()
		if err != nil {
			b.Fatalf("Failed to generate key pair: %v", err)
		}
	}
}
