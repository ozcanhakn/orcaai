package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "io"
)

func getKeyBytes(secret string) ([]byte, error) {
    // expect hex-encoded 32 bytes (64 chars)
    b, err := hex.DecodeString(secret)
    if err != nil {
        return nil, err
    }
    if len(b) != 32 {
        return nil, errors.New("secret key must be 32 bytes (AES-256)")
    }
    return b, nil
}

func EncryptAESGCM(plaintext []byte, secretHex string) (string, error) {
    key, err := getKeyBytes(secretHex)
    if err != nil { return "", err }
    block, err := aes.NewCipher(key)
    if err != nil { return "", err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return "", err }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return "", err }
    sealed := gcm.Seal(nonce, nonce, plaintext, nil)
    return hex.EncodeToString(sealed), nil
}

func DecryptAESGCM(cipherHex string, secretHex string) ([]byte, error) {
    key, err := getKeyBytes(secretHex)
    if err != nil { return nil, err }
    data, err := hex.DecodeString(cipherHex)
    if err != nil { return nil, err }
    block, err := aes.NewCipher(key)
    if err != nil { return nil, err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return nil, err }
    if len(data) < gcm.NonceSize() { return nil, errors.New("ciphertext too short") }
    nonce := data[:gcm.NonceSize()]
    ct := data[gcm.NonceSize():]
    return gcm.Open(nil, nonce, ct, nil)
}


