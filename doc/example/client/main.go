package main

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/utrack/clay/doc/example/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:12345", grpc.WithInsecure())
	if err != nil {
		logrus.Fatal(err)
	}
	client := sumpb.NewSummatorClient(conn)
	rsp, err := client.Sum(context.Background(), &sumpb.SumRequest{A: 1, B: 2})
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(rsp)
	}

	rsp, err = client.Sum(context.Background(), &sumpb.SumRequest{A: 0, B: 2})
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(rsp)
	}

	rsp, err = client.Sum(context.Background(), &sumpb.SumRequest{A: 1, B: 65536})
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(rsp)
	}

}
