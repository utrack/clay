package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	desc "github.com/utrack/clay/doc/example/pb"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:12345", grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}
	client := desc.NewSummatorClient(conn)

	makeRequest(client, &desc.SumRequest{A: 1, B: &desc.NestedB{B: 2}})
	makeRequest(client, &desc.SumRequest{A: 1, B: &desc.NestedB{B: 65536}})
	makeRequest(client, &desc.SumRequest{A: 0, B: &desc.NestedB{B: 2}})
	makeRequest(client, &desc.SumRequest{A: 1, B: &desc.NestedB{B: 3}})
}

func makeRequest(client desc.SummatorClient, req *desc.SumRequest) {
	logrus.Infof("Request `%v`", req)
	rsp, err := client.Sum(context.Background(), req)
	if err != nil {
		logrus.Errorf("server responded with error: `%v`", err)
	} else {
		logrus.Infof("Response: `%v`", rsp)
	}
}
