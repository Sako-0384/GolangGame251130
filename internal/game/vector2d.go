package game

type Vector2d struct {
	X float32
	Y float32
}

func (v Vector2d) Add(v1 Vector2d) Vector2d {
	return Vector2d{v.X + v1.X, v.Y + v1.Y}
}

func (v Vector2d) Sub(v1 Vector2d) Vector2d {
	return Vector2d{v.X - v1.X, v.Y - v1.Y}
}

func (v Vector2d) Multiply(f float32) Vector2d {
	return Vector2d{v.X * f, v.Y * f}
}

func (v Vector2d) Dot(v1 Vector2d) float32 {
	return v.X*v1.X + v.Y*v1.Y
}

func (v Vector2d) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y
}


