package game

// defaultRNG is the package-level default random number generator.
var defaultRNG = NewRNG(123456789)

// SetRandomSeed sets the seed for the global default random number generator.
func SetRandomSeed(seed uint32) {
	defaultRNG = NewRNG(seed)
}

func RandomFloat32() float32 {
	return defaultRNG.Float32()
}

// RandomIntn accepts a max integer and returns a random integer in [0, max) using the default RNG.
func RandomIntn(max int) int {
	return defaultRNG.Intn(max)
}

// Round rounds a float32 to the nearest integer.
func Round(x float32) int {
	if x < 0 {
		return int(x - 0.5)
	}
	return int(x + 0.5)
}

const PI = 3.1415926535
const TWO_PI = 6.283185307

// 簡易Sin関数 (mathパッケージを使わないため)
func fastSin(x float64) float64 {
	// x を 0以上 TWO_PI未満にする
	for x < 0 {
		x += TWO_PI
	}
	for x >= TWO_PI {
		x -= TWO_PI
	}

	// Bhaskara I approximation
	sign := 1.0
	if x > PI {
		x -= PI
		sign = -1.0
	}

	numerator := 16 * x * (PI - x)
	denominator := 5*PI*PI - 4*x*(PI-x)

	return sign * (numerator / denominator)
}

func easeHalfLinear(t float32, mid float32) float32 {
	if t < 0.5 {
		return t * 2.0 * mid
	}
	return mid + (t-0.5)*2.0*(1.0-mid)
}

func intToString(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + intToString(-i)
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}

// Clamp clamps v to the range [min, max].
func Clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
