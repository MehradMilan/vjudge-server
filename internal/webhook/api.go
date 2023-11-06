package webhook

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"vjudge/pkg/judge"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type CloneRepoData struct {
	RepoUrl   string
	OwnerName string
	RepoName  string
}

func CloneRepositoryFromGithub(repoData CloneRepoData) (string, error) {
	httpsAuth := &http.BasicAuth{
		Username: config.GitUsername, // this can be anything except an empty string
		Password: config.GitPassword, // ideally, the GitHub token
	}

	tmpDir, tmpDirErr := os.MkdirTemp(config.TmpDirectory, repoData.OwnerName+"-"+repoData.RepoName)
	if tmpDirErr != nil {
		log.Fatal(tmpDirErr)
		return "", tmpDirErr
	}
	_, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:      repoData.RepoUrl,
		Progress: os.Stdout,
		Auth:     httpsAuth,
	})
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return tmpDir, nil
}

func RunJudgeProcess(payload githubPayload) {
	if len(os.Args) > 2 {
		readConfig(os.Args[2])
	} else {
		readConfig("config/config-judge.json")
	}
	repoData := &CloneRepoData{
		RepoUrl:   payload.Repository.Url,
		OwnerName: payload.Repository.Owner.Name,
		RepoName:  payload.Repository.Name,
	}
	tmpDir, err := CloneRepositoryFromGithub(*repoData)
	if err != nil {
		log.Fatal(err)
	}

	defer cleanupTempDir(tmpDir)

	repo, err := git.PlainOpen(tmpDir)
	if err != nil {
		log.Fatal(err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	for _, hwDir := range config.HWDirectories {
		judgeResult := judge.JudgeCode(filepath.Join(tmpDir, config.SRCDirectory, hwDir), filepath.Join(config.TestDirectory, hwDir))
		// Write grade.txt file
		gradeFilePath := filepath.Join(tmpDir, config.SRCDirectory, hwDir, "grade.txt")
		writeErr := writeGradeFile(judgeResult, gradeFilePath)
		if writeErr != nil {
			log.Fatal(err)
		}
	}

	// Commit, push, and cleanup
	err = commitAndPushChanges(repo, worktree, "Updated `grade.txt`")
	if err != nil {
		log.Fatal(err)
	}
}

func writeGradeFile(judgeResult *judge.JudgeResult, filePath string) error {
	// Open the file in write mode. Create it if it doesn't exist.
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write human-friendly data to the file
	_, err = file.WriteString("Status: " + judgeResult.Status.Message + "\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString("Passed: " + strconv.FormatBool(judgeResult.Passed) + "\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString("Tests Count: " + strconv.Itoa(judgeResult.TestsCount) + "\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString("Passed Tests Count: " + strconv.Itoa(judgeResult.PassedTestsCount) + "\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString("Score: " + strconv.FormatFloat(judgeResult.Score, 'f', -1, 64) + "\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString("Test Cases:\n")
	if err != nil {
		return err
	}

	// Write test cases data
	for _, testcase := range judgeResult.Testcases {
		_, err = file.WriteString("- Name: " + testcase.Name + "\n")
		if err != nil {
			return err
		}
		_, err = file.WriteString("  Passed: " + strconv.FormatBool(testcase.Passed) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func commitAndPushChanges(repo *git.Repository, worktree *git.Worktree, commitMessage string) error {
	// Stage all changes (similar to `git add .`)
	_, err := worktree.Add(".")
	if err != nil {
		return err
	}

	// Commit changes
	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		All:    true,
		Author: &object.Signature{Name: config.JudgeName, Email: config.JudgeEmail, When: time.Now()},
	})
	if err != nil {
		return err
	}

	httpAuth := &http.BasicAuth{
		Username: config.GitUsername, // this can be anything except an empty string
		Password: config.GitPassword, // ideally, the GitHub token
	}

	// Push changes to the remote repository
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       httpAuth,
	})
	if err != nil {
		return err
	}

	return nil
}

func cleanupTempDir(tmpDir string) error {
	// Remove the temporary directory and its contents
	err := os.RemoveAll(tmpDir)
	if err != nil {
		return err
	}

	return nil
}
