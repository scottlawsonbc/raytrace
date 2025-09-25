// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// Cylinder represents a finite cylinder with a base center, axis, radius, and height.
type Cylinder struct {
	Origin    r3.Point // Center of the cylinder's base in world units.
	Direction r3.Vec   // Axis direction (does not need to be a unit vector).
	Radius    Distance // Radius of the cylinder in world units.
	Height    Distance // Height of the cylinder in world units.
}

func (c Cylinder) Validate() error {
	if c.Radius <= 0 {
		return fmt.Errorf("invalid radius: %v (has it been set?)", c.Radius)
	}
	if c.Height <= 0 {
		return fmt.Errorf("invalid height: %v (has it been set?)", c.Height)
	}
	if c.Direction.IsZero() {
		return fmt.Errorf("invalid direction: %v (has it been set?)", c.Direction)
	}
	eps := 1e-6
	if c.Direction.Length() < 1-eps || c.Direction.Length() > 1+eps {
		length := c.Direction.Length()
		return fmt.Errorf("direction should be a unit vector, got: %v which has length %v", c.Direction, length)
	}
	return nil
}

// Collide determines if the ray intersects with the finite cylinder.
// It returns a boolean indicating a hit and the collision details.
func (c Cylinder) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	d := c.Direction.Unit() // Ensure the axis is a unit vector.
	oc := r.origin.Sub(c.Origin)

	dDotRd := d.Dot(r.direction)
	dDotOc := d.Dot(oc)

	// Components perpendicular to the cylinder axis.
	rdPerp := r.direction.Sub(d.Muls(dDotRd))
	ocPerp := oc.Sub(d.Muls(dDotOc))

	a := rdPerp.Dot(rdPerp)
	b := 2.0 * rdPerp.Dot(ocPerp)
	cVal := ocPerp.Dot(ocPerp) - float64(c.Radius*c.Radius)

	var closestT float64 = math.MaxFloat64
	var closestCollision collision
	hit := false

	// Check for intersections with the cylindrical surface.
	if a > eps { // Avoid division by zero for parallel rays.
		discriminant := b*b - 4*a*cVal
		if discriminant >= 0 {
			sqrtD := math.Sqrt(discriminant)
			t1 := (-b - sqrtD) / (2.0 * a)
			t2 := (-b + sqrtD) / (2.0 * a)

			for _, t := range []float64{t1, t2} {
				if t < float64(tmin) || t > float64(tmax) {
					continue
				}
				// Compute the y coordinate along the axis.
				y := dDotOc + t*dDotRd
				if y >= 0 && y <= float64(c.Height) {
					if t < closestT {
						at := r.at(Distance(t))
						normal := at.Sub(c.Origin.Add(d.Muls(y))).Unit()
						closestT = t
						closestCollision = collision{
							t:      Distance(t),
							at:     at,
							normal: normal,
						}
						hit = true
					}
				}
			}
		}
	}

	// Define the top and bottom caps.
	caps := []struct {
		center r3.Point
		normal r3.Vec
	}{
		{
			center: c.Origin,
			normal: d.Muls(-1), // Bottom cap normal.
		},
		{
			center: c.Origin.Add(d.Muls(float64(c.Height))),
			normal: d, // Top cap normal.
		},
	}

	// Check for intersections with the caps.
	for _, cap := range caps {
		denom := cap.normal.Dot(r.direction)
		if math.Abs(denom) < eps {
			// Ray is parallel to the cap.
			continue
		}
		t := cap.normal.Dot(cap.center.Sub(r.origin)) / denom
		if t < float64(tmin) || t > float64(tmax) {
			continue
		}
		// Compute the intersection point.
		p := r.at(Distance(t))
		// Check if the point is within the cap's radius.
		if p.Sub(cap.center).Dot(p.Sub(cap.center)) <= float64(c.Radius*c.Radius) {
			if t < closestT {
				closestT = t
				closestCollision = collision{
					t:      Distance(t),
					at:     p,
					normal: cap.normal,
					uv:     r2.Point{X: 0.5, Y: 0.5},
				}
				hit = true
			}
		}
	}

	return hit, closestCollision
}

func (c Cylinder) Bounds() AABB {
	d := c.Direction.Unit()
	var orthogonal r3.Vec
	if math.Abs(d.X) > math.Abs(d.Y) {
		orthogonal = r3.Vec{X: -d.Z, Y: 0, Z: d.X}.Unit()
	} else {
		orthogonal = r3.Vec{X: 0, Y: d.Z, Z: -d.Y}.Unit()
	}

	u := orthogonal
	v := d.Cross(u)

	// Compute all 8 corners
	var points []r3.Point
	for i := 0; i <= 1; i++ {
		base := c.Origin.Add(d.Muls(float64(i) * float64(c.Height)))
		for theta := 0.0; theta < 2*math.Pi; theta += math.Pi / 4 { // More points
			circPoint := base.Add(u.Muls(float64(c.Radius) * math.Cos(theta))).Add(v.Muls(float64(c.Radius) * math.Sin(theta)))
			points = append(points, circPoint)
		}
	}

	// Initialize min and max with first point
	min := points[0]
	max := points[0]
	for _, p := range points[1:] {
		min.X = math.Min(min.X, p.X)
		min.Y = math.Min(min.Y, p.Y)
		min.Z = math.Min(min.Z, p.Z)
		max.X = math.Max(max.X, p.X)
		max.Y = math.Max(max.Y, p.Y)
		max.Z = math.Max(max.Z, p.Z)
	}

	return AABB{Min: min, Max: max}
}

func init() {
	RegisterInterfaceType(Cylinder{})
}
