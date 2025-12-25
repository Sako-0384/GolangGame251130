package game

import (
	"github.com/sorucoder/tic80"
)

type GeneratorFactory func() LevelGenerator

type TitleScene struct {
	sceneManager *SceneManager
	ticks        int
	genFactory   GeneratorFactory
}

func NewTitleScene(sm *SceneManager, genFactory GeneratorFactory) *TitleScene {
	return &TitleScene{
		sceneManager: sm,
		ticks:        0,
		genFactory:   genFactory,
	}
}

func (s *TitleScene) OnEnter() {
	// BGM 0 (Title) Loop
	tic80.Music(tic80.NewMusicOptions().SetTrack(0))
}

func (s *TitleScene) Update(dt float32) {
	s.ticks++

	// Zボタン (Aボタン) でゲーム開始
	if tic80.Btnp(tic80.BUTTON_A, 60000, 60000) {
		newGame := NewGame(s.genFactory)

		newGame.SetSceneManager(s.sceneManager)
		s.sceneManager.ChangeScene(newGame)
	}
}

func (s *TitleScene) Draw() {
	tic80.Cls(0)

	// タイトルロゴ
	// "Gopher the Channel Miner" (24 chars * 6 = 144px width). (240-144)/2 = 48
	tic80.Print("Gopher the Channel Miner", 48, 30, tic80.NewPrintOptions().SetColor(3).SetScale(1))

	// 点滅する "PRESS B TO START"
	if (s.ticks/30)%2 == 0 {
		tic80.Print("PRESS A TO START", 75, 60, tic80.NewPrintOptions().SetColor(12))
	}

	// 操作説明など
	tic80.Print("A: MOVE UPPER PLAYER", 60, 85, tic80.NewPrintOptions().SetColor(11))
	tic80.Print("B: MOVE LOWER PLAYER", 60, 95, tic80.NewPrintOptions().SetColor(9))
	tic80.Print("X: SWAP PICKAXE", 60, 105, tic80.NewPrintOptions().SetColor(4))

	// Gopher Copyright
	tic80.Print("The Go gopher was designed", 45, 120, tic80.NewPrintOptions().SetColor(14))
	tic80.Print("by Renee French", 75, 128, tic80.NewPrintOptions().SetColor(14))
}
