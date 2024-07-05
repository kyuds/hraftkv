# hraftkv
Simple KV Store on HashiCorp Raft

### Setup
Run the following:
```
./build.sh
```

### Launching a Cluster
```
# from the command line
./hraftkv -mem -aaddr=localhost:10000 -raddr=localhost:20000 -id=node0
```
or generate commands for a `n` node cluster with attached script
```
./cmd.sh <n>
```

### KV Store Operations
```
curl -X GET "localhost:10000/kv?key=test"
curl -X POST "localhost:10000/kv" -d '{"test": "success"}'
curl -X DELETE "localhost:10000/kv?key=test"
```
