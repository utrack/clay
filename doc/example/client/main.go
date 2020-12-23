package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/utrack/clay/doc/example/pb"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:12345", grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}
	client := sumpb.NewSummatorClient(conn)
	rsp, err := client.Sum(context.Background(), &sumpb.SumRequest{A: 1, B: &sumpb.NestedB{B: 2}})
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(rsp)
	}

	rsp, err = client.Sum(context.Background(), &sumpb.SumRequest{A: 0, B: &sumpb.NestedB{B: 2}})
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(rsp)
	}

	rsp, err = client.Sum(context.Background(), &sumpb.SumRequest{A: 1, B: &sumpb.NestedB{B: 65536}})
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(rsp)
	}

}
