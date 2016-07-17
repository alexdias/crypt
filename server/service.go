package server

import (
	"sync"

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

type inMemoryStoreService struct {
	mtx sync.RWMutex
	m   map[string][]byte
}

func NewInMemoryStoreService() StoreService {
	return &inMemoryStoreService{
		m: make(map[string][]byte),
	}
}

func (s *inMemoryStoreService) GetData(ctx context.Context, id string, key []byte) ([]byte, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return key, nil
}

func (s *inMemoryStoreService) PostData(ctx context.Context, data PlainTextData) ([]byte, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return []byte{104, 101, 108, 108, 111}, nil
}
