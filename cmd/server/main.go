package main

import (
	"context"
	"log"
	"net"
	"vjudge/pkg/judge"

	"google.golang.org/grpc"
)

type server struct {
	judge.UnimplementedCodeJudgeServer
}

func (s *server) JudgeCode(ctx context.Context, in *judge.SubmissionRequest) (*judge.JudgementReply, error) {
	judgement := judgeCodeFromGitHubURL(in.GetGithubUrl())
	return judgement, nil
}

func judgeCodeFromGitHubURL(gitHubURL string) *judge.JudgementReply {
	// TODO: Your logic to pull code and judge it
	return &judge.JudgementReply{
		Score: 90,
		TestCaseResults: []*judge.TestCaseResult{
			{Id: 0, Name: "Test1", Passed: true},
			{Id: 1, Name: "Test2", Passed: true},
		},
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	judge.RegisterCodeJudgeServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
