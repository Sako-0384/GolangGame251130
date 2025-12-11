package game

// 簡易乱数生成器 (math/randを使わない)
var rngState uint32 = 123456789

func xorShiftUint() uint32 {
	rngState ^= rngState << 13
	rngState ^= rngState >> 17
	rngState ^= rngState << 5
	return rngState
}

func xorShift() float32 {
	return float32(xorShiftUint()) / 4294967295.0
}

// Intn accepts a max integer and returns a random integer in [0, max).
func Intn(max int) int {
	if max <= 0 { return 0 }
	return int(xorShiftUint() % uint32(max))
}

// Round rounds a float32 to the nearest integer.
func Round(x float32) int {
	if x < 0 {
		return int(x - 0.5)
	}
	return int(x + 0.5)
}

// 簡易Sin関数 (mathパッケージを使わないため)
func fastSin(x float64) float64 {
	// 0 ~ 2PI に正規化 (簡易版)
	const PI = 3.1415926535
	const TWO_PI = 6.283185307
	
	// x を 0以上 TWO_PI未満にする
	for x < 0 { x += TWO_PI }
	for x >= TWO_PI { x -= TWO_PI }
	
	// Bhaskara I approximation
	sign := 1.0
	if x > PI {
		x -= PI
		sign = -1.0
	}
	
	numerator := 16 * x * (PI - x)
	denominator := 5 * PI * PI - 4 * x * (PI - x)
	
	return sign * (numerator / denominator)
}

func minF(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func maxF(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func easeHalfLinear(t float32, mid float32) float32 {
	if t < 0.5 {
		return t * 2.0 * mid
	}
	return mid + (t - 0.5) * 2.0 * (1.0 - mid)
}
