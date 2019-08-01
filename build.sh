#!/bin/bash

SYSTEMS=(windows linux freebsd darwin)
ARCHS=(amd64 386)

clean=$(git status --porcelain --untracked-files=no)
if [ -n "$clean" ]; then
   echo "There are uncommited changes"
   exit 1
fi

rev=$(git describe --tags --always)
if [ -e "$rev" ]; then
    rm -rf "$rev"
fi
mkdir -p "./bin/$rev"

echo "Revision is ${rev}"
for os in ${SYSTEMS[@]}; do
    for arch in ${ARCHS[@]}; do
        echo "Building GOOS=$os GOARCH=$arch..."
        out="prom_exporter_${os}_${arch}"
        if [ $os = "windows" ]; then
            out="${out}.exe"
        fi
        CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o "./bin/${rev}/${out}" ./cmd/prom_exporter
        (
            cd "./bin/$rev"
            sha256sum "$out" > "$out".sha256
        )
    done
done

(
    cd "./bin/$rev"
    sha256sum -c --strict *.sha256
)