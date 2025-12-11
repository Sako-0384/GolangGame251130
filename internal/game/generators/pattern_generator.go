package generators

import (
	"GolangGame251130/internal/game"
)

// PatternItem はパターン内の個々のアイテム定義
type PatternItem struct {
	Type    int     // 0: Removed (WideRock), 1: Rock, 2: Food, 3: GoldRock
	LineIdx int     // 0: Upper, 1: Lower
	Lane    int     // 0, 1
	OffsetX float32 // パターン開始位置からのオフセット
}

// Pattern はアイテムの配置パターン
type Pattern struct {
	Items []PatternItem
	Width float32 // パターンの全長（次の生成までの間隔）
}

// PatternGenerator はパターンベースのレベル生成
type PatternGenerator struct {
	patterns   []Pattern
	lastSpawnX float32 // 最後にスポーン判定を行った位置（重複防止）
	nextSpawnX float32 // 次にスポーンすべき位置
	lastFoodX  float32 // 最後にFoodを生成した位置
}

func NewPatternGenerator() *PatternGenerator {
	s := &PatternGenerator{
		lastSpawnX: 0,
		nextSpawnX: 300, // 初期スポーン位置
		lastFoodX:  0,
	}
	s.initPatterns()
	return s
}

func (s *PatternGenerator) initPatterns() {
	// パターン定義
	s.patterns = []Pattern{
		// 1. 基本的な岩（シングル） - 高密度化
		{
			Items: []PatternItem{
				{Type: 1, LineIdx: 0, Lane: 0, OffsetX: 0},
			},
			Width: 40,
		},
		{
			Items: []PatternItem{
				{Type: 1, LineIdx: 1, Lane: 1, OffsetX: 0},
			},
			Width: 40,
		},
		// 2. 2列の岩（同時） - 高密度化
		{
			Items: []PatternItem{
				{Type: 1, LineIdx: 0, Lane: 0, OffsetX: 0},
				{Type: 1, LineIdx: 1, Lane: 0, OffsetX: 0},
			},
			Width: 60,
		},
		// 3. 交互の岩（ジグザグ） - 間隔短縮
		{
			Items: []PatternItem{
				{Type: 1, LineIdx: 0, Lane: 1, OffsetX: 0},
				{Type: 1, LineIdx: 1, Lane: 0, OffsetX: 25},
				{Type: 1, LineIdx: 0, Lane: 1, OffsetX: 50},
			},
			Width: 80,
		},
		// 4. WideRock (Legacy replaced with 2 Rocks)
		{
			Items: []PatternItem{
				{Type: 1, LineIdx: 0, Lane: 0, OffsetX: 0},
				{Type: 1, LineIdx: 0, Lane: 1, OffsetX: 0},
			},
			Width: 60,
		},
		// 5. GoldRockチャンス（岩の後ろにGoldRock） - 密着させる
		{
			Items: []PatternItem{
				{Type: 1, LineIdx: 0, Lane: 0, OffsetX: 0},
				{Type: 3, LineIdx: 0, Lane: 0, OffsetX: 18}, // 岩のすぐ後ろ(16px+2px)
			},
			Width: 50,
		},
		// 6. 連続岩 - 間隔短縮
		{
			Items: []PatternItem{
				{Type: 1, LineIdx: 1, Lane: 0, OffsetX: 0},
				{Type: 1, LineIdx: 1, Lane: 0, OffsetX: 20},
				{Type: 1, LineIdx: 1, Lane: 0, OffsetX: 40},
			},
			Width: 70,
		},
		// 7. [NEW] GoldRock Rush (GoldRockが連続するボーナスパターン)
		{
			Items: []PatternItem{
				{Type: 3, LineIdx: 0, Lane: 0, OffsetX: 0},
				{Type: 3, LineIdx: 1, Lane: 1, OffsetX: 20},
				{Type: 3, LineIdx: 0, Lane: 0, OffsetX: 40},
			},
			Width: 80,
		},
		// 8. [NEW] Hard Wall (WideRock replaced + Rock)
		{
			Items: []PatternItem{
				{Type: 1, LineIdx: 0, Lane: 0, OffsetX: 0},   // 上の壁 (Rock 1)
				{Type: 1, LineIdx: 0, Lane: 1, OffsetX: 0},   // 上の壁 (Rock 2)
				{Type: 1, LineIdx: 1, Lane: 1, OffsetX: 20}, // 下の岩
			},
			Width: 70,
		},
	}
}

func (s *PatternGenerator) ShouldSpawn(g *game.Game) bool {
	cameraX := g.GetCameraX()
	spawnThreshold := cameraX + 300 // 画面外で生成

	// 次の生成位置に達したか
	return s.nextSpawnX < spawnThreshold
}

func (s *PatternGenerator) SpawnItem(g *game.Game) {
	// Foodの頻度チェック
	level := g.GetLevel()
	foodInterval := float32(125 + (level-1)*25)
	
	var pattern Pattern
	
	if s.nextSpawnX - s.lastFoodX > foodInterval {
		// Foodパターン (シングル・完全ランダム)
		pattern = Pattern{
			Items: []PatternItem{
				{Type: 2, LineIdx: game.Intn(2), Lane: game.Intn(2), OffsetX: 0},
			},
			Width: 60,
		}
	} else {
		// ランダムにパターン選択
		idx := game.Intn(len(s.patterns))
		pattern = s.patterns[idx]
	}

	// パターン全体のランダム反転 (GoldRock Rush以外に適用)
	isRusher := len(pattern.Items) > 0 && pattern.Items[0].Type == 3
	flipLane := !isRusher && game.Intn(2) == 0
	flipLine := !isRusher && game.Intn(2) == 0

	lines := g.GetLines()

	// パターン内のアイテムを生成
	for _, item := range pattern.Items {
		spawnX := s.nextSpawnX + item.OffsetX
		
		// パターン反転適用
		currentLineIdx := item.LineIdx
		currentLane := item.Lane

		if flipLine {
			currentLineIdx = 1 - currentLineIdx
		}
		if flipLane {
			currentLane = 1 - currentLane
		}

		// GoldRock Rush (Type 3 rush) の場合は個別にランダム化
		if isRusher && item.Type == 3 {
			currentLineIdx = game.Intn(2)
			currentLane = game.Intn(2)
		}

		// アイテム変異 (Rock -> GoldRock)
		itemType := item.Type
		if itemType == 1 {
			// 10%の確率でGoldRock化
			if game.Intn(10) == 0 {
				itemType = 3
			}
		}
		
		targetLine := lines[currentLineIdx]
		var newItem game.Item

		switch itemType {
		
		case 1: // Rock
			newItem = game.NewRock(targetLine, spawnX, currentLane)
		case 2: // Food
			newItem = game.NewFood(targetLine, spawnX, currentLane)
			if spawnX > s.lastFoodX {
				s.lastFoodX = spawnX
			}
		case 3: // GoldRock
			newItem = game.NewGoldRock(targetLine, spawnX, currentLane)
		case 4: // HardRock
			newItem = game.NewHardRock(targetLine, spawnX, currentLane)
		}

		if newItem != nil {
			targetLine.AddItem(newItem)
		}
	}

	// 次の生成位置を更新
	s.nextSpawnX += pattern.Width
	
	// ランダムな間隔を追加（バラつきを出す）
	s.nextSpawnX += float32(game.Intn(50))
}

func (s *PatternGenerator) OnCoordinateReset(offset float32) {
	s.nextSpawnX -= offset
	s.lastSpawnX -= offset
	s.lastFoodX -= offset
}
