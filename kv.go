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
	Put(key, value string) error

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
	k.raft.Shutdown().Error()
}

// return empty if is leader.
// Otherwise, return hint for leader.
func (k *RaftKV) IsLeader() (bool, string) {
	if k.raft.State() == raft.Leader {
		return true, ""
	}
	addr, _ := k.raft.LeaderWithID()
	return false, string(addr)
}

//////////////////////////
//     KVStore Impl     //
//////////////////////////

type cmd struct {
	Operation string
	Key       string
	Value     string
}

func (k *RaftKV) Get(key string) (string, error) {
	isLeader, hint := k.IsLeader()
	if !isLeader {
		return "", fmt.Errorf("not leader. Leader is at %s", hint)
	}
	k.lock.Lock()
	defer k.lock.Unlock()
	return k.kv[key], nil
}

func (k *RaftKV) Put(key, value string) error {
	isLeader, hint := k.IsLeader()
	if !isLeader {
		return fmt.Errorf("not leader. Leader is at %s", hint)
	}

	c := &cmd{
		Operation: "put",
		Key:       key,
		Value:     value,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f := k.raft.Apply(b, time.Second)
	return f.Error()
}

func (k *RaftKV) Delete(key string) error {
	isLeader, hint := k.IsLeader()
	if !isLeader {
		return fmt.Errorf("not leader. Leader is at %s", hint)
	}

	k.lock.Lock()
	defer k.lock.Unlock()
	c := &cmd{
		Operation: "del",
		Key:       key,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f := k.raft.Apply(b, time.Second)
	return f.Error()
}

func (k *RaftKV) Join(nodeID, address string) error {
	return nil
}

//////////////////////////
//         fsm          //
//////////////////////////

func (k *RaftKV) Apply(l *raft.Log) interface{} {
	var c cmd
	err := json.Unmarshal(l.Data, &c)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}
	k.lock.Lock()
	defer k.lock.Unlock()

	if c.Operation == "put" {
		k.kv[c.Key] = c.Value
	} else if c.Operation == "del" {
		delete(k.kv, c.Key)
	} else {
		panic(fmt.Sprintf("unknown operation: %s", c.Operation))
	}
	return nil
}

// Snapshot returns a snapshot of the key-value store.
func (k *RaftKV) Snapshot() (raft.FSMSnapshot, error) {
	k.lock.Lock()
	defer k.lock.Unlock()
	snp := make(map[string]string)
	for k, v := range k.kv {
		snp[k] = v
	}
	return &snapshot{kv: snp}, nil
}

// Restore stores the key-value store to a previous state.
// Not called concurrently, per godoc
func (k *RaftKV) Restore(rc io.ReadCloser) error {
	tmp := make(map[string]string)
	err := json.NewDecoder(rc).Decode(&tmp)
	if err != nil {
		return err
	}
	k.kv = tmp
	return nil
}

type snapshot struct {
	kv map[string]string
}

func (ss *snapshot) Persist(sink raft.SnapshotSink) error {
	b, err := json.Marshal(ss.kv)
	if err != nil {
		sink.Cancel()
		return err
	}
	_, err = sink.Write(b)
	if err != nil {
		sink.Cancel()
		return err
	}
	err = sink.Close()
	if err != nil {
		sink.Cancel()
	}
	return err
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

func (d *DummyKV) Put(key, value string) error {
	fmt.Printf("dummy put: %s, %s.\n", key, value)
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
