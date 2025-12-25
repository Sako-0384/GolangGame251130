package game

import "github.com/sorucoder/tic80"

// プレイヤー。右に掘り進みながら縦横に動く
type Player struct {
	line     *Line
	velocity Vector2d
	position Vector2d
	animTime float32

	// Visual Effects
	hurtTimer   float32
	lastHolePos Vector2d
}

func NewPlayer(line *Line) *Player {
	return &Player{
		line:        line,
		velocity:    Vector2d{0, 0},
		position:    Vector2d{120, line.GetY()}, // 初期X座標を120に変更
		lastHolePos: Vector2d{120, line.GetY()},
	}
}

func (p *Player) Update(dt float32) {
	dx := p.line.game.Speed() * dt
	dy := float32(0.0)

	p.velocity = Vector2d{dx, dy}
	p.position = p.position.Add(p.velocity)

	// Y座標をLineのレーンに同期（スムーズな移動）
	targetY := p.line.GetY()
	const moveSpeedY = 4.0 // 1フレームあたりの移動ピクセル数

	if p.position.Y < targetY {
		p.position.Y += moveSpeedY
		if p.position.Y > targetY {
			p.position.Y = targetY
		}
	} else if p.position.Y > targetY {
		p.position.Y -= moveSpeedY
		if p.position.Y < targetY {
			p.position.Y = targetY
		}
	}

	// 穴掘り処理: エフェクト生成（背景の穴）
	if p.position.Sub(p.lastHolePos).LengthSquared() > 64.0 { // 8px以上移動したら
		spawnPos := p.position.Add(Vector2d{8, 8})
		p.line.game.AddBackgroundEffect(NewHoleEffect(spawnPos.X, spawnPos.Y))
		p.lastHolePos = p.position
	}

	p.animTime += dt
	if p.animTime >= 0.2 {
		p.animTime -= 0.2
	}

	// ダメージ演出タイマー更新
	if p.hurtTimer > 0 {
		p.hurtTimer -= dt
	}
}

func (p *Player) getAnimFrame() int {
	// 下のプレイヤー（lineIndex = 1）は異なるスプライトを使用
	if p.line.lineIndex == 1 {
		switch {
		case p.hurtTimer > 0:
			return 292
		case p.animTime < 0.1:
			return 288
		case p.animTime < 0.2:
			return 290
		default:
			return 288
		}
	}

	// 上のプレイヤー（lineIndex = 0）は従来のスプライト
	switch {
	case p.hurtTimer > 0:
		return 260
	case p.animTime < 0.1:
		return 256
	case p.animTime < 0.2:
		return 258
	default:
		return 256
	}
}

func (p *Player) Draw(camera *Camera) {
	// Wiggle Effect (ダメージ時)
	drawPos := p.position
	if p.hurtTimer > 0 {
		magnitude := p.hurtTimer * 8.0
		offsetX := (RandomFloat32() - 0.5) * magnitude
		offsetY := (RandomFloat32() - 0.5) * magnitude

		drawPos = drawPos.Add(Vector2d{offsetX, offsetY})
	}

	// カメラのメソッドを使ってワールド座標をスクリーン座標に変換
	screenPos := camera.WorldToScreen(drawPos)

	// スプライト描画 (Roundを使って座標丸め)
	tic80.Spr(p.getAnimFrame(), Round(screenPos.X), Round(screenPos.Y), tic80.NewSpriteOptions().AddTransparentColor(14).SetScale(1).SetSize(2, 2))

	// ツルハシ描画
	if p.line.game.HasPickaxe(p.line.lineIndex) {
		// プレイヤーのアニメーションに合わせて上下させる
		sprite := 268
		if p.animTime >= 0.1 {
			sprite = 269
		}
		// スプライト268（16x8）を描画。プレイヤーの右側に配置
		tic80.Spr(sprite, Round(screenPos.X)+16, Round(screenPos.Y), tic80.NewSpriteOptions().AddTransparentColor(0).SetScale(1).SetSize(1, 2))
	}
}

// 衝突判定用の矩形情報を取得
func (p *Player) GetBounds() (pos Vector2d, width, height int) {
	const playerWidth = 16
	const playerHeight = 16
	return p.position, playerWidth, playerHeight
}
