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
	r, err := c.JudgeCode(ctx, &judge.SubmissionRequest{GithubUrl: "https://github.com/Org218/asm-chart-sorousherafat"})
	if err != nil {
		log.Fatalf("could not judge code: %v", err)
	}
	log.Printf("Submission: %t", r.GetSubmitted())
}
