package main

import (
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	v := NewViper()
	relayService := NewRelayService()
	config := ConfigFromViper(v)
	if err := relayService.Start(config); err != nil {
		logrus.Fatalf("failed to start relay service: %+v", err)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
