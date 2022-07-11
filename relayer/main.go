package main

import (
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	config := configFromViper(newViper())
	emitter := newEmitter()

	relayService := NewRelayService(emitter)
	if err := relayService.Start(config); err != nil {
		logrus.Fatalf("failed to start relay service: %+v", err)
	}

	grpcServer := NewGrpcServer(emitter)
	if err := grpcServer.Start(config); err != nil {
		logrus.Fatalf("failed to start grpc server: %+v", err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
