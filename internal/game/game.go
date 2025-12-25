package game

import (
	"github.com/sorucoder/tic80"
)

type Updatable interface {
	Update(dt float32)
}

type Drawable interface {
	Draw(camera *Camera)
}

type Collidable interface {
	OnCollide(collidable Collidable)
}

type Game struct {
	score         float32 // スコア（時間経過で増加）
	speed         float32
	lines         []*Line
	camera        Camera
	spawner       LevelGenerator
	genFactory    GeneratorFactory // タイトル画面に戻るために必要
	pickaxeOwner  int              // ツルハシの所持者 (0=プレイヤー1, 1=プレイヤー2)
	energy        float32          // エネルギー（ライフ）
	gameOver      bool             // ゲームオーバーフラグ
	goalDistance  float32          // ゴールまでの距離
	totalDistance float32          // 実際の総移動距離
	level         int              // 現在のレベル（周回数 + 1）
	effects          *EffectManager
	bgEffects        *EffectManager
	sceneManager     *SceneManager
	gameOverTimer    float32 // ゲームオーバー後の経過時間
	canReturnToTitle bool    // タイトルに戻れるかどうか
}

func NewGame(genFactory GeneratorFactory) *Game {
	g := &Game{
		score:         0,
		speed:         64, // スピードを32から64に増加
		lines:         []*Line{},
		camera:        Camera{Position: Vector2d{0, 0}, Scale: 1.0},
		spawner:       genFactory(),
		genFactory:    genFactory,
		pickaxeOwner:  0,   // 初期はプレイヤー1がツルハシを所持
		energy:        100, // 初期エネルギー
		gameOver:      false,
		goalDistance:  3000, // ゴール地点 (3000ピクセル)
		totalDistance:    0,
		level:            1,
		effects:          NewEffectManager(),
		bgEffects:        NewEffectManager(),
		gameOverTimer:    0,
		canReturnToTitle: false,
	}

	// 上下2つのラインを作成
	g.lines = append(g.lines, NewLine(g, 0)) // 上ライン
	g.lines = append(g.lines, NewLine(g, 1)) // 下ライン

	return g
}

func (g *Game) OnEnter() {
	// BGM 1 (Game) Loop
	tic80.Music(tic80.NewMusicOptions().SetTrack(1))
}

func (g *Game) SetGameOver() {
	if g.gameOver {
		return
	}
	g.gameOver = true
	// Stop music
	tic80.Music(tic80.NewMusicOptions().SetTrack(-1))
}

func (g *Game) Speed() float32 {
	return g.speed
}

func (g *Game) GetLevel() int {
	return g.level
}

func (g *Game) GetLines() []*Line {
	return g.lines
}

func (g *Game) GetCameraX() float32 {
	return g.camera.Position.X
}

func (g *Game) Update(dt float32) {
	// ゲームオーバーまたはクリア時は更新しない
	if g.gameOver {
		g.gameOverTimer += dt

		// 約1.5秒後にタイトルに戻れるようにする
		if g.gameOverTimer > 1.5 {
			g.canReturnToTitle = true
		}

		// ゲームオーバー時にボタン入力でタイトルへ
		if g.canReturnToTitle {
			if tic80.Btnp(tic80.BUTTON_A, 60000, 60000) || tic80.Btnp(tic80.BUTTON_B, 60000, 60000) {
				if g.sceneManager != nil {
					g.sceneManager.ChangeScene(NewTitleScene(g.sceneManager, g.genFactory))
				}
			}
		}
		return
	}

	// 総移動距離の更新
	g.totalDistance += g.speed * dt

	// ゴール判定 (totalDistanceを使用)
	// 無限ループ機能: ゴールに到達したら距離をリセットして続行
	if g.totalDistance >= g.goalDistance {
		g.totalDistance -= g.goalDistance
		g.level++

		// レベルアップ処理
		// スピード上昇: レベルごとに +8
		g.speed = 64.0 + float32(g.level-1)*8.0

		// エフェクト表示 (画面中央付近に)
		centerX := g.camera.Position.X
		g.AddEffect(NewPoppingTextEffect("LEVEL "+intToString(g.level), centerX, 60, 12))

		// SFX: Level Up (24)
		tic80.Sfx(tic80.NewSoundEffectOptions().SetId(24).SetNote(52))

		// レベルアップボーナススコア
		g.score += 1000
	}

	// スコアとエネルギーの更新
	g.score += dt * 10.0 // 1秒あたり10ポイント
	g.energy -= dt * 5.0 // 1秒あたり5減少

	if g.energy <= 0 {
		g.energy = 0
		g.SetGameOver()
		return
	}

	// ボタン入力処理
	if tic80.Btnp(tic80.BUTTON_A, 60000, 60000) && len(g.lines) > 0 {
		g.lines[0].ToggleLane()
	}
	if tic80.Btnp(tic80.BUTTON_B, 60000, 60000) && len(g.lines) > 1 {
		g.lines[1].ToggleLane()
	}

	// ツルハシ受け渡しボタン
	if tic80.Btnp(tic80.BUTTON_X, 60000, 60000) {
		oldOwner := g.pickaxeOwner
		g.pickaxeOwner = 1 - g.pickaxeOwner // 0→1, 1→0 に切り替え

		// 受け渡しエフェクト発生
		// 両プレイヤーの位置を取得
		if len(g.lines) >= 2 {
			p1 := g.lines[oldOwner].player.position
			p2 := g.lines[g.pickaxeOwner].player.position

			// ツルハシの位置（プレイヤー右側）に合わせる
			offset := Vector2d{20, 8}
			p1 = p1.Add(offset)
			p2 = p2.Add(offset)

			g.AddEffect(NewTransferEffect(p1, p2, g.speed))

			// SFX: Pickaxe Transfer (14) Note: 57
			tic80.Sfx(tic80.NewSoundEffectOptions().SetId(14).SetNote(57))
		}
	}

	// アイテムスポーン管理（Game全体で管理）
	g.spawnItems()

	for i := range g.lines {
		g.lines[i].Update(dt)
	}

	// エフェクト更新
	g.bgEffects.Update(dt)
	g.effects.Update(dt)

	// カメラを最前のプレイヤーに追従させる
	if len(g.lines) > 0 {
		var frontmostX float32 = -999999
		for i := range g.lines {
			if g.lines[i].player != nil {
				if g.lines[i].player.position.X > frontmostX {
					frontmostX = g.lines[i].player.position.X
				}
			}
		}
		offset := float32(60) // プレイヤーの右側60ピクセルをカメラ中心に
		g.camera.Position.X = frontmostX + offset
	}

	// 座標リセット（float丸め誤差対策）
	// カメラX座標が1000を超えたら、全ての座標を平行移動
	if g.camera.Position.X > 1000 {
		resetOffset := g.camera.Position.X - 100 // カメラを100付近に戻す

		// カメラをリセット
		g.camera.Position.X -= resetOffset

		// Spawnerに通知
		g.spawner.OnCoordinateReset(resetOffset)

		// 全プレイヤーの座標をリセット
		for i := range g.lines {
			if g.lines[i].player != nil {
				g.lines[i].player.position.X -= resetOffset
			}

			// 全アイテムの座標をリセット
			for j := range g.lines[i].items {
				pos := g.lines[i].items[j].GetPosition()
				pos.X -= resetOffset
				g.lines[i].items[j].SetPosition(pos)
			}
		}

		// エフェクトの座標をリセット
		g.effects.OnCoordinateReset(resetOffset)
		g.bgEffects.OnCoordinateReset(resetOffset)
	}
}

// SetSceneManager sets the scene manager for the game
func (g *Game) SetSceneManager(sm *SceneManager) {
	g.sceneManager = sm
}

// HasPickaxe は指定したラインのプレイヤーがツルハシを所持しているかを返す
func (g *Game) HasPickaxe(lineIndex int) bool {
	return g.pickaxeOwner == lineIndex
}

// AddEnergy はエネルギーを追加する（Foodの取得時など）
func (g *Game) AddEnergy(amount float32) {
	g.energy += amount
	if g.energy > 300 {
		g.energy = 300
	}
	if g.energy <= 0 {
		g.energy = 0
		g.SetGameOver()
	}
}

// アイテムスポーン管理（spawnerに委譲）
func (g *Game) spawnItems() {
	if g.spawner.ShouldSpawn(g) {
		g.spawner.SpawnItem(g)
	}
}

func (g *Game) Draw() {
	tic80.Cls(13)

	// 背景マップ描画
	// カメラ位置に合わせてマップを表示範囲分だけ描画
	// WorldToScreenは画面中央(120)にCameraPosが来るようになっているため、
	// マップ描画の開始位置（左端）は CameraPos - 120 となる。
	startWorldX := g.camera.Position.X - 120
	// Roundを使って整数座標に丸める（スプライトと合わせるため）
	startWorldX_Int := Round(startWorldX)

	tileX := startWorldX_Int / 8
	offsetX := startWorldX_Int % 8
	// 負の剰余の補正
	if offsetX < 0 {
		offsetX += 8
		tileX -= 1 // 負の方向に1タイル分ずらす
	}

	// マップの幅は240タイル。無限スクロールのためにラップアラウンド処理を行う
	mapX := tileX % 240
	if mapX < 0 {
		mapX += 240
	}
	tilesToDraw := 31 // 画面幅(240px) / 8px = 30タイル + バッファ1

	if mapX+tilesToDraw <= 240 {
		// 通常描画（ラップなし）
		tic80.Map(tic80.NewMapOptions().SetOffset(mapX, 0).SetSize(tilesToDraw, 18).SetPosition(-offsetX, 0))
	} else {
		// ラップアラウンド描画（右端まで描画し、残りを左端から描画）
		firstChunkWidth := 240 - mapX
		secondChunkWidth := tilesToDraw - firstChunkWidth

		// 1. 右端部分
		tic80.Map(tic80.NewMapOptions().SetOffset(mapX, 0).SetSize(firstChunkWidth, 18).SetPosition(-offsetX, 0))

		// 2. 左端部分（折り返し）
		// 描画位置は -offsetX + (firstChunkWidth * 8)
		tic80.Map(tic80.NewMapOptions().SetOffset(0, 0).SetSize(secondChunkWidth, 18).SetPosition(-offsetX+(firstChunkWidth*8), 0))
	}

	// 背景エフェクト描画
	g.bgEffects.Draw(&g.camera)

	for i := range g.lines {
		g.lines[i].Draw(&g.camera)
	}

	// エフェクト描画
	g.effects.Draw(&g.camera)

	// UI描画
	g.DrawUI()
}

// AddEffect はエフェクトを追加する
func (g *Game) AddEffect(e Effect) {
	g.effects.Add(e)
}

// AddBackgroundEffect は背景エフェクトを追加する
func (g *Game) AddBackgroundEffect(e Effect) {
	g.bgEffects.Add(e)
}
