package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
)

// This is not production grade
//
// assumptions
// - underlying Raft instance will not crash
//
// some settings
// - storage is always in memory

func main() {
	// parse command line arguments
	var aaddr string
	var raddr string
	var join string
	var id string
	var dir string

	flag.StringVar(&aaddr, "aaddr", "", "API address")
	flag.StringVar(&raddr, "raddr", "", "Raft address")
	flag.StringVar(&join, "join", "", "Set join address, if any")
	flag.StringVar(&id, "id", "", "Node ID. Default to Raft address")
	flag.StringVar(&dir, "dir", "", "Raft storage directory. Default to nodeID")

	flag.Parse()

	if len(aaddr) == 0 || len(raddr) == 0 {
		fmt.Fprintf(os.Stderr, "Some required parameters not specified\n")
		os.Exit(1)
	}
	if len(id) == 0 {
		id = raddr
	}
	if len(dir) == 0 {
		dir = id
	}

	// kvstore settings
	fmt.Printf("Starting hraftkv node <%s> on %s\n", id, aaddr)
	fmt.Printf("Backing raft on %s in directory <%s>\n", raddr, dir)

	// start kvstore
	kv := NewRaftKV(raddr, dir, id)
	if err := kv.Start(join); err != nil {
		fmt.Fprintf(os.Stderr, "Raft failed to start: %s\n", err.Error())
	}
	server := NewServer(aaddr, kv)
	server.Start()
	fmt.Println("Started successfully")

	// wait for HTTP server to terminate.
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	server.Close()
	kv.Stop()
	fmt.Println("terminated")
}
