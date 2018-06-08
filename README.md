# clay
Minimal server platform for gRPC and REST+Swagger APIs

Using clay you can automatically spin up HTTP handlers for your gRPC server with complete Swagger defs with a few lines of code.

## Why?
There's an excellent [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) proxy generator, 
but it requires you to spin up (at least) one proxy instance in addition to your services.
`clay` allows you to serve HTTP traffic by server instances themselves for easier debugging/testing. 

## How?
See [example server](https://github.com/utrack/clay/blob/master/doc/example/main.go).
First, generate your Go code using protoc:
```
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. --goclay_out=:. ./sum.proto
```
Then finish your gRPC service implementation as usual:

```
// SumImpl is an implementation of SummatorService.
type SumImpl struct{}

// Sum implements SummatorServer.Sum.
func (s *SumImpl) Sum(ctx context.Context, r *pb.SumRequest) (*pb.SumResponse, error) {
	if r.GetA() == 0 {
		return nil, errors.New("a is zero")
	}

	if r.GetB() == 65536 {
		panic(errors.New("we've got a problem!"))
	}

	sum := r.GetA() + r.GetB()
	return &pb.SumResponse{
		Sum: sum,
	}, nil
}
```

Then, add one method to the implementation, so it would implement the `"github.com/utrack/clay/transport".Service`:
```
// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (s *SumImpl) GetDescription() transport.ServiceDesc {
	return pb.NewSummatorServiceDesc(s)
}
```

Now, you can [run the server](https://github.com/utrack/clay/blob/master/doc/example/main.go#L68). 
Swagger definition will be served at `/swagger.json`.

`clay.Server` is easily extendable, as you can pass any options gRPC server can use, 
but if it's not extendable enough then you can use the `.GetDescription()` method 
of your implementation to register the service in your own custom server 
(see [ServiceDesc](https://github.com/utrack/clay/blob/master/transport/handlers.go#L17)).
