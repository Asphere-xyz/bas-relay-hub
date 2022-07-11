package main

import (
	"context"
	"fmt"
	"github.com/Ankr-network/bas-relay-hub/relayer/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/olebedev/emitter"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
	"time"
)

type GrpcServer struct {
	emitter *emitter.Emitter

	server   *grpc.Server
	listener net.Listener

	proto.RelayHubServer
}

func NewGrpcServer(emitter *emitter.Emitter) *GrpcServer {
	return &GrpcServer{emitter: emitter}
}

func (s *GrpcServer) Start(config *Config) error {
	listener, err := net.Listen("tcp", config.Relayer.GrpcAddress)
	if err != nil {
		return errors.Wrapf(err, "can't listen address (%s)", config.Relayer.GrpcAddress)
	}
	var serverOps []grpc.ServerOption
	s.server = grpc.NewServer(serverOps...)
	s.listener = listener
	proto.RegisterRelayHubServer(s.server, s)
	go func() {
		log.Infof("gRPC server is listening on address (%s)", config.Relayer.GrpcAddress)
		if err = s.server.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()
	for _, url := range config.Relayer.RelayerUrls {
		conn, err := grpc.Dial(url)
		if err != nil {
			log.WithError(err).Warnf("failed to establish gRPC connection with relayer (%s)", url)
		}
		client := proto.NewRelayHubClient(conn)
		go func(url string) {
			for {
				err := s.listenForIncomingMessages(client)
				if err != nil {
					log.WithError(err).WithField("relayer", url).Errorf("failed to stream signatures from relayer")
				}
				time.Sleep(15 * time.Second)
			}
		}(url)
	}
	return nil
}

func (s *GrpcServer) listenForIncomingMessages(client proto.RelayHubClient) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()
	stream, err := client.SignCheckpointProof(ctx, &proto.SignCheckpointProofRequest{}, nil)
	if err != nil {
		return err
	}
	for {
		reply, err := stream.Recv()
		if err != nil {
			return err
		} else if len(reply.Validators) != len(reply.Signatures) {
			return fmt.Errorf("corrupted data, count of validators is not equal to signers (%d != %d)", len(reply.Validators), len(reply.Signatures))
		}
		signatures := make(map[common.Address][]byte)
		for i, validator := range reply.Validators {
			signatures[common.BytesToAddress(validator)] = reply.Signatures[i]
		}
		proof := &checkpointProof{
			epochNumber:    reply.EpochNumber,
			rawEpochBlock:  reply.RawEpochBlock,
			blockHash:      common.BytesToHash(reply.BlockHash),
			checkpointHash: common.BytesToHash(reply.CheckpointHash),
			signatures:     signatures,
		}
		<-s.emitter.Emit(checkpointProofReceived, proof)
	}
}

func (s *GrpcServer) SignCheckpointProof(_ *proto.SignCheckpointProofRequest, result proto.RelayHub_SignCheckpointProofServer) error {
	event := s.emitter.On(checkpointProofSigned)
	ctx := result.Context()
	for {
		select {
		case e := <-event:
			message, ok := e.Args[0].(*checkpointProof)
			if !ok {
				return fmt.Errorf("can't retrieve epoch signed message from event emitter")
			}
			reply := &proto.SignCheckpointProofReply{
				EpochNumber:    message.epochNumber,
				RawEpochBlock:  message.rawEpochBlock,
				BlockHash:      message.blockHash.Bytes(),
				CheckpointHash: message.checkpointHash.Bytes(),
			}
			for k, v := range message.signatures {
				reply.Validators = append(reply.Validators, k.Bytes())
				reply.Signatures = append(reply.Signatures, v)
			}
			err := result.Send(reply)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (s *GrpcServer) Stop() {
	s.server.Stop()
	_ = s.listener.Close()
}
