package main

import (
  _ "context"

  _ "github.com/caring/producers/internal/handlers"
  "github.com/caring/producers/pb"
  _ "github.com/caring/go-packages/pkg/errors"
  _ "google.golang.org/grpc/codes"
)

type service struct {
}

func (s *service) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
    l.Printf("Received: %v", in.Data)
    resp := "Data: " + in.Data
    
    return &pb.PingResponse{Data: resp}, nil
    
}