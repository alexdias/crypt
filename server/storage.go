package server

import (
	"errors"
	"sync"
)

type Storage interface {
	put(string, []byte) error
	get(string) ([]byte, error)
}

type inMemoryStorage struct {
	mtx sync.RWMutex
	m   map[string][]byte
}

var (
	ErrAlreadyInStorage = errors.New("already in storage")
	ErrNotInStorage     = errors.New("not in storage")
)

func NewInMemoryStorage() Storage {
	return &inMemoryStorage{
		m: make(map[string][]byte),
	}
}

func (s *inMemoryStorage) put(id string, value []byte) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; ok {
		return ErrAlreadyInStorage
	}
	s.m[id] = value
	return nil
}

func (s *inMemoryStorage) get(id string) ([]byte, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	value, ok := s.m[id]
	if !ok {
		return nil, ErrNotInStorage
	}
	return value, nil
}
