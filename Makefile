SOURCE = $(wildcard *.go)
PROTO = pb/net.proto
PROTO_GEN = $(wildcard pb/*.pb.go)
OUTPUT = hraftkv

build: proto $(SOURCE)
	go build -o $(OUTPUT)

proto: $(PROTO)
	protoc --go_out=. --go_opt=paths=source_relative \
    	   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	   $(PROTO)

clean:
	rm $(OUTPUT)
	rm $(PROTO_GEN)
