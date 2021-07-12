package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var (
	setPath = regexp.MustCompile(`^\/set\/(\w+)]*$`)
	getPath  = regexp.MustCompile(`^\/retrieve\/(\w+)]*$`)
	existsPath = regexp.MustCompile(`^\/exists\/(\w+)]*$`)
	removePath = regexp.MustCompile(`^\/remove\/(\w+)]*$`)
)

type Store interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte) error
	Remove(key string) error
}

type storeHandler struct {
	store Store
}

func (h *storeHandler) set(w http.ResponseWriter, r *http.Request) {
	params := setPath.FindStringSubmatch(r.URL.Path)
	body, err :=  ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.store.Set(params[1], body)
	if err != nil {
		switch err {
		case StoreFull:
			w.WriteHeader(http.StatusInsufficientStorage)
		case KeyExists:
			w.WriteHeader(http.StatusConflict)
		case DataExceedsLimit:
			w.WriteHeader(http.StatusRequestEntityTooLarge)
		default:
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

func (h *storeHandler) retrieve(w http.ResponseWriter, r *http.Request) {
	params := getPath.FindStringSubmatch(r.URL.Path)
	value, exists := h.store.Get(params[1])
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(value)
}

func (h *storeHandler) exists(w http.ResponseWriter, r *http.Request) {
	params := existsPath.FindStringSubmatch(r.URL.Path)
	response := "false"
	if _, exists := h.store.Get(params[1]); exists {
		response = "true"
	}
	w.Write([]byte(response))
}

func (h *storeHandler) remove(w http.ResponseWriter, r *http.Request) {
	params := removePath.FindStringSubmatch(r.URL.Path)
	err := h.store.Remove(params[1])
	if err != nil && err == KeyNotExist {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *storeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && setPath.MatchString(r.URL.Path):
		h.set(w, r)
		return
	case r.Method == http.MethodGet && getPath.MatchString(r.URL.Path):
		h.retrieve(w, r)
		return
	case r.Method == http.MethodGet && existsPath.MatchString(r.URL.Path):
		h.exists(w, r)
		return
	case r.Method == http.MethodDelete && removePath.MatchString(r.URL.Path):
		h.remove(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func main() {
	log.Println("starting application")
	mux := http.NewServeMux()
	mux.Handle("/", &storeHandler{store: NewInMemoryStore()})
	http.ListenAndServe(":80", mux)
}