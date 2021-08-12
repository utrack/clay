package main

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	desc "github.com/utrack/clay/doc/example/pb"
)

func main() {
	client := desc.NewSummatorHTTPClient(http.DefaultClient, "http://localhost:12345")

	makeRequest(client, &desc.SumRequest{A: 1, B: &desc.NestedB{B: 2}})
	makeRequest(client, &desc.SumRequest{A: 1, B: &desc.NestedB{B: 65536}})
	makeRequest(client, &desc.SumRequest{A: 0, B: &desc.NestedB{B: 2}})
	makeRequest(client, &desc.SumRequest{A: 1, B: &desc.NestedB{B: 3}})
}

func makeRequest(client *desc.Summator_httpClient, req *desc.SumRequest) {
	logrus.Infof("Request `%v`", req)
	rsp, err := client.Sum(context.Background(), req)
	if err != nil {
		logrus.Errorf("server responded with error: `%v`", err)
	} else {
		logrus.Infof("Response: `%v`", rsp)
	}
}
