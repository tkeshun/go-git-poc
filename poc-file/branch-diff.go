package main

import (
	"fmt"
	"log/slog"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

func main() {
	r, err := git.PlainOpen("./poc")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	ref1, err := r.Reference(plumbing.NewBranchReferenceName("main"), true)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	ref2, err := r.Reference(plumbing.NewBranchReferenceName("test"), true)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	commit1Hash := ref1.Hash()
	commit2Hash := ref2.Hash()
	t1, _ := r.CommitObject(commit1Hash)
	t2, _ := r.CommitObject(commit2Hash)
	p, _ := t1.Patch(t2)
	print(p.String())
	for _, fp := range p.FilePatches() {
		_, to := fp.Files()
		if to != nil {
			t := to.Path()
			fmt.Println("file: " + t)
		}

		for _, c := range fp.Chunks() {
			fmt.Println(c.Content())
		}
	}

}
