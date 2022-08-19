package main

import (
	"errors"
	"sync"
)

// store - general variable

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

// errors variables

var (
	// ErrorNoSuchKey using for no such key
	ErrorNoSuchKey = errors.New("no such key")
)

// functions

// Put key to storage
func Put(key, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()

	return nil
}

// Get key from storage
func Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

// Delete key from storage
func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	return nil
}
