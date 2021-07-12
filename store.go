package main

import (
	"errors"
	"regexp"
)

const maxDataSize = 50
const maxKeyCount = 1000
const validKeyPattern = `^[a-z_]{1}[a-z0-9-_]{1,15}$`

var (
	StoreFull = errors.New("store is full")
	KeyFormatInvalid = errors.New("invalid key")
	KeyExists = errors.New("key already exists")
	DataIsEmpty = errors.New("data must not be empty")
	DataExceedsLimit = errors.New("data length exceeds limit")
	KeyNotExist = errors.New("key does not exist")
)

type InMemoryStore struct {
	data map[string][]byte
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string][]byte),
	}
}

func (s *InMemoryStore) Get(key string) ([]byte, bool) {
	value, exists := s.data[key]
	return value, exists
}

func (s *InMemoryStore) Set(key string, value []byte) error {
	if len(s.data) >= maxKeyCount {
		return StoreFull
	}
	if _, exists := s.data[key]; exists {
		return KeyExists
	}
	if len(value) == 0 {
		return DataIsEmpty
	}
	if len(value) > maxDataSize {
		return DataExceedsLimit
	}
	if keyValid, _ := regexp.MatchString(validKeyPattern, key); !keyValid {
		return KeyFormatInvalid
	}

	s.data[key] = value
	return nil
}

func (s *InMemoryStore) Remove(key string) error {
	_, exists := s.data[key]
	if !exists {
		return KeyNotExist
	}

	delete(s.data, key)
	return nil
}
