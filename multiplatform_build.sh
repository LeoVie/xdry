#!/bin/bash

#for platform in "aix/ppc64" "darwin/amd64" "darwin/arm64" "dragonfly/amd64" "freebsd/386" "freebsd/amd64" "freebsd/arm" "freebsd/arm64" "illumos/amd64" "js/wasm" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "linux/mips" "linux/mips64" "linux/mips64le" "linux/mipsle" "linux/ppc64" "linux/ppc64le" "linux/riscv64" "linux/s390x" "netbsd/386" "netbsd/amd64" "netbsd/arm" "netbsd/arm64" "openbsd/386" "openbsd/amd64" "openbsd/arm" "openbsd/arm64" "openbsd/mips64" "plan9/386" "plan9/amd64" "plan9/arm" "solaris/amd64" "windows/386" "windows/amd64" "windows/arm" "windows/arm64"
for platform in "linux/arm64"
do
  IFS='/' read -ra platformArray <<< "$platform"

  binaryName=$(pwd)/build/xdry-${platformArray[0]}-${platformArray[1]}

  if [ ${platformArray[0]} == "windows" ]
  then
    binaryName=${binaryName}.exe
  fi

  make build goos=${platformArray[0]} goarch=${platformArray[1]} binaryFile=${binaryName} version=${1}
done