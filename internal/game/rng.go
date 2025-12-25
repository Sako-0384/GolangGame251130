package game

// RNG is a random number generator using the XorShift algorithm.
type RNG struct {
	state uint32
}

// NewRNG creates a new RNG with the given seed.
func NewRNG(seed uint32) *RNG {
	if seed == 0 {
		seed = 123456789
	}
	return &RNG{state: seed}
}

// Uint32 returns a random uint32.
func (r *RNG) Uint32() uint32 {
	r.state ^= r.state << 13
	r.state ^= r.state >> 17
	r.state ^= r.state << 5
	return r.state
}

// Float32 returns a random float32 in [0.0, 1.0).
func (r *RNG) Float32() float32 {
	return float32(r.Uint32()) / 4294967295.0
}

// RandomIntn returns a random integer in [0, max).
func (r *RNG) Intn(max int) int {
	if max <= 0 {
		return 0
	}
	return int(r.Uint32() % uint32(max))
}
