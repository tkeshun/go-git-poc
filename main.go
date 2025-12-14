package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func main() {
	r, err := git.PlainOpen("./poc")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	remoteRefName := plumbing.NewRemoteReferenceName("origin", "test")

	// before は無い場合もある（初回など）のでゼロハッシュ扱いにする
	beforeHash := plumbing.ZeroHash
	if beforeRef, err := r.Reference(remoteRefName, true); err == nil {
		beforeHash = beforeRef.Hash()
	}

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

	if beforeHash == afterHash && !beforeHash.IsZero() {
		fmt.Println("diff: no changes")
		return
	}

	afterCommit, err := r.CommitObject(afterHash)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// 変更された .sql のパス一覧を作る
	filepaths, err := changedSQLPaths(r, beforeHash, afterHash)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if len(filepaths) == 0 {
		fmt.Println("psql: no .sql files to apply")
		return
	}

	// afterCommit の内容から該当SQLを一時ディレクトリに書き出す
	tmpDir, err := os.MkdirTemp("", "sqlapply-*")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer os.RemoveAll(tmpDir)

	sqlFilesAbs := make([]string, 0, len(filepaths))
	for _, fp := range filepaths {
		f, err := afterCommit.File(fp)
		if err != nil {
			slog.Error("file not found in after commit", "path", fp, "err", err)
			return
		}

		body, err := f.Contents()
		if err != nil {
			slog.Error("read file contents failed", "path", fp, "err", err)
			return
		}

		dst := filepath.Join(tmpDir, filepath.FromSlash(fp))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			slog.Error("mkdir failed", "dir", filepath.Dir(dst), "err", err)
			return
		}
		if err := os.WriteFile(dst, []byte(body), 0o644); err != nil {
			slog.Error("write file failed", "path", dst, "err", err)
			return
		}

		fmt.Println("file:", fp)
		sqlFilesAbs = append(sqlFilesAbs, dst)
	}

	args := []string{
		"-h", "localhost",
		"-p", "5432",
		"-U", "app",
		"-d", "appdb",
		"-v", "ON_ERROR_STOP=1",
		"-1",
	}
	for _, abs := range sqlFilesAbs {
		args = append(args, "-f", abs) // 絶対パスで渡す
	}

	cmd := exec.Command("psql", args...)
	cmd.Env = append(os.Environ(), "PGPASSWORD=app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		slog.Error("psql failed", "err", err)
		return
	}
}

func changedSQLPaths(r *git.Repository, beforeHash, afterHash plumbing.Hash) ([]string, error) {
	afterCommit, err := r.CommitObject(afterHash)
	if err != nil {
		return nil, err
	}

	// 初回などで before が無い場合：afterCommit に含まれる .sql を全適用
	if beforeHash.IsZero() {
		return allSQLInCommit(afterCommit)
	}

	beforeCommit, err := r.CommitObject(beforeHash)
	if err != nil {
		return nil, err
	}

	patch, err := beforeCommit.Patch(afterCommit)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, fp := range patch.FilePatches() {
		_, to := fp.Files()
		if to == nil {
			continue
		}
		p := to.Path()
		if strings.HasSuffix(p, ".sql") {
			paths = append(paths, p)
		}
	}

	sort.Strings(paths)
	return paths, nil
}

func allSQLInCommit(c *object.Commit) ([]string, error) {
	tree, err := c.Tree()
	if err != nil {
		return nil, err
	}

	var paths []string
	err = tree.Files().ForEach(func(f *object.File) error {
		if strings.HasSuffix(f.Name, ".sql") {
			paths = append(paths, f.Name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(paths)
	return paths, nil
}
