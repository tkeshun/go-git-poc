package main

import (
	"fmt"
	"log/slog"
	"os"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func main() {
	// clone
	r, err := git.PlainClone("poc/", &git.CloneOptions{
		URL:      "https://github.com/tkeshun/go-git-poc",
		Progress: os.Stdout, // 進捗を表示したくなければ nil でOK
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	ref, err := r.Head()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err := cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	}); err != nil {
		slog.Error(err.Error())
		return
	}

	if err != nil {
		slog.Error(err.Error())
		return
	}
}
