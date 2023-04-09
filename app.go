package main

import (
	"coh3-replay-manager-go/modules/replay"
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

type ListResult struct {
	Replays      map[string]string // Replace 'Replay' with the actual type of 'replays'
	LocalReplays map[string]string // Replace 'Replay' with the actual type of 'localReplays'
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("App started")
}

func (a *App) shutdown(ctx context.Context) {}

func (a *App) List() ListResult {
	replays := replay.ListDownloaded()
	localReplays := replay.ListLocal()

	return ListResult{
		Replays:      replays,
		LocalReplays: localReplays,
	}
}

func (a *App) Play(id string) {
	fileName := fmt.Sprintf("downloaded-replay-%s.rec", id)
	replay.Play(fileName)
}

func (a *App) PlayLocal(filename string) {
	replay.Play(filename)
}

func (a *App) Remove(id string) {
	fileName := fmt.Sprintf("downloaded-replay-%s.rec", id)
	replay.Remove(fileName)
}

func (a *App) RemoveLocal(filename string) {
	replay.Remove(filename)
}
