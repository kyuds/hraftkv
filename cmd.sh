#!/bin/bash

NUM="$1"

for (( i=0; i<$NUM; i++ )); do
    echo "./hraftkv -mem -aaddr=localhost:1000$i -raddr=localhost:2000$i -id=node$i"
done
