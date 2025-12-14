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

	commit1Hash := plumbing.NewHash("cacf82fd8a097b8266fb7a7c84f02a9a1a599a4c")
	commit2Hash := plumbing.NewHash("b20da7389753009b3cd6be1ae0cc66e5551e8916")
	t1, _ := r.CommitObject(commit1Hash)
	t2, _ := r.CommitObject(commit2Hash)
	p, _ := t1.Patch(t2)
	print(p.String())
	for _, fp := range p.FilePatches() {
		from, to := fp.Files()
		if from != nil {
			f := from.Path()
			fmt.Println(f)
		}

		if to != nil {
			t := to.Path()
			fmt.Println("file: " + t)
		}

		for _, c := range fp.Chunks() {
			fmt.Println(c.Content())
		}
	}

}
