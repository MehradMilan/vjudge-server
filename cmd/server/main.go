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

func (s *server) JudgeCode(ctx context.Context, in *judge.SubmissionRequest) (*judge.JudgementReply, error) {
	judgement := judgeCodeFromGitHubURL(in.GetGithubUrl())
	return judgement, nil
}

func cloneCodeFromGitHub(gitHubURL string) bool {
	httpsAuth := &http.BasicAuth{
		Username: "MehradMilan",                              // this can be anything except an empty string
		Password: "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", // ideally, your GitHub token
	}
	// urlParts := strings.Split(gitHubURL, "/")
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

// func authenticateInGitHub(username string, passToken string) bool {
// }

func judgeCodeFromGitHubURL(gitHubURL string) *judge.JudgementReply {
	// TODO: Your logic to pull code and judge it
	judge.JudgeCode("Meow")
	// cloneCodeFromGitHub(gitHubURL)
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
