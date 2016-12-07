package gs_position

import "math"

type len_t int64

type Vec3 struct {
	X len_t
	Y len_t
	Z len_t
}

func (p Vec3) DistanceTo(p2 Vec3) len_t {
	dx := p.X - p2.X
	dy := p.Y - p2.Y
	dz := p.Z - p2.Z
	return len_t(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

func (p Vec3) DistanceSquareTo(p2 Vec3) len_t {
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
