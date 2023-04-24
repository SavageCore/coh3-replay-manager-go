package main

import (
	"coh3-replay-manager-go/modules/game"
	"coh3-replay-manager-go/modules/replay"
	"context"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {}

func (a *App) List() []replay.Replay {
	return replay.List()
}

func (a *App) Play(fileName string) {
	replay.Play(fileName)
}

func (a *App) PlayLocal(fileName string) {
	replay.Play(fileName)
}

func (a *App) Remove(fileName string) {
	replay.Remove(fileName)
}

func (a *App) RemoveLocal(filename string) {
	replay.Remove(filename)
}

// Function to return current game version
func (a *App) GetGameVersion() string {
	return game.GetGameVersion()
}
