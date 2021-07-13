package main

import (
	"errors"
	"regexp"
	"sync"
)

const maxDataSize = 50
const maxKeyCount = 1000
var validKeyPattern = regexp.MustCompile(`^[a-z_]{1}[a-z0-9-_]{1,15}$`)

var (
	StoreFull = errors.New("store is full")
	KeyFormatInvalid = errors.New("invalid key")
	KeyExists = errors.New("key already exists")
	DataIsEmpty = errors.New("data must not be empty")
	DataExceedsLimit = errors.New("data length exceeds limit")
	KeyNotExist = errors.New("key does not exist")
)

type InMemoryStore struct {
	sync.RWMutex
	data map[string][]byte
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string][]byte),
	}
}

func (s *InMemoryStore) Get(key string) ([]byte, bool) {
	s.RLock()
	defer s.RUnlock()
	value, exists := s.data[key]
	return value, exists
}

func (s *InMemoryStore) Set(key string, value []byte) error {
	if len(value) == 0 {
		return DataIsEmpty
	}
	if len(value) > maxDataSize {
		return DataExceedsLimit
	}
	if !validKeyPattern.MatchString(key) {
		return KeyFormatInvalid
	}

	s.Lock()
	defer s.Unlock()
	if len(s.data) >= maxKeyCount {
		return StoreFull
	}
	if _, exists := s.data[key]; exists {
		return KeyExists
	}
	s.data[key] = value

	return nil
}

func (s *InMemoryStore) Remove(key string) error {
	s.Lock()
	defer s.Unlock()
	_, exists := s.data[key]
	if !exists {
		return KeyNotExist
	}

	delete(s.data, key)
	return nil
}
