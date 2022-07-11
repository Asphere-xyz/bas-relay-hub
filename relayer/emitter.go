package main

import "github.com/olebedev/emitter"

func newEmitter() *emitter.Emitter {
	return emitter.New(10000)
}

var checkpointProofSigned = "checkpointProofSigned"
var checkpointProofReceived = "checkpointProofReceived"
