package game

import (
	"github.com/sorucoder/tic80"
)

type GeneratorFactory func() LevelGenerator

type TitleScene struct {
	sceneManager    *SceneManager
	ticks           int
	genFactory      GeneratorFactory
	isTransitioning bool
	transitionTimer float32
}

func NewTitleScene(sm *SceneManager, genFactory GeneratorFactory) *TitleScene {
	return &TitleScene{
		sceneManager:    sm,
		ticks:           0,
		genFactory:      genFactory,
		isTransitioning: false,
		transitionTimer: 0,
	}
}

func (s *TitleScene) OnEnter() {
	// BGM 0 (Title) Loop
	tic80.Music(tic80.NewMusicOptions().SetTrack(0))
}

func (s *TitleScene) Update(dt float32) {
	s.ticks++

	if s.isTransitioning {
		s.transitionTimer += dt
		if s.transitionTimer > 1.0 {
			newGame := NewGame(s.genFactory)
			newGame.SetSceneManager(s.sceneManager)
			s.sceneManager.ChangeScene(newGame)
		}
		return
	}

	// Zボタン (Aボタン) でゲーム開始
	if tic80.Btnp(tic80.BUTTON_A, 60000, 60000) {
		s.isTransitioning = true
		s.transitionTimer = 0

		// BGM停止
		tic80.Music(tic80.NewMusicOptions().SetTrack(-1))
		tic80.Sfx(tic80.NewSoundEffectOptions().SetId(8).SetNote(64))
	}
}

func (s *TitleScene) Draw() {
	tic80.Cls(0)

	// map描画(BG)
	tic80.Map(tic80.NewMapOptions().SetOffset(0, 17).SetSize(31, 12).SetPosition(0, 0))

	// Gopher Sprites
	leftSprite := 256
	rightSprite := 288
	offsetY := 0.0

	// 決定時の演出
	if s.isTransitioning {
		leftSprite = 258
		rightSprite = 290

		// ジャンプ演出 (0.0 -> 0.5秒でもぐる)
		if s.transitionTimer < 1.0 {
			t := s.transitionTimer * 2.0 // 0.0 -> 1.0
			jumpHeight := 20.0
			offsetY = float64((4.0*t*(1.0-t)*float32(jumpHeight) - t*16.0))
		}
	}

	// Left Gopher (Offset X: 48 - 20 = 28, Y: 26)
	tic80.Spr(leftSprite, 28, 80-int(offsetY), tic80.NewSpriteOptions().AddTransparentColor(14).SetScale(1).SetSize(2, 2))

	// Right Gopher (Offset X: 48 + 144 + 4 = 196, Y: 26)
	// Flip horizontally
	tic80.Spr(rightSprite, 196, 80-int(offsetY), tic80.NewSpriteOptions().AddTransparentColor(14).SetScale(1).SetSize(2, 2).FlipHorizontally())

	// map描画(FG)
	tic80.Map(tic80.NewMapOptions().SetOffset(0, 29).SetSize(31, 5).SetPosition(0, 96))

	// タイトルロゴ
	DrawOutlinedText("Gopher the Channel Miner", 56, 30, 3, 15)

	// 点滅する "PRESS A TO START"
	if (s.ticks/30)%2 == 0 {
		DrawOutlinedText("PRESS A TO START", 80, 40, 12, 15)
	}

	// 操作説明など
	tic80.Print("A: MOVE UPPER PLAYER", 68, 65, tic80.NewPrintOptions().SetColor(11))
	tic80.Print("B: MOVE LOWER PLAYER", 68, 75, tic80.NewPrintOptions().SetColor(9))
	tic80.Print("X: SWAP PICKAXE", 68, 85, tic80.NewPrintOptions().SetColor(4))

	// Gopher Copyright
	tic80.Print("The Go gopher was designed", 48, 110, tic80.NewPrintOptions().SetColor(13))
	tic80.Print("by Renee French", 84, 120, tic80.NewPrintOptions().SetColor(13))

	// Transition Effect
	if s.isTransitioning {
		alpha := s.transitionTimer / 1.0 // 1.0秒で完了
		DrawDitheredBlack(alpha)
	}
}
