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
		t := Clamp(float64(time*2.0)*PI-float64(i)*0.2, 0, PI)
		offsetY := -fastSin(t) * 10

		charStr := string(char)
		curX := x + i*width
		curY := y + int(offsetY)

		DrawOutlinedText(charStr, curX, curY, color, outlineColor)
	}
}

// 4x4 Bayer Matrix
var bayerMatrix = [16]float32{
	0.0 / 16.0, 8.0 / 16.0, 2.0 / 16.0, 10.0 / 16.0,
	12.0 / 16.0, 4.0 / 16.0, 14.0 / 16.0, 6.0 / 16.0,
	3.0 / 16.0, 11.0 / 16.0, 1.0 / 16.0, 9.0 / 16.0,
	15.0 / 16.0, 7.0 / 16.0, 13.0 / 16.0, 5.0 / 16.0,
}

// DrawDitheredBlack はディザリングのかかった黒を描画する
func DrawDitheredBlack(alpha float32) {
	if alpha <= 0.0 {
		return
	}
	if alpha >= 1.0 {
		tic80.Cls(0)
		return
	}

	for y := 0; y < 136; y++ {
		for x := 0; x < 240; x++ {
			threshold := bayerMatrix[(x%4)+(y%4)*4]

			if alpha > threshold {
				tic80.Pix(x, y, 0)
			}
		}
	}
}
