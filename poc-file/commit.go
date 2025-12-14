package main

import (
	"log/slog"

	git "github.com/go-git/go-git/v6"
)

func main() {
	// clone
	// r, err := git.PlainClone("poc/", &git.CloneOptions{
	// 	URL:      "https://github.com/tkeshun/go-git-poc",
	// 	Progress: os.Stdout, // 進捗を表示したくなければ nil でOK
	// })
	r, err := git.PlainOpen("./poc")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	wt, err := r.Worktree()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// filepath := "./poc/empty.txt"
	// f, err := os.OpenFile(
	// 	filepath,
	// 	os.O_CREATE|os.O_WRONLY,
	// 	0644,
	// )
	// if err != nil {
	// 	slog.Error("Create File Error: " + err.Error())
	// 	return
	// }
	// defer f.Close()

	// addpath := "empty.txt"
	// if _, err := wt.Add(addpath); err != nil {
	// 	slog.Error("Add Error: " + err.Error())
	// 	return
	// }
	// st, _ := wt.Status()
	// print(st.String())

	// head, err := r.Head()
	// if err != nil {
	// 	slog.Error("Head Error: " + err.Error())
	// }

	// if err := wt.Reset(&git.ResetOptions{
	// 	Commit: head.Hash(),
	// 	Mode:   git.MixedReset, // index を HEAD に戻す（作業ツリーはそのまま）
	// }); err != nil {
	// 	slog.Error(err.Error())
	// }
	// st, _ = wt.Status()
	// print(st.String())

	// if _, err := wt.Add(addpath); err != nil {
	// 	slog.Error("Add Error: " + err.Error())
	// 	return
	// }
	// st, _ = wt.Status()
	// print(st.String())

	// if _, err = wt.Commit("test", &git.CommitOptions{
	// 	Author: &object.Signature{
	// 		Name:  "go-git-invalid",
	// 		Email: "",
	// 		When:  time.Now(),
	// 	},
	// }); err != nil {
	// 	slog.Error("commit Error: " + err.Error())
	// }

	// println("commit finish")

	headRef, _ := r.Head()
	headCommit, _ := r.CommitObject(headRef.Hash())
	parent, _ := headCommit.Parent(0)
	wt.Reset(&git.ResetOptions{
		Commit: parent.Hash,
		Mode:   git.MixedReset,
	})
	println("reset commit")
}
