package main

import (
	"GolangGame251130/internal/game"
	"GolangGame251130/internal/game/generators"
)

var sm *game.SceneManager

//go:export BOOT
func BOOT() {
	sm = game.NewSceneManager()

	// Level Generator Factory define
	// どちらを使うかここで切り替え
	genFactory := func() game.LevelGenerator {
		// return generators.NewPatternGenerator()
		return generators.NewRuleBasedGenerator()
	}

	// TitleSceneにFactoryを渡す
	sm.ChangeScene(game.NewTitleScene(sm, genFactory))
}

//go:export TIC
func TIC() {
	sm.Update(1.0 / 60)
	sm.Draw()
}
