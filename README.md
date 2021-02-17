# clay
[![Build Status](https://travis-ci.org/utrack/clay.svg?branch=master)](https://travis-ci.org/utrack/clay)

Minimal server platform for gRPC and REST+Swagger APIs in Go

Using clay you can automatically spin up HTTP handlers for your gRPC server with
complete Swagger defs with a few lines of code.

This project provides the HTTP+Swagger handler generator and optional server that you
can use to serve your handlers via any protocol.

# Deprecation / Looking for maintainers
This project supports only legacy Protobuf API ([golang/protobuf](https://github.com/golang/protobuf) and protoc-gen-go <v1.20.0).
These APIs were superseded by module [google.golang.org/protobuf](https://pkg.go.dev/mod/google.golang.org/protobuf).

Read more about this migration here: https://blog.golang.org/protobuf-apiv2 .

[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) is an excellent alternative to clay.

[golang/protobuf](https://github.com/golang/protobuf) >=v1.40.0 (a shim between APIs v1/v2) makes clay compatible with v2 APIs, though there's no promises on future compatibility.

It is possible to upgrade clay to v2 APIs, but there's no other maintainers and my own work has gone elsewhere. PM me or fork the repo if you want to become a new maintainer.

## Requirements

Since new [Semantic Import Versioning](https://research.swtch.com/vgo-import) is used, you are required to
use [Go1.10.3+](https://golang.org/doc/devel/release.html#go1.10)

## How?
Check out an [example server](https://github.com/utrack/clay/wiki/Build-and-run-an-example-SummatorService-using-clay-Server)
for a quick start if you're experienced with gRPC, or dive into [step-by-step docs](https://github.com/utrack/clay/wiki/Describe-and-create-your-own-API)
for a full guide.

## Contributing
You may contribute in several ways like creating new features, fixing bugs, 
improving documentation and/or examples using GitHub pull requests.
