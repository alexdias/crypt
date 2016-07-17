package server

import (
	"errors"
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

var (
	ErrAlreadyExists = errors.New("id already present in store")
	ErrNotFound      = errors.New("id not found in store")
)

func NewInMemoryStoreService() StoreService {
	return &inMemoryStoreService{
		m: make(map[string][]byte),
	}
}

func (s *inMemoryStoreService) GetData(ctx context.Context, id string, key []byte) ([]byte, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	data, ok := s.m[id]
	if !ok {
		return nil, ErrNotFound
	}
	return data, nil
}

func (s *inMemoryStoreService) PostData(ctx context.Context, data PlainTextData) ([]byte, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[data.ID]; ok {
		return nil, ErrAlreadyExists
	}
	s.m[data.ID] = []byte(data.Data)
	return []byte{104, 101, 108, 108, 111}, nil
}
