package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlesSetKey(t *testing.T) {
	store := &mockStore{}
	handler := &storeHandler{store: store}
	data := []byte("data")
	req := httptest.NewRequest("POST", "/set/abc", bytes.NewReader(data))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Error("expected created response")
	}
	if store.key != "abc" || bytes.Compare(data, store.value) != 0 {
		t.Error("handler failed to set data to store")
	}
}

func TestHandlesGetsKey(t *testing.T) {
	data := []byte("data")
	handler := &storeHandler{store: &mockStore{value: data}}
	req := httptest.NewRequest("GET", "/retrieve/abc", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if bytes.Compare(body, data) != 0 {
		t.Error("expected data in response body")
	}
}

func TestHandlesRemoveKey(t *testing.T) {
	store := &mockStore{}
	handler := &storeHandler{store: store}
	req := httptest.NewRequest("DELETE", "/remove/abc", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNoContent {
		t.Error("expected deleted response")
	}
	if store.deleted != "abc" {
		t.Error("expected key to be deleted from store")
	}
}

func TestHandlesKeyExists(t *testing.T) {
	handler := &storeHandler{store: &mockStore{}}
	req := httptest.NewRequest("GET", "/exists/abc", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if bytes.Compare(body, []byte("true")) != 0 {
		t.Error("expected true in response")
	}
}

type mockStore struct {
	key string
	value []byte
	deleted string
}

func (s *mockStore) Get(key string) ([]byte, bool) {
	return s.value, true
}
func (s *mockStore) Set(key string, value []byte) error {
	s.key = key
	s.value = value
	return nil
}
func (s *mockStore) Remove(key string) error {
	s.deleted = key
	return nil
}
