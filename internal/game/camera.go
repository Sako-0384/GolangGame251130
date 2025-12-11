package game

const (
	ScreenWidth  = 240
	ScreenHeight = 136
)

type Camera struct {
	Position Vector2d
	Scale    float32 // ズーム用（1.0が等倍）
}

// ワールド座標をスクリーン座標に変換
func (c *Camera) WorldToScreen(worldPos Vector2d) Vector2d {
	screenCenterX := float32(ScreenWidth / 2)
	screenCenterY := float32(0.0)

	return Vector2d{
		X: (worldPos.X-c.Position.X)*c.Scale + screenCenterX,
		Y: (worldPos.Y-c.Position.Y)*c.Scale + screenCenterY, // Y座標もカメラ位置を考慮
	}
}

// ワールドX座標をスクリーンX座標に変換
func (c *Camera) WorldToScreenX(worldX float32) float32 {
	screenCenterX := float32(ScreenWidth / 2)
	return (worldX-c.Position.X)*c.Scale + screenCenterX
}

// ワールドY座標をスクリーンY座標に変換
func (c *Camera) WorldToScreenY(worldY float32) float32 {
	screenCenterY := float32(ScreenHeight / 2)
	return (worldY-c.Position.Y)*c.Scale + screenCenterY
}

// 現在のスケールを取得
func (c *Camera) GetScale() float32 {
	return c.Scale
}
