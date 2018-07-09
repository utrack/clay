#!/usr/bin/env bash

set -ex

die() {
    echo "$@" >&2
    exit 1
}

cd /home/travis

case "$PROTOBUF_VERSION" in
3*)
    basename=protoc-$PROTOBUF_VERSION-linux-x86_64
    wget https://github.com/google/protobuf/releases/download/v$PROTOBUF_VERSION/$basename.zip
    unzip $basename.zip
    ;;
*)
    die "unknown protobuf version: $PROTOBUF_VERSION"
    ;;
esac
