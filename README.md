# clay

[![Build Status](https://travis-ci.org/utrack/clay.svg?branch=master)](https://travis-ci.org/utrack/clay)

Minimal server platform for gRPC and REST+Swagger APIs in Go

Using clay you can automatically spin up HTTP handlers for your gRPC server with
complete Swagger defs with a few lines of code.

This project provides the HTTP+Swagger handler generator and optional server that you
can use to serve your handlers via any protocol.

# protobuf generator compatibility

clay/v3 uses new protobuf generator [google.golang.org/protobuf](https://pkg.go.dev/mod/google.golang.org/protobuf).
If you're using old generator ([golang/protobuf](https://github.com/golang/protobuf) and protoc-gen-go <v1.20.0) then consider using `clay/v2` instead.

Read more about this migration here: https://blog.golang.org/protobuf-apiv2 .

# Migration from v2

Migration from v2 to v3 is intended to be as straightforward as possible - just replace generator v2 with v3 and you're good to go.
Don't forget to change your protobuf generator as well - see previous section for details.

NB: to use new protobuf generator, you need to change `go_package` directives to be an absolute module path (i.e. `github.com/foo/bar/baz` instead of `baz`).

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
