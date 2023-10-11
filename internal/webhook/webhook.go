package webhook

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"vjudge/pkg/judge"
	"vjudge/pkg/util"

	"github.com/gin-gonic/gin"
)

var Secret []byte

// Webhook is the function which gin should call when GitHub accesses it
func Webhook(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")
	logger := slog.With(
		slog.String("id", c.GetHeader("X-GitHub-Delivery")),
		slog.String("event", event),
		slog.String("ip", c.GetHeader("CF-Connecting-IP")))
	// Read the body to validate hash
	body := make([]byte, 64*1024)
	readBytes, err := c.Request.Body.Read(body)
	if err != nil {
		logger.With(util.SlogError(err)).Error("cannot read body of request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	body = body[:readBytes]
	// Validate the hash
	expectedHash := c.GetHeader("X-Hub-Signature-256")
	if len(expectedHash) > 7 {
		expectedHash = expectedHash[7:]
	}
	if !util.VerifyGithubSignature(Secret, body, expectedHash) {
		logger.Warn("signature mismatch")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// Check action
	if event != "push" {
		logger.Warn("unknown event: " + event)
		return
	}
	// Parse the body
	var payload githubPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		logger.With(util.SlogError(err)).Error("cannot parse payload")
		return
	}
	// Accept main pushes only
	if payload.Ref != "refs/heads/main" {
		logger.With(slog.String("ref", payload.Ref)).Debug("ignored non main ref")
		return
	}

	// Push the job
	judge.JudgeCode("/home/mehrad/Projects/vjudge/vjudge-core/test/srcdir/")
}