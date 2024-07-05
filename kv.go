package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/raft"
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

type RaftKV struct {
	// settings
	raddr string
	dir   string
	id    raft.ServerID

	// for concurrent requests
	lock sync.Mutex
	kv   map[string]string

	// backing raft
	raft *raft.Raft
}

func NewRaftKV(raddr, dir, id string) *RaftKV {
	return &RaftKV{
		kv:    make(map[string]string),
		raddr: raddr,
		dir:   dir,
		id:    raft.ServerID(id),
	}
}

func (k *RaftKV) Start(join string) error {
	config := raft.DefaultConfig()
	config.LocalID = k.id

	raddr, err := net.ResolveTCPAddr("tcp", k.raddr)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(k.raddr, raddr, 3, 100*time.Millisecond, os.Stdout)
	if err != nil {
		return err
	}
	snapshots, err := raft.NewFileSnapshotStore(k.dir, 3, os.Stdout)
	if err != nil {
		return err
	}
	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()

	k.raft, err = raft.NewRaft(config, k, logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("raft error: %s", err)
	}

	if join == "" {
		// no join indicates cluster is starting from one server
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		k.raft.BootstrapCluster(configuration)
	} else {
		// TODO: join to cluster
	}

	return nil
}

func (k *RaftKV) Stop() {

}

// return empty if is leader.
// Otherwise, return hint for leader.
func (k *RaftKV) IsLeader() string {
	return ""
}

//////////////////////////
//     KVStore Impl     //
//////////////////////////

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
//         fsm          //
//////////////////////////

func (k *RaftKV) Apply(l *raft.Log) interface{} {
	return nil
}

// Snapshot returns a snapshot of the key-value store.
func (k *RaftKV) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

// Restore stores the key-value store to a previous state.
func (k *RaftKV) Restore(rc io.ReadCloser) error {
	return nil
}

type snapshot struct {
	kv map[string]string
}

func (ss *snapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (ss *snapshot) Release() {}

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
