package main

import "errors"

// store - general variable

var store = make(map[string]string)

// errors variables

var (
	ErrorNoSuchKey = errors.New("no such key")
)

// functions

func Put(key, value string) error {
	store[key] = value

	return nil
}

func Get(key string) (string, error) {
	value, ok := store[key]
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func Delete(key string) error {
	delete(store, key)

	return nil
}
