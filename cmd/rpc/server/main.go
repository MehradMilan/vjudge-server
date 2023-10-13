package main

import (
	"context"
	"log"
	"net"
	"os"
	"vjudge/pkg/judge"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"google.golang.org/grpc"
	"gopkg.in/src-d/go-git.v4"
)

type server struct {
	judge.UnimplementedCodeJudgeServer
}

func (s *server) JudgeCode(ctx context.Context, in *judge.SubmissionRequest) (*judge.SubmissionReply, error) {
	go judgeCodeFromGitHubURL(in.GetGithubUrl())
	return &judge.SubmissionReply{
		Submitted: true,
	}, nil
}

func cloneCodeFromGitHub(gitHubURL string) bool {
	httpsAuth := &http.BasicAuth{
		Username: "", // this can be anything except an empty string
		Password: "", // ideally, the GitHub token
	}
	_, err := git.PlainClone("../../resources", false, &git.CloneOptions{
		URL:      gitHubURL,
		Progress: os.Stdout,
		Auth:     httpsAuth,
	})
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func judgeCodeFromGitHubURL(gitHubURL string) {
	// TODO: Your logic to pull code and judge it
	cloneCodeFromGitHub(gitHubURL)
	// judge.JudgeCode("Meow")
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
