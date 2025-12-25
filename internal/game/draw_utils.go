package game

import (
	"github.com/sorucoder/tic80"
)

// DrawOutlinedText は枠線付きテキストを描画する
func DrawOutlinedText(text string, x, y int, color, outlineColor int) {
	// 枠線（上下左右斜め）
	tic80.Print(text, x-1, y, tic80.NewPrintOptions().SetColor(outlineColor))
	tic80.Print(text, x+1, y, tic80.NewPrintOptions().SetColor(outlineColor))
	tic80.Print(text, x, y-1, tic80.NewPrintOptions().SetColor(outlineColor))
	tic80.Print(text, x, y+1, tic80.NewPrintOptions().SetColor(outlineColor))

	// 本体
	tic80.Print(text, x, y, tic80.NewPrintOptions().SetColor(color))
}

// DrawPoppingText は波打つテキストを描画する
func DrawPoppingText(text string, x, y int, color, outlineColor int, time float32) {
	width := 6 // 1文字あたりの概算幅（TIC-80のフォントサイズによる）

	for i, char := range text {
		t := Clamp(float64(time*2.0)*PI - float64(i)*0.2, 0, PI)
		offsetY := -fastSin(t)*10

		charStr := string(char)
		curX := x + i*width
		curY := y + int(offsetY)

		DrawOutlinedText(charStr, curX, curY, color, outlineColor)
	}
}
