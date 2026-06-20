package loader

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/emyasa/tools-to-go/realtime-config/config-loader/internal/redis"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func LoadFromGit(
	pemFilePath string,
	repoURL string,
	branch string,
	redisAddr string,
	channels ...string,
) {
	auth, _ := ssh.NewPublicKeysFromFile("git", pemFilePath, "")
	repoDir := "/tmp/realtime-config/config-loader/" + branch + "/"
	ensureRepo(auth, repoURL, branch, repoDir)

	rdb := redis.New(redisAddr)
	entries := readRepoFiles(repoDir)
	redis.SetEntries(context.Background(), rdb, "", entries)

	for _, ch := range channels {
		rdb.Publish(context.Background(), ch, "")
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
