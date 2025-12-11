package game

// LevelGenerator はレベル生成ロジックを抽象化するインターフェース
type LevelGenerator interface {
	// ShouldSpawn は新しいアイテムを生成すべきかどうかを判定
	ShouldSpawn(game *Game) bool

	// SpawnItem はアイテムを生成してLineに追加
	SpawnItem(game *Game)

	// OnCoordinateReset は座標リセット時に呼ばれる
	OnCoordinateReset(offset float32)
}
