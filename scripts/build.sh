#!/usr/bin/env bash

OS=(darwin linux)
ARCH=(amd64)

for os in ${OS[@]}; do
    for arch in ${ARCH[@]}; do
        mkdir -p release/filecollector_${os}_${arch}
        env GOOS=$os GOARCH=$arch go build -o release/filecollector_${os}_${arch}
        (cd release && zip -r filecollector_${os}_${arch}.zip filecollector_${os}_${arch})
        rm -rf release/filecollector_${os}_${arch}
    done
done
