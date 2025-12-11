package game

import (
	"github.com/sorucoder/tic80"
)

type Item interface {
	Updatable
	Drawable

	GetPosition() Vector2d
	SetPosition(pos Vector2d)
	Width() int
	Height() int
	IsObstacle() bool
	IsExpired() bool
	CollidesWith(pos Vector2d, width, height int) bool
}

// 基本的なアイテム構造体
type BaseItem struct {
	line     *Line
	Position Vector2d // エクスポート（座標リセット用）
	width    int
	height   int
}

func (i *BaseItem) Update(dt float32) {
	// アイテムは動かない（相対的に左に流れるように見えるが、実際にはカメラが右に進む）
	// ただし、もしアイテム自体を動かす必要があるならここで処理
}

func (i *BaseItem) GetPosition() Vector2d {
	return i.Position
}

func (i *BaseItem) SetPosition(pos Vector2d) {
	i.Position = pos
}

func (i *BaseItem) Width() int {
	return i.width
}

func (i *BaseItem) Height() int {
	return i.height
}

func (i *BaseItem) IsExpired() bool {
	// カメラの左端よりさらに左に行ったら削除
	cameraX := i.line.game.GetCameraX()
	return i.Position.X < cameraX-140 // 画面幅(240)/2 + マージン
}

// AABB衝突判定
// AABB衝突判定
func (i *BaseItem) CollidesWith(pos Vector2d, width, height int) bool {
	return i.Position.X < pos.X+float32(width) &&
		i.Position.X+float32(i.width) > pos.X &&
		i.Position.Y < pos.Y+float32(height) &&
		i.Position.Y+float32(i.height) > pos.Y
}

// 岩。障害物。
type Rock struct {
	BaseItem
}

func NewRock(line *Line, x float32, lane int) *Rock {
	y := line.GetLaneY(lane)
	return &Rock{
		BaseItem: BaseItem{
			line:     line,
			Position: Vector2d{x, y},
			width:    16,
			height:   16,
		},
	}
}

func (r *Rock) IsObstacle() bool {
	return true
}

func (r *Rock) Draw(camera *Camera) {
	screenPos := camera.WorldToScreen(r.Position)
	tic80.Spr(386, Round(screenPos.X), Round(screenPos.Y), tic80.NewSpriteOptions().AddTransparentColor(2).SetScale(1).SetSize(2, 2))
}

// 食物。スコアになる。
type Food struct {
	BaseItem
}

func NewFood(line *Line, x float32, lane int) *Food {
	y := line.GetLaneY(lane)
	return &Food{
		BaseItem: BaseItem{
			line:     line,
			Position: Vector2d{x, y},
			width:    16,
			height:   16,
		},
	}
}

func (f *Food) IsObstacle() bool {
	return false
}

func (f *Food) Draw(camera *Camera) {
	screenPos := camera.WorldToScreen(f.Position)
	tic80.Spr(384, Round(screenPos.X), Round(screenPos.Y), tic80.NewSpriteOptions().AddTransparentColor(14).SetScale(1).SetSize(2, 2))
}

// 金塊岩。壊すと高得点。
type GoldRock struct {
	BaseItem
}

func NewGoldRock(line *Line, x float32, lane int) *GoldRock {
	y := line.GetLaneY(lane)
	return &GoldRock{
		BaseItem: BaseItem{
			line:     line,
			Position: Vector2d{x, y},
			width:    16,
			height:   16,
		},
	}
}

func (g *GoldRock) IsObstacle() bool {
	return true
}

func (g *GoldRock) Draw(camera *Camera) {
	screenPos := camera.WorldToScreen(g.Position)
	tic80.Spr(388, Round(screenPos.X), Round(screenPos.Y), tic80.NewSpriteOptions().AddTransparentColor(2).SetScale(1).SetSize(2, 2))
}

// 硬い岩。壊せない障害物。
type HardRock struct {
	BaseItem
}

func NewHardRock(line *Line, x float32, lane int) *HardRock {
	y := line.GetLaneY(lane)
	return &HardRock{
		BaseItem: BaseItem{
			line:     line,
			Position: Vector2d{x, y},
			width:    16,
			height:   16,
		},
	}
}

func (h *HardRock) IsObstacle() bool {
	return true
}

func (h *HardRock) Draw(camera *Camera) {
	screenPos := camera.WorldToScreen(h.Position)
	tic80.Spr(390, Round(screenPos.X), Round(screenPos.Y), tic80.NewSpriteOptions().AddTransparentColor(2).SetScale(1).SetSize(2, 2))
}
