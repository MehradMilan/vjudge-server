package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	judge.JudgeCode("Meow")
}

func main() {
	r := gin.Default()
	r.POST("/webhook", Webhook)
	srv := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT)
	<-quit
	slog.Info("Shutting down the server...")
	_ = srv.Shutdown(context.Background())
}

type githubPayload struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Repository struct {
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Name              string      `json:"name"`
			Email             interface{} `json:"email"`
			Login             string      `json:"login"`
			Id                int         `json:"id"`
			NodeId            string      `json:"node_id"`
			AvatarUrl         string      `json:"avatar_url"`
			GravatarId        string      `json:"gravatar_id"`
			Url               string      `json:"url"`
			HtmlUrl           string      `json:"html_url"`
			FollowersUrl      string      `json:"followers_url"`
			FollowingUrl      string      `json:"following_url"`
			GistsUrl          string      `json:"gists_url"`
			StarredUrl        string      `json:"starred_url"`
			SubscriptionsUrl  string      `json:"subscriptions_url"`
			OrganizationsUrl  string      `json:"organizations_url"`
			ReposUrl          string      `json:"repos_url"`
			EventsUrl         string      `json:"events_url"`
			ReceivedEventsUrl string      `json:"received_events_url"`
			Type              string      `json:"type"`
			SiteAdmin         bool        `json:"site_admin"`
		} `json:"owner"`
		HtmlUrl                  string        `json:"html_url"`
		Description              string        `json:"description"`
		Fork                     bool          `json:"fork"`
		Url                      string        `json:"url"`
		ForksUrl                 string        `json:"forks_url"`
		KeysUrl                  string        `json:"keys_url"`
		CollaboratorsUrl         string        `json:"collaborators_url"`
		TeamsUrl                 string        `json:"teams_url"`
		HooksUrl                 string        `json:"hooks_url"`
		IssueEventsUrl           string        `json:"issue_events_url"`
		EventsUrl                string        `json:"events_url"`
		AssigneesUrl             string        `json:"assignees_url"`
		BranchesUrl              string        `json:"branches_url"`
		TagsUrl                  string        `json:"tags_url"`
		BlobsUrl                 string        `json:"blobs_url"`
		GitTagsUrl               string        `json:"git_tags_url"`
		GitRefsUrl               string        `json:"git_refs_url"`
		TreesUrl                 string        `json:"trees_url"`
		StatusesUrl              string        `json:"statuses_url"`
		LanguagesUrl             string        `json:"languages_url"`
		StargazersUrl            string        `json:"stargazers_url"`
		ContributorsUrl          string        `json:"contributors_url"`
		SubscribersUrl           string        `json:"subscribers_url"`
		SubscriptionUrl          string        `json:"subscription_url"`
		CommitsUrl               string        `json:"commits_url"`
		GitCommitsUrl            string        `json:"git_commits_url"`
		CommentsUrl              string        `json:"comments_url"`
		IssueCommentUrl          string        `json:"issue_comment_url"`
		ContentsUrl              string        `json:"contents_url"`
		CompareUrl               string        `json:"compare_url"`
		MergesUrl                string        `json:"merges_url"`
		ArchiveUrl               string        `json:"archive_url"`
		DownloadsUrl             string        `json:"downloads_url"`
		IssuesUrl                string        `json:"issues_url"`
		PullsUrl                 string        `json:"pulls_url"`
		MilestonesUrl            string        `json:"milestones_url"`
		NotificationsUrl         string        `json:"notifications_url"`
		LabelsUrl                string        `json:"labels_url"`
		ReleasesUrl              string        `json:"releases_url"`
		DeploymentsUrl           string        `json:"deployments_url"`
		CreatedAt                int           `json:"created_at"`
		UpdatedAt                time.Time     `json:"updated_at"`
		PushedAt                 int           `json:"pushed_at"`
		GitUrl                   string        `json:"git_url"`
		SshUrl                   string        `json:"ssh_url"`
		CloneUrl                 string        `json:"clone_url"`
		SvnUrl                   string        `json:"svn_url"`
		Homepage                 string        `json:"homepage"`
		Size                     int           `json:"size"`
		StargazersCount          int           `json:"stargazers_count"`
		WatchersCount            int           `json:"watchers_count"`
		Language                 string        `json:"language"`
		HasIssues                bool          `json:"has_issues"`
		HasProjects              bool          `json:"has_projects"`
		HasDownloads             bool          `json:"has_downloads"`
		HasWiki                  bool          `json:"has_wiki"`
		HasPages                 bool          `json:"has_pages"`
		HasDiscussions           bool          `json:"has_discussions"`
		ForksCount               int           `json:"forks_count"`
		MirrorUrl                interface{}   `json:"mirror_url"`
		Archived                 bool          `json:"archived"`
		Disabled                 bool          `json:"disabled"`
		OpenIssuesCount          int           `json:"open_issues_count"`
		License                  interface{}   `json:"license"`
		AllowForking             bool          `json:"allow_forking"`
		IsTemplate               bool          `json:"is_template"`
		WebCommitSignoffRequired bool          `json:"web_commit_signoff_required"`
		Topics                   []interface{} `json:"topics"`
		Visibility               string        `json:"visibility"`
		Forks                    int           `json:"forks"`
		OpenIssues               int           `json:"open_issues"`
		Watchers                 int           `json:"watchers"`
		DefaultBranch            string        `json:"default_branch"`
		Stargazers               int           `json:"stargazers"`
		MasterBranch             string        `json:"master_branch"`
		Organization             string        `json:"organization"`
	} `json:"repository"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Organization struct {
		Login            string      `json:"login"`
		Id               int         `json:"id"`
		NodeId           string      `json:"node_id"`
		Url              string      `json:"url"`
		ReposUrl         string      `json:"repos_url"`
		EventsUrl        string      `json:"events_url"`
		HooksUrl         string      `json:"hooks_url"`
		IssuesUrl        string      `json:"issues_url"`
		MembersUrl       string      `json:"members_url"`
		PublicMembersUrl string      `json:"public_members_url"`
		AvatarUrl        string      `json:"avatar_url"`
		Description      interface{} `json:"description"`
	} `json:"organization"`
	Sender struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"sender"`
	Created bool        `json:"created"`
	Deleted bool        `json:"deleted"`
	Forced  bool        `json:"forced"`
	BaseRef interface{} `json:"base_ref"`
	Compare string      `json:"compare"`
	Commits []struct {
		Id        string    `json:"id"`
		TreeId    string    `json:"tree_id"`
		Distinct  bool      `json:"distinct"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"commits"`
	HeadCommit struct {
		Id        string    `json:"id"`
		TreeId    string    `json:"tree_id"`
		Distinct  bool      `json:"distinct"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"head_commit"`
}
