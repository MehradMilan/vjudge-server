package main

import (
	"context"
	"log"
	"time"
	"vjudge/pkg/judge"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := judge.NewCodeJudgeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.JudgeCode(ctx, &judge.SubmissionRequest{GithubUrl: "https://github.com/MehradMilan/"})
	if err != nil {
		log.Fatalf("could not judge code: %v", err)
	}
	log.Printf("Score: %d", r.GetScore())
	for _, tc := range r.GetTestCaseResults() {
		log.Printf("Test Case %s: Passed=%v", tc.GetName(), tc.GetPassed())
	}
}
