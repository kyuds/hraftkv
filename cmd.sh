#!/bin/bash

NUM="$1"

echo "./hraftkv -aaddr=localhost:10000 -raddr=localhost:20000 -id=node0"

for (( i=1; i<$NUM; i++ )); do
    echo "./hraftkv -aaddr=localhost:1000$i -raddr=localhost:2000$i -id=node$i -join=localhost:10000"
done
