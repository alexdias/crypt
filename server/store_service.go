package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/net/context"
)

type StoreService interface {
	PostData(ctx context.Context, data PlainTextData) ([]byte, error)
	GetData(ctx context.Context, id string, key []byte) ([]byte, error)
}

type PlainTextData struct {
	ID   string `json:"id"`
	Data string `json:"plaintext"`
}

type storeService struct {
	storage Storage
}

var (
	ErrAlreadyExists  = errors.New("id already present in store")
	ErrNotFound       = errors.New("id not found in store")
	ErrGeneratingKey  = errors.New("error generating key")
	ErrEncryptingData = errors.New("error encrypting the data")
	ErrDecryptingData = errors.New("error decrypting the data")
)

func NewStoreService(storage Storage) StoreService {
	return &storeService{
		storage: storage,
	}
}

func (s *storeService) GetData(ctx context.Context, id string, key []byte) ([]byte, error) {
	ciphertext, storage_err := s.storage.get(id)
	if storage_err != nil {
		return nil, ErrNotFound
	}
	plaintext, err := decrypt(key, ciphertext)
	if err != nil {
		return nil, ErrDecryptingData
	}
	return plaintext, nil
}

func (s *storeService) PostData(ctx context.Context, data PlainTextData) ([]byte, error) {
	key := make([]byte, 32) // 32 bytes for AES-256
	_, err := rand.Read(key)
	if err != nil {
		return nil, ErrGeneratingKey
	}
	ciphertext, encrypt_err := encrypt(key, []byte(data.Data))
	if encrypt_err != nil {
		return nil, ErrEncryptingData
	}
	storage_err := s.storage.put(data.ID, ciphertext)
	if storage_err != nil {
		return nil, ErrAlreadyExists
	}
	return key, nil
}

func encrypt(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	base_string := base64.StdEncoding.EncodeToString(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(base_string))
	iv := ciphertext[:aes.BlockSize]
	// generate a random IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	byte_string := []byte(base_string)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], byte_string)
	return ciphertext, nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrDecryptingData
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, ErrDecryptingData
	}
	iv := ciphertext[:aes.BlockSize]
	encrypted_text := make([]byte, len(ciphertext)-aes.BlockSize)
	copy(encrypted_text, ciphertext[aes.BlockSize:])
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(encrypted_text, encrypted_text)
	string_text := string(encrypted_text)
	data, err := base64.StdEncoding.DecodeString(string_text)
	if err != nil {
		return nil, ErrDecryptingData
	}
	return data, nil
}
