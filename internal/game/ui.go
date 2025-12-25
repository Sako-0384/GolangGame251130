package game

import (
	"github.com/sorucoder/tic80"
)

// DrawUI はゲームのUIを描画する
func (g *Game) DrawUI() {
	// UIのベースY座標（上下ラインの間）
	baseY := 65

	// 3カラム構成
	// Total Width: 240
	// Col 1 (Score): 0-80
	// Col 2 (Progress): 80-160
	// Col 3 (Energy): 160-240

	// --- Column 1: Score ---
	scoreText := "SC:" + intToString(int(g.score))
	
	// ゲームオーバーアニメーション中（1.0秒以降）はデフォルトの表示を隠す
	shouldDrawDefaultScore := true
	if g.gameOver && g.gameOverTimer > 1.0 {
		shouldDrawDefaultScore = false
	}

	if shouldDrawDefaultScore {
		DrawOutlinedText(scoreText, 2, baseY, 4, 14)
	}

	// --- Column 2: Progress Bar ---
	progressWidth := 70
	progressX := 85
	progressHeight := 6

	// 背景
	tic80.Rect(progressX, baseY, progressWidth, progressHeight, 0)
	// 枠線
	tic80.Rectb(progressX-1, baseY-1, progressWidth+2, progressHeight+2, 12)

	// 進捗
	progress := g.totalDistance / g.goalDistance
	if progress > 1.0 {
		progress = 1.0
	}
	if progress < 0 {
		progress = 0
	}
	fillWidth := int(float32(progressWidth) * progress)
	tic80.Rect(progressX, baseY, fillWidth, progressHeight, 11)

	// ゴールアイコン
	tic80.Print("G", progressX+progressWidth+2, baseY, tic80.NewPrintOptions().SetColor(12).TogglePage())

	// --- Column 3: Energy Bar ---
	energyWidth := 70
	energyX := 165
	energyHeight := 6

	// 背景
	tic80.Rect(energyX, baseY, energyWidth, energyHeight, 0)
	// 枠線
	tic80.Rectb(energyX-1, baseY-1, energyWidth+2, energyHeight+2, 12)

	// エネルギーバー
	if g.energy <= 100 {
		// 0-100: 緑
		tic80.Rect(energyX, baseY, int(float32(energyWidth)*(g.energy/100.0)), energyHeight, 5)
	} else if g.energy <= 200 {
		// 0-100: 緑
		tic80.Rect(energyX, baseY, energyWidth, energyHeight, 5)
		// 100-200: 黄
		overWidth := int(float32(energyWidth) * ((g.energy - 100.0) / 100.0))
		tic80.Rect(energyX, baseY, overWidth, energyHeight, 4)
	} else {
		// 100-200: 黄
		tic80.Rect(energyX, baseY, energyWidth, energyHeight, 4)
		// 200-300: 橙
		overWidth := int(float32(energyWidth) * ((g.energy - 200.0) / 100.0))
		tic80.Rect(energyX, baseY, overWidth, energyHeight, 3)
	}

	// 数値
	tic80.Print(intToString(int(g.energy)), energyX+2, baseY+1, tic80.NewPrintOptions().SetColor(0).TogglePage())

	// --- Overlays ---

	// ゲームオーバー表示とアニメーション
	if g.gameOver {
		// 1. GAME OVER テキスト (0.3秒後に表示)
		if g.gameOverTimer > 0.3 {
			text := "GAME OVER"
			DrawOutlinedText(text, 96, 50, 6, 12)
		}

		// 2. スコア位置のアニメーション
		targetX := float32(240/2 - 20)
		targetY := float32(66)

		currentX := float32(2)
		currentY := float32(baseY)

		if g.gameOverTimer > 1.0 {
			t := (g.gameOverTimer - 1.0) / 0.5
			t = float32(Clamp(float64(t), 0.0, 1.0))
			easeT := EaseInOutCubic(t)

			drawX := Lerp(currentX, targetX, easeT)
			drawY := Lerp(currentY, targetY, easeT)

			DrawOutlinedText(scoreText, int(drawX), int(drawY), 4, 14)

			if g.canReturnToTitle {
				color := 12
				if (int(g.gameOverTimer*2) % 2) == 0 {
					color = 0
				}
				prompt := "PRESS BUTTON"
				promptWidth := tic80.Print(prompt, 0, -10, tic80.NewPrintOptions())
				tic80.Print(prompt, (240-promptWidth)/2, 80, tic80.NewPrintOptions().SetColor(color))
			}
		}
	}
}
