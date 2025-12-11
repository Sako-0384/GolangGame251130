package game

import (
	"github.com/sorucoder/tic80"
)

// Effect は視覚効果のインターフェース
type Effect interface {
	Update(dt float32) bool // 戻り値: つづけるならtrue, 終了ならfalse
	Draw(camera *Camera)
	OnCoordinateReset(offsetX float32)
}

// EffectManager はエフェクトを管理する
type EffectManager struct {
	effects []Effect
}

func NewEffectManager() *EffectManager {
	return &EffectManager{
		effects: []Effect{},
	}
}

func (em *EffectManager) Add(e Effect) {
	em.effects = append(em.effects, e)
}

func (em *EffectManager) Update(dt float32) {
	activeEffects := em.effects[:0]
	for _, e := range em.effects {
		if e.Update(dt) {
			activeEffects = append(activeEffects, e)
		}
	}
	em.effects = activeEffects
}

func (em *EffectManager) Draw(camera *Camera) {
	for _, e := range em.effects {
		e.Draw(camera)
	}
}

func (em *EffectManager) OnCoordinateReset(offsetX float32) {
	for _, e := range em.effects {
		e.OnCoordinateReset(offsetX)
	}
}

// FloatingTextEffect は上に浮かび上がるテキスト
type FloatingTextEffect struct {
	text     string
	position Vector2d
	color    int
	lifeTime float32
	maxLife  float32
}

func NewFloatingTextEffect(text string, x, y float32, color int) *FloatingTextEffect {
	return &FloatingTextEffect{
		text:     text,
		position: Vector2d{x, y},
		color:    color,
		lifeTime: 0,
		maxLife:  1.5, // 1.5秒で消える
	}
}

func (e *FloatingTextEffect) Update(dt float32) bool {
	e.lifeTime += dt
	e.position.Y -= 20 * dt // 上昇
	return e.lifeTime < e.maxLife
}

func (e *FloatingTextEffect) Draw(camera *Camera) {
	// 点滅効果（終了間際）
	if e.lifeTime > e.maxLife*0.8 && int(e.lifeTime*20)%2 == 0 {
		return
	}
	
	screenPos := camera.WorldToScreen(e.position)
	// DrawOutlinedTextはdraw_utils.goで定義する
	DrawWavyText(e.text, Round(screenPos.X), Round(screenPos.Y), e.color, 0, e.lifeTime)
}

func (e *FloatingTextEffect) OnCoordinateReset(offsetX float32) {
	e.position.X -= offsetX
}

// ParticleEffect は飛び散るパーティクル
type ParticleEffect struct {
	position Vector2d
	velocity Vector2d
	color    int
	lifeTime float32
	maxLife  float32
}

func NewParticleEffect(x, y float32, color int) *ParticleEffect {
	// ランダムな速度 (XorShift)
	vx := (xorShift() - 0.5) * 60
	vy := (xorShift() - 0.5) * 60 - 30 // 少し上向き
	return &ParticleEffect{
		position: Vector2d{x, y},
		velocity: Vector2d{vx, vy},
		color:    color,
		lifeTime: 0,
		maxLife:  0.8,
	}
}

func (e *ParticleEffect) Update(dt float32) bool {
	e.lifeTime += dt
	e.velocity.Y += 200 * dt // 重力
	e.position = e.position.Add(e.velocity.Multiply(dt))
	return e.lifeTime < e.maxLife
}

func (e *ParticleEffect) Draw(camera *Camera) {
	screenPos := camera.WorldToScreen(e.position)
	tic80.Pix(Round(screenPos.X), Round(screenPos.Y), e.color)
}

func (e *ParticleEffect) OnCoordinateReset(offsetX float32) {
	e.position.X -= offsetX
}

// TransferEffect はプレイヤー間の受け渡しエフェクト（稲妻のようなライン）
type TransferEffect struct {
	start    Vector2d
	end      Vector2d
	speed    float32 // プレイヤー移動速度に追従するため
	lifeTime float32
	maxLife  float32
}

func NewTransferEffect(start, end Vector2d, speed float32) *TransferEffect {
	return &TransferEffect{
		start:    start,
		end:      end,
		speed:    speed,
		lifeTime: 0,
		maxLife:  0.2, // 0.2秒間表示
	}
}

func (e *TransferEffect) Update(dt float32) bool {
	e.lifeTime += dt
	// プレイヤーと同じ速度で移動させる（簡易実装）
	move := e.speed * dt
	e.start.X += move
	e.end.X += move
	return e.lifeTime < e.maxLife
}

func (e *TransferEffect) Draw(camera *Camera) {
	startPos := camera.WorldToScreen(e.start)
	endPos := camera.WorldToScreen(e.end)
	
	gradient := []int{12, 8, 3, 4, 5, 6, 10}
	
	count := len(gradient)
	dx := (endPos.X - startPos.X)
	dy := (endPos.Y - startPos.Y)
	t := (e.lifeTime / e.maxLife)
	
	for i := 0; i < count; i++ {
		// 位置を計算 (後ろは遅れる)
		// 矩形を描く
		x := startPos.X + dx*float32(i)
		y1 := startPos.Y + dy * easeHalfLinear(t, float32(count - i) / float32(count) * 0.5 + 0.5)
		y2 := startPos.Y + dy * easeHalfLinear(t, float32(count - i) / float32(count) * 0.5)

		var top, bottom float32
		if y1 < y2 {
			top = y1
			bottom = y2
		} else {
			top = y2
			bottom = y1
		}
		
		// サイズは適当に
		size := 8
		
		tic80.Rect(Round(x)-size/2, Round(top), size, Round(bottom-top), gradient[i])
	}
	
	centerX := (startPos.X + endPos.X) / 2
	centerY := (startPos.Y + dy * (0.25 + t * 0.1))
	
	// tの値(0.0〜1.0)をグラデーション配列のインデックスにマッピング
	colorIndex := int(t * float32(count - 1))
	

	DrawOutlinedText("<-", Round(centerX)-4, Round(centerY)-3, gradient[colorIndex], 0)
}

func (e *TransferEffect) OnCoordinateReset(offsetX float32) {
	e.start.X -= offsetX
	e.end.X -= offsetX
}

// HoleEffect は掘った跡（穴）
type HoleEffect struct {
	position Vector2d
	radius   int
	color    int
	lifeTime float32
	maxLife  float32
}

func NewHoleEffect(x, y float32) *HoleEffect {
	return &HoleEffect{
		position: Vector2d{x, y},
		radius:   7,
		color:    0, // 黒
		lifeTime: 0,
		maxLife:  5.0, // 5秒間表示
	}
}

func (e *HoleEffect) Update(dt float32) bool {
	e.lifeTime += dt
	return e.lifeTime < e.maxLife
}

func (e *HoleEffect) Draw(camera *Camera) {
	screenPos := camera.WorldToScreen(e.position)
	// 画面外なら描画しない簡易チェック
	if screenPos.X < -20 || screenPos.X > 260 {
		return
	}

	// 半径のアニメーション計算
	currentRadius := float32(e.radius)

	const fadeInTime = 0.1
	const fadeOutTime = 1.5

	if e.lifeTime < fadeInTime {
		// 拡大
		currentRadius = currentRadius * (e.lifeTime / fadeInTime)
	} else if e.lifeTime > (e.maxLife - fadeOutTime) {
		// 縮小
		remaining := e.maxLife - e.lifeTime
		currentRadius = currentRadius * (remaining / fadeOutTime)
	}

	// 最小半径は0 (負にならないように)
	if currentRadius < 0 {
		currentRadius = 0
	}

	tic80.Circ(Round(screenPos.X), Round(screenPos.Y), int(currentRadius), e.color)
}

func (e *HoleEffect) OnCoordinateReset(offsetX float32) {
	e.position.X -= offsetX
}
