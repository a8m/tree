#!/bin/bash

while read -r line; do
    parts=(${line//\// })
    export GOOS=${parts[0]}
    export GOARCH=${parts[1]}
    echo Try GOOS=${GOOS} GOARCH=${GOARCH}
    go install
done < <(go tool dist list)
