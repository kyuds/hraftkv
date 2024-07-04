package main

import "fmt"

type KVStore interface {
	// get op
	Get(key string) (string, error)

	// put op (or replace)
	Put(key, value string) error

	// delete op
	Delete(key string) error

	// join existing cluster
	Join(nodeID, address string) error
}

type RaftKV struct {
}

// for testing purposes

type DummyKV struct{}

func (d *DummyKV) Get(key string) (string, error) {
	fmt.Printf("dummy get: %s!\n", key)
	return "dummy get!", nil
}

func (d *DummyKV) Put(key, value string) error {
	fmt.Printf("dummy put: %s, %s!\n", key, value)
	return nil
}

func (d *DummyKV) Delete(key string) error {
	fmt.Printf("dummy delete: %s!\n", key)
	return nil
}

func (d *DummyKV) Join(nodeID, address string) error {
	fmt.Printf("dummy put: %s, %s!\n", nodeID, address)
	return nil
}
