package game

import "github.com/sorucoder/tic80"

// 列。
type Line struct {
	game        *Game
	player      *Player
	items       []Item
	lineIndex   int // ライン番号（0=上、1=下）
	currentLane int // 現在のレーン（0または1）
}

func NewLine(game *Game, lineIndex int) *Line {
	l := &Line{
		game:        game,
		lineIndex:   lineIndex,
		currentLane: 0, // 最初はレーン0から開始
	}

	l.player = NewPlayer(l)

	return l
}

func (l *Line) Update(dt float32) {
	l.player.Update(dt)

	// アイテム更新と削除（生成はGameで管理）
	// 衝突判定も同時に行う
	activeItems := l.items[:0]
	playerPos, playerWidth, playerHeight := l.player.GetBounds()

	for i := range l.items {
		l.items[i].Update(dt)

		// 衝突判定
		if l.items[i].CollidesWith(playerPos, playerWidth, playerHeight) {
			// デバッグ: 衝突情報を表示
			itemType := "Food"
			if l.items[i].IsObstacle() {
				itemType = "Rock"
			}
			hasPickaxe := l.game.HasPickaxe(l.lineIndex)

			// 衝突した場合の処理
			if l.items[i].IsObstacle() {
				// HardRock判定
				_, isHardRock := l.items[i].(*HardRock)
				
				// Rock/GoldRock: ツルハシ所持状態をチェック
				if !isHardRock && hasPickaxe {
					// ツルハシ所持: Rockを破壊（削除）
					if _, ok := l.items[i].(*GoldRock); ok {
						// GoldRock破壊ボーナス
						l.game.score += 500
						l.game.AddEffect(NewFloatingTextEffect("+500", playerPos.X, playerPos.Y-10, 12)) // 12=Light Blue
						// SFX: GoldRock (11)
						tic80.Sfx(tic80.NewSoundEffectOptions().SetId(11).SetNote(64))
						// パーティクルを散らす
						for k := 0; k < 10; k++ {
							l.game.AddEffect(NewParticleEffect(playerPos.X+8, playerPos.Y+8, 14)) // 14=Yellow
						}
						tic80.Trace("P"+intToString(l.lineIndex+1)+": GoldRock DESTROYED! (+500 Score)", tic80.NewTraceOptions())
					} else {
						// Normal Rock
						tic80.Trace("P"+intToString(l.lineIndex+1)+": "+itemType+" HIT! (Pickaxe: DESTROY)", tic80.NewTraceOptions())
						// SFX: Normal Rock Destroy (12)
						tic80.Sfx(tic80.NewSoundEffectOptions().SetId(12).SetNote(64))
						// パーティクルを散らす (グレー: 13)
						for k := 0; k < 10; k++ {
							l.game.AddEffect(NewParticleEffect(playerPos.X+8, playerPos.Y+8, 13)) 
						}
					}
					continue
				} else {
					// ツルハシ非所持 または HardRock: エネルギー減少
					tic80.Trace("P"+intToString(l.lineIndex+1)+": "+itemType+" HIT! (DAMAGE)", tic80.NewTraceOptions())
					l.game.AddEnergy(-30) // エネルギー減少
					l.game.AddEffect(NewFloatingTextEffect("-30", playerPos.X, playerPos.Y-10, 6)) // 6=Red
					// SFX: Rock Damage (10)
					tic80.Sfx(tic80.NewSoundEffectOptions().SetId(10).SetNote(40))
					l.player.hurtTimer = 0.5 // 0.5秒間揺れる
					continue
				}
			} else {
				// Food: 取得して削除
				tic80.Trace("P" + intToString(l.lineIndex+1) + ": " + itemType + " GET!", tic80.NewTraceOptions())
				l.game.AddEnergy(20) // エネルギー回復
				l.game.AddEffect(NewFloatingTextEffect("+20", playerPos.X, playerPos.Y-10, 11)) // 11=Cyan/Greenish
				// SFX: Food (08)
				tic80.Sfx(tic80.NewSoundEffectOptions().SetId(8).SetNote(64))
				continue
			}
		}

		// 期限切れチェック
		if !l.items[i].IsExpired() {
			activeItems = append(activeItems, l.items[i])
		}
	}
	l.items = activeItems
}

// レーンを切り替える
func (l *Line) ToggleLane() {
	if l.currentLane == 0 {
		l.currentLane = 1
	} else {
		l.currentLane = 0
	}
	// SFX: Movement (13) Note: 33
	tic80.Sfx(tic80.NewSoundEffectOptions().SetId(13).SetNote(33))
}

// 現在のY座標を計算（ラインとレーンに基づく）
func (l *Line) GetY() float32 {
	return l.GetLaneY(l.currentLane)
}

// 指定したレーンのY座標を取得
func (l *Line) GetLaneY(lane int) float32 {
	const lineSpacing = 56.0 // ライン間のスペース
	const laneSpacing = 16.0 // レーン間のスペース
	baseY := float32(24.0)   // 最初のラインのベースY座標

	lineY := baseY + float32(l.lineIndex)*lineSpacing
	laneOffset := float32(lane) * laneSpacing

	return lineY + laneOffset
}


func (l *Line) Draw(camera *Camera) {
	// TODO: draw bg

	// アイテム描画
	for i := range l.items {
		l.items[i].Draw(camera)
	}

	l.player.Draw(camera)
}

func (l *Line) AddItem(item Item) {
	l.items = append(l.items, item)
}
