package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/Brickchain/go-crypto.v1"
	jose "gopkg.in/square/go-jose.v1"
)

func TestStoredKey_Encrypt(t *testing.T) {
	type test struct {
		name    string
		key     *jose.JsonWebKey
		prepare func(*testing.T, *test)
		kek     []byte
		wantErr bool
	}

	tests := []test{
		{
			name: "New",
			prepare: func(t *testing.T, tt *test) {
				key, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					t.Error(err)
				}

				tt.key = &jose.JsonWebKey{
					Key:       key,
					Algorithm: "RS256",
				}
			},
			kek:     crypto.NewSymmetricKey(jose.A256KW),
			wantErr: false,
		},
		{
			name: "Broken_Key",
			prepare: func(t *testing.T, tt *test) {
				key, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					t.Error(err)
				}

				tt.key = &jose.JsonWebKey{
					Key:       key,
					Algorithm: "RS256",
				}
			},
			kek:     []byte("broken"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(t, &tt)
			}
			s := NewStoredKey("test")
			if err := s.Encrypt(tt.key, tt.kek); (err != nil) != tt.wantErr {
				t.Errorf("TestStoredKey_Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoredKey_Decrypt(t *testing.T) {
	type test struct {
		name    string
		key     *StoredKey
		prepare func(*testing.T, *test)
		kek     []byte
		wantErr bool
	}

	tests := []test{
		{
			name: "New",
			prepare: func(t *testing.T, tt *test) {
				tt.key = NewStoredKey("test")
				tt.kek = crypto.NewSymmetricKey(jose.A256KW)

				key, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					t.Error(err)
				}

				jwk := &jose.JsonWebKey{
					Key:       key,
					Algorithm: "RS256",
				}
				if err := tt.key.Encrypt(jwk, tt.kek); err != nil {
					t.Error(err)
				}
			},
			wantErr: false,
		},
		{
			name: "Wrong_Key",
			prepare: func(t *testing.T, tt *test) {
				tt.key = NewStoredKey("test")
				key, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					t.Error(err)
				}

				jwk := &jose.JsonWebKey{
					Key:       key,
					Algorithm: "RS256",
				}
				if err := tt.key.Encrypt(jwk, crypto.NewSymmetricKey(jose.A256KW)); err != nil {
					t.Error(err)
				}
			},
			kek:     crypto.NewSymmetricKey(jose.A256KW),
			wantErr: true,
		},
		{
			name: "Broken_Key",
			prepare: func(t *testing.T, tt *test) {
				tt.key = NewStoredKey("test")
				key, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					t.Error(err)
				}

				jwk := &jose.JsonWebKey{
					Key:       key,
					Algorithm: "RS256",
				}
				if err := tt.key.Encrypt(jwk, crypto.NewSymmetricKey(jose.A256KW)); err != nil {
					t.Error(err)
				}
			},
			kek:     []byte("broken"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(t, &tt)
			}
			if _, err := tt.key.Decrypt(tt.kek); (err != nil) != tt.wantErr {
				t.Errorf("TestStoredKey_Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
