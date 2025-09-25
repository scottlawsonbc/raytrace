// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

// Camera defines an interface for objects that can emit rays.
// Given normalized screen coordinates (s, t), where 0 ≤ s, t ≤ 1.
// Camera uses the same convention as the PBR book. X is right, Y is up, and
// the camera looks in the negative Z direction.
// Parameters s and t are used to generate a ray that starts at the camera
// position and goes through the screen coordinates (s, t).
// The parameters are normalized to the range [0, 1].
type Camera interface {
	Cast(s, t float64, rand *Rand) ray
	Validate() error
}
