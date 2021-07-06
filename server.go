package main

import (
	"context"
	"raftpractice/grpcproto"
)

type Server struct {
}

func (s *Server) RequestVote(ctx context.Context, in *grpcproto.RequestVoteMsg) (*grpcproto.RequestVoteReply, error) {
	return nil, nil
}

func (s *Server) AppendEntries(ctx context.Context, in *grpcproto.AppendEntriesMsg) (*grpcproto.AppendEntriesReply, error) {
	return nil, nil
}
