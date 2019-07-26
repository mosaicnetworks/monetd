package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pborman/uuid"
)

const (

	// Version 6219 corresponds to Monet. We just added the public key as one of the
	// fields in the JSON object because it makes our lives easier when working with
	// Babble.
	monetVersion = 6219

	ethereumVersion = 3
)

type encryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	Id      string              `json:"id"`
	Version int                 `json:"version"`
}

type encryptedKeyJSONMonet struct {
	Address   string              `json:"address"`
	PublicKey string              `json:"pub"`
	Crypto    keystore.CryptoJSON `json:"crypto"`
	Id        string              `json:"id"`
	Version   int                 `json:"version"`
}

// EncryptKey encrypts a key using the specified scrypt parameters into a json
// blob that can be decrypted later on.
func EncryptKey(key *keystore.Key, auth string, scryptN, scryptP int) ([]byte, error) {
	keyBytes := math.PaddedBigBytes(key.PrivateKey.D, 32)
	cryptoStruct, err := keystore.EncryptDataV3(keyBytes, []byte(auth), scryptN, scryptP)
	if err != nil {
		return nil, err
	}
	encryptedKeyJSONMonet := encryptedKeyJSONMonet{
		hex.EncodeToString(key.Address[:]),
		hex.EncodeToString(crypto.FromECDSAPub(&key.PrivateKey.PublicKey)),
		cryptoStruct,
		key.Id.String(),
		monetVersion,
	}
	return json.Marshal(encryptedKeyJSONMonet)
}

// DecryptKey decrypts a key from a json blob, returning the private key itself.
func DecryptKey(keyjson []byte, auth string) (*keystore.Key, error) {
	// Parse the json into a simple map to fetch the key version
	m := make(map[string]interface{})
	if err := json.Unmarshal(keyjson, &m); err != nil {
		return nil, err
	}

	// Depending on the version try to parse one way or another
	var (
		keyBytes []byte
		err      error
	)

	k := new(encryptedKeyJSONV3)
	if err := json.Unmarshal(keyjson, k); err != nil {
		return nil, err
	}
	keyBytes, _, err = decryptKeyV3(k, auth)

	// Handle any decryption errors and return the key
	if err != nil {
		return nil, err
	}
	privateKey := crypto.ToECDSAUnsafe(keyBytes)

	key := WrapKey(privateKey)

	return key, nil
}

func decryptKeyV3(keyProtected *encryptedKeyJSONV3, auth string) (keyBytes []byte, keyId []byte, err error) {
	if keyProtected.Version != monetVersion && keyProtected.Version != ethereumVersion {
		return nil, nil, fmt.Errorf("Version not supported: %v", keyProtected.Version)
	}
	keyId = uuid.Parse(keyProtected.Id)
	plainText, err := keystore.DecryptDataV3(keyProtected.Crypto, auth)
	if err != nil {
		return nil, nil, err
	}
	return plainText, keyId, err
}

// WrapKey Create the keyfile object with a random UUID. It would be preferable
// to create the key manually, rather then calling this function, but we cannot
// use pborman/uuid directly because it is vendored in go-ethereum. That package
// defines the type of keystore.Key.Id.
func WrapKey(privateKey *ecdsa.PrivateKey) *keystore.Key {
	key := keystore.NewKeyForDirectICAP(rand.Reader)
	key.Address = crypto.PubkeyToAddress(privateKey.PublicKey)
	key.PrivateKey = privateKey
	return key
}
