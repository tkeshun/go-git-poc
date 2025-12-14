package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
)

func main() {
	r, err := git.PlainOpen("./poc")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	remoteRefName := plumbing.NewRemoteReferenceName("origin", "test")
	beforeRef, err := r.Reference(remoteRefName, true)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	beforeHash := beforeRef.Hash()
	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			"+refs/heads/test:refs/remotes/origin/test",
		},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		slog.Error(err.Error())
		return
	}

	afterRef, err := r.Reference(remoteRefName, true)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	afterHash := afterRef.Hash()
	if beforeHash == afterHash {
		fmt.Println("diff: no changes")
		return
	}

	beforeCommit, err := r.CommitObject(beforeHash)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	afterCommit, err := r.CommitObject(afterHash)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	patch, err := beforeCommit.Patch(afterCommit)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	filepaths := make([]string, 0, 10)
	for _, fp := range patch.FilePatches() {
		_, to := fp.Files()

		if to != nil {
			t := to.Path()
			if !strings.HasSuffix(t, ".sql") {
				continue
			}

			filepaths = append(filepaths, t)
			fmt.Println("file: " + t)
		}
	}

	if len(filepaths) == 0 {
		fmt.Println("psql: no .sql files to apply")
		return
	}

	args := []string{
		"-h", "localhost",
		"-p", "5432",
		"-U", "app",
		"-d", "appdb",
		"-v", "ON_ERROR_STOP=1",
		"-1",
	}
	for _, fp := range filepaths {
		args = append(args, "-f", fp)
	}

	cmd := exec.Command("psql", args...)
	cmd.Dir = "./poc"
	cmd.Env = append(os.Environ(), "PGPASSWORD=app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		slog.Error("psql failed", "err", err)
		return
	}
}
