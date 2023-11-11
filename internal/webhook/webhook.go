package webhook

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"vjudge/pkg/util"

	"github.com/gin-gonic/gin"
)

var Secret []byte

// Webhook is the function which gin should call when GitHub accesses it
func Webhook(c *gin.Context) {
	if len(os.Args) > 2 {
		readConfig(os.Args[2])
	} else {
		readConfig("config/config-judge.json")
	}

	event := c.GetHeader("X-GitHub-Event")
	logger := slog.With(
		slog.String("id", c.GetHeader("X-GitHub-Delivery")),
		slog.String("event", event),
		slog.String("ip", c.GetHeader("CF-Connecting-IP")))
	// Read the body to validate hash
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.With(util.SlogError(err)).Error("cannot read body of request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
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

	// Don't accept grading pushes by the judge
	if payload.Pusher.Name == config.GitUsername {
		logger.With(slog.String("judged", payload.Ref)).Debug("ignored grading push")
		return
	}

	homeworkName, err := getHomeworkName(payload.Repository.Name)
	if err != nil {
		logger.With(slog.String("name", payload.Ref)).Warn(err.Error())
		return
	}

	homework, ok := getHomework(homeworkName)
	if !ok {
		logger.With(slog.String("name", payload.Ref)).Warn("no homework with the given name")
		return
	}

	RunJudgeProcess(payload, homeworkName, homework)
	// Push the job
}

func getHomeworkName(repositoryName string) (string, error) {
	firstDashIndex := strings.Index(repositoryName, "-")
	if firstDashIndex == -1 {
		return "", errors.New("could not extract homework name")
	}

	return repositoryName[:firstDashIndex], nil
}

func getHomework(name string) (*Homework, bool) {
	for _, homework := range config.Homeworks {
		if homework.Name == name {
			return &homework, true
		}
	}

	return nil, false
}
