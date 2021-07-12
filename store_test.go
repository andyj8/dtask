package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

func TestRejectsInvalidKeys(t *testing.T) {
	var invalids = []string {
		"",
		"thiskeyiswaytoolongtobevalid",
		"inval!dch@rs",
		"1numberatstart",
	}

	store := NewInMemoryStore()
	for _, key := range invalids {
		if store.Set(key, []byte("data")) != KeyFormatInvalid {
			t.Errorf("%s key should not be accepted", key)
		}
	}
}

func TestRejectsDuplicateKey(t *testing.T) {
	store := NewInMemoryStore()
	store.Set("test", []byte("data"))
	if store.Set("test", []byte("data")) != KeyExists {
		t.Error("expected to received key exists error")
	}
}

func TestRejectsBodyOverLimit(t *testing.T) {
	store := NewInMemoryStore()
	data := make([]byte, maxDataSize + 1)
	rand.Read(data)
	if store.Set("key", data) != DataExceedsLimit {
		t.Error("expected to received data exceeds limit error")
	}
}

func TestRejectsAllWhenAtCapacity(t *testing.T) {
	store := NewInMemoryStore()
	for i := 1; i <= maxKeyCount; i++ {
		store.Set(fmt.Sprintf("_%d", i), []byte("data"))
	}

	err := store.Set("key", []byte("data"))
	if err != StoreFull {
		t.Error("expected to receive store full error")
	}
}

func TestSetsKey(t *testing.T) {
	store := NewInMemoryStore()
	if store.Set("key", []byte("data")) != nil {
		t.Error("unexpected set error")
	}
}

func TestGetsKey(t *testing.T) {
	store := NewInMemoryStore()
	data := []byte("data")
	store.Set("key", data)
	retrieved, found := store.Get("key")
	if bytes.Compare(data, retrieved) != 0 || found != true {
		t.Error("unexpected error getting key")
	}
}

func TestRemovesKey(t *testing.T) {
	store := NewInMemoryStore()
	store.Set("key", []byte("data"))
	store.Remove("key")
	_, found := store.Get("key")
	if found != false {
		t.Error("expected key to have been removed")
	}
}

