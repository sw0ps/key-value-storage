package main

import (
	"errors"
	"testing"
)

func TestPut(t *testing.T) {
	const (
		key   = "create-key"
		value = "create-value"
	)
	var (
		err      error
		contains bool
		val      interface{}
	)

	defer delete(store, key)

	if _, contains = store[key]; contains {
		t.Error("key/value already exists")
	}

	err = Put(key, value)
	if err != nil {
		t.Error(err)
	}

	val, contains = store[key]
	if !contains {
		t.Error("create failed")
	}

	if val != value {
		t.Error("value is incorrect")
	}
}

func TestGet(t *testing.T) {
	const (
		key   = "read-key"
		value = "read-value"
	)
	var (
		err error
		val interface{}
	)

	defer delete(store, key)

	val, err = Get(key)
	if err == nil {
		t.Error("expected error")
	}
	if !errors.Is(err, ErrorNoSuchKey) {
		t.Error("unexpected error:", err)
	}

	store[key] = value

	val, err = Get(key)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if val != value {
		t.Error("value is incorrect")
	}
}

func TestDelete(t *testing.T) {
	const (
		key   = "delete-key"
		value = "delete-value"
	)
	var (
		contains bool
		err      error
	)

	defer delete(store, key)

	store[key] = value

	err = Delete(key)
	if err != nil {
		t.Error(err)
	}

	_, contains = store[key]
	if contains {
		t.Error("delete failed")
	}
}
