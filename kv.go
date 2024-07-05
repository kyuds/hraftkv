package main

import (
	"encoding/json"
	"fmt"
)

type KVStore interface {
	// get op
	Get(key string) (string, error)

	// put op (or replace)
	Put(pair map[string]string) error

	// delete op
	Delete(key string) error

	// join existing cluster
	Join(nodeID, address string) error
}

type RaftKV struct{}

func NewRaftKV(raddr, join, id, dir string) *RaftKV {
	return &RaftKV{}
}

func (k *RaftKV) Start() error {
	return nil
}

func (k *RaftKV) Stop() {

}

func (k *RaftKV) Get(key string) (string, error) {
	return "", nil
}

func (k *RaftKV) Put(pair map[string]string) error {
	return nil
}

func (k *RaftKV) Delete(key string) error {
	return nil
}

func (k *RaftKV) Join(nodeID, address string) error {
	return nil
}

//////////////////////////
// for testing purposes //
//////////////////////////

type DummyKV struct{}

func (d *DummyKV) Get(key string) (string, error) {
	fmt.Printf("dummy get: %s.\n", key)
	return "dummy get!", nil
}

func (d *DummyKV) Put(pair map[string]string) error {
	bs, _ := json.Marshal(pair)
	fmt.Printf("dummy put: %s.\n", string(bs))
	return nil
}

func (d *DummyKV) Delete(key string) error {
	fmt.Printf("dummy delete: %s.\n", key)
	return nil
}

func (d *DummyKV) Join(nodeID, address string) error {
	fmt.Printf("dummy put: %s, %s.\n", nodeID, address)
	return nil
}
