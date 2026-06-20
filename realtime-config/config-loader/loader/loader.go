package loader

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/redis/go-redis/v9"
)

type GitOptions struct {
	pemFilePath string
	repoURL string
	branch string
}

func NewGitOptions(
	pemFilePath string,
	repoURL string,
	branch string,
) (*GitOptions, error) {
	if pemFilePath == "" {
		return nil, errors.New("pemFilePath is required")
	}

	if repoURL == "" {
		return nil, errors.New("repoURL is required")
	}

	if branch == "" {
		return nil, errors.New("branch is required")
	}

	return &GitOptions{
		pemFilePath: pemFilePath,
		repoURL: repoURL,
		branch: branch,
	}, nil
}

type RedisOptions struct {
	client *redis.Client
	prefixKey string
	channels []string
}

func NewRedisOptions (
	client *redis.Client,
	prefixKey string,
	channels []string,
) (*RedisOptions, error) {
	if client == nil {
		return nil, errors.New("client is required")
	}

	return &RedisOptions{
		client: client,
		prefixKey: prefixKey,
		channels: channels,
	}, nil
}

func LoadFromGit(
	ctx context.Context,
	gitOptions *GitOptions,
	redisOptions *RedisOptions,
) {
	if gitOptions == nil {
		fmt.Println("gitOptions is required")
		return
	}

	if redisOptions == nil {
		fmt.Println("redisOptions is required")
		return
	}

	auth, _ := ssh.NewPublicKeysFromFile("git", gitOptions.pemFilePath, "")
	repoDir := "/tmp/realtime-config/config-loader/" + gitOptions.branch + "/"
	ensureRepo(auth, gitOptions.repoURL, gitOptions.branch, repoDir)

	entries := readRepoFiles(repoDir)
	setEntries(ctx, redisOptions.client, redisOptions.prefixKey, entries)

	for _, ch := range redisOptions.channels {
		redisOptions.client.Publish(context.Background(), ch, "")
	}
}

func ensureRepo(auth *ssh.PublicKeys, repoURL string, branch string, repoDir string) {
	_, err := git.PlainClone(repoDir, false, &git.CloneOptions{
		URL: repoURL,
		Auth: auth,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})

	if err == git.ErrRepositoryAlreadyExists {
		r, _ := git.PlainOpen(repoDir)
		w, _ := r.Worktree()
		w.Pull(&git.PullOptions{
			Auth: auth,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
		})

		return
	}

	if err != nil {
		panic(err)
	}
}

func readRepoFiles(targetPath string) (map[string]string) {
	files := make(map[string]string)
	ignorePatterns := []string{".*"}

	matchesIgnore := func(name string) bool {
		for _, pattern := range ignorePatterns {
			matched, _ := filepath.Match(pattern, name)
			if matched {
				return true
			}
		}

		return false
	}

	filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := d.Name()
		if matchesIgnore(name) {
			if d.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if d.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(targetPath, path)
		if err != nil {
			return err
		}

		files[relPath] = string(data)
		return nil
	})

	return files
}

func setEntries(ctx context.Context, rdb *redis.Client, prefix string, data map[string]string) error {
	pipe := rdb.Pipeline()

	for k, v := range data {
		pipe.Set(ctx, prefix+k, v, 0)
	}

	_, err := pipe.Exec(ctx)
	return err
}
