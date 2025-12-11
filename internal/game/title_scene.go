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

func (s *TitleScene) Update(dt float32) {
	s.ticks++

	// Zボタン (Aボタン) でゲーム開始
	if tic80.Btnp(tic80.BUTTON_A, 60000, 60000) {
		// ゲームシーンへ遷移
		// GeneratorFactoryを使って新しいGeneratorを生成
		// GeneratorFactoryを直接渡す(Game内で生成される)
		newGame := NewGame(s.genFactory)
		
		newGame.SetSceneManager(s.sceneManager)
		s.sceneManager.ChangeScene(newGame)
	}
}

func (s *TitleScene) Draw() {
	tic80.Cls(0)

	// タイトルロゴ的な表示
	tic80.Print("GOLANG GAME 251130", 60, 50, tic80.NewPrintOptions().SetColor(12).SetScale(1))

	// 点滅する "PRESS Z TO START"
	if (s.ticks/30)%2 == 0 {
		tic80.Print("PRESS Z TO START", 75, 80, tic80.NewPrintOptions().SetColor(15))
	}

	// 操作説明など
	tic80.Print("B: MOVE UPPER PLAYER", 10, 110, tic80.NewPrintOptions().SetColor(13))
	tic80.Print("A: MOVE LOWER PLAYER", 10, 120, tic80.NewPrintOptions().SetColor(13))
	tic80.Print("X: SWAP PICKAXE", 10, 130, tic80.NewPrintOptions().SetColor(13))
}
