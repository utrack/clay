# clay
[![Build Status](https://travis-ci.org/utrack/clay.svg?branch=master)](https://travis-ci.org/utrack/clay)

Minimal server platform for gRPC and REST+Swagger APIs in Go

Using clay you can automatically spin up HTTP handlers for your gRPC server with
complete Swagger defs with a few lines of code.

This project provides the HTTP+Swagger handler generator and optional server that you
can use to serve your handlers via any protocol.

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
