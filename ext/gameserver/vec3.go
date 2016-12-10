package gameserver

import (
	"fmt"
	"math"
)

type Len_t int64

type Vec3 struct {
	X Len_t
	Y Len_t
	Z Len_t
}

func (p Vec3) DistanceTo(p2 Vec3) Len_t {
	dx := p.X - p2.X
	dy := p.Y - p2.Y
	dz := p.Z - p2.Z
	return Len_t(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

func (p Vec3) DistanceSquareTo(p2 Vec3) Len_t {
	dx := p.X - p2.X
	dy := p.Y - p2.Y
	dz := p.Z - p2.Z
	return dx*dx + dy*dy + dz*dz
}

func (p Vec3) Add(p2 Vec3) Vec3 {
	return Vec3{
		X: p.X + p2.X,
		Y: p.Y + p2.Y,
		Z: p.Z + p2.Z,
	}
}

func (p Vec3) Sub(p2 Vec3) Vec3 {
	return Vec3{
		X: p.X - p2.X,
		Y: p.Y - p2.Y,
		Z: p.Z - p2.Z,
	}
}

func (p *Vec3) Assign(x, y, z Len_t) {
	p.X = x
	p.Y = y
	p.Z = z
}

func (p Vec3) String() string {
	return fmt.Sprintf("(%v,%v,%v)", p.X, p.Y, p.Z)
}
