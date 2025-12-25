package main

import (
	"github.com/sorucoder/tic80"

	"GolangGame251130/internal/game"
	"GolangGame251130/internal/game/generators"
)

var sm *game.SceneManager

//go:export BOOT
func BOOT() {
	tic80.Initialize()

	// reset rng seed
	ts := tic80.Tstamp()
	game.SetRandomSeed(ts)
	sm = game.NewSceneManager()

	genFactory := func() game.LevelGenerator {
		return generators.NewPathGenerator()
	}

	sm.ChangeScene(game.NewTitleScene(sm, genFactory))
}

//go:export TIC
func TIC() {
	sm.Update(1.0 / 60)
	sm.Draw()
}
