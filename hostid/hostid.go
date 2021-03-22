package hostid

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
)

const KEY_BIT_SIZE = 2048

type HostIdentity struct {
	HostId string
	PrivateKey []byte
	PublicKey []byte
	Version int
}

type PublicHostIdentity struct {
	HostId string
	PublicKey []byte
	Version int
}

func GenerateIdentity() (*HostIdentity, error) {
	privateKey, publicKey, err := generateKeys()
	if err != nil {
		return nil, err
	}

	return &HostIdentity{
		HostId: hostIdFromPublicKey(publicKey),
		PrivateKey: privateKey,
		PublicKey: publicKey,
		Version: 1,
	}, nil
}

func (h *HostIdentity) SignNonce(nonce string) ([]byte, error) {
	hashEngine := sha256.New()
	hashEngine.Write([]byte(nonce))
	hash := hashEngine.Sum(nil)

	privateKey, err := loadKeyFromBytes(h.PrivateKey)
	if err != nil {
		return nil, err
	}
	res, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *HostIdentity) ToPublicHostIdentity() *PublicHostIdentity {
	return &PublicHostIdentity{
		HostId: h.HostId,
		PublicKey: h.PublicKey,
		Version: h.Version,
	}
}

func hostIdFromPublicKey(public []byte) string {
	hashEngine := md5.New()
	hashEngine.Write(public)
	hash := hashEngine.Sum(nil)

	return string(hash)
}

func generateKeys() (private []byte, public []byte, err error) {
	reader := rand.Reader

	key, err := rsa.GenerateKey(reader, KEY_BIT_SIZE)
	if err != nil {
		return
	}

	private, err = keyToPrivateKeyBytes(key)
	if err != nil {
		return
	}
	public, err = keyToPublicKeyString(key.PublicKey)
	if err != nil {
		return
	}

	return
}

func keyToPublicKeyString(pubkey rsa.PublicKey) ([]byte, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(&pubkey)
	if err != nil {
		return nil, err
	}

	result := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubkey_bytes,
	})

	return result, nil
}

func keyToPrivateKeyBytes(key *rsa.PrivateKey) ([]byte, error) {
	privateGobBuf := new(bytes.Buffer)
	encoder := gob.NewEncoder(privateGobBuf)
	err := encoder.Encode(key)
	if err != nil {
		return nil, err
	}

	return privateGobBuf.Bytes(), nil
}

func loadKeyFromBytes(input []byte) (*rsa.PrivateKey, error) {
	var result rsa.PrivateKey
	privateGobBuf := bytes.NewReader(input)
	decoder := gob.NewDecoder(privateGobBuf)
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
