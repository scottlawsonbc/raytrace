// // Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// BVH represents a bounding volume hierarchy of shapes.
type BVH struct {
	Left   Shape
	Right  Shape
	bounds AABB
}

// Ensure BVH implements the Shape interface.
var _ Shape = (*BVH)(nil)

// Validate checks if the BVH is valid.
func (b *BVH) Validate() error {
	if b.Left == nil || b.Right == nil {
		return fmt.Errorf("BVH nodes must not be nil")
	}
	if err := b.Left.Validate(); err != nil {
		return fmt.Errorf("BVH Left child is invalid: %v", err)
	}
	if err := b.Right.Validate(); err != nil {
		return fmt.Errorf("BVH Right child is invalid: %v", err)
	}
	return nil
}

// Bounds returns the bounding box of the BVH node.
func (b *BVH) Bounds() AABB {
	return b.bounds
}

// Collide checks for collision between a ray and the BVH.
func (b *BVH) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	if !b.bounds.hit(r, tmin, tmax) {
		return false, collision{}
	}

	var hitLeft, hitRight bool
	var collLeft, collRight collision

	// Early termination and tmax update
	if b.Left != nil {
		hitLeft, collLeft = b.Left.Collide(r, tmin, tmax)
		if hitLeft {
			tmax = Distance(math.Min(float64(tmax), float64(collLeft.t)))
		}
	}

	if b.Right != nil && b.Right != b.Left {
		hitRight, collRight = b.Right.Collide(r, tmin, tmax)
	}

	if !hitLeft && !hitRight {
		return false, collision{}
	}

	if hitLeft && hitRight {
		if collLeft.t < collRight.t {
			return true, collLeft
		}
		return true, collRight
	} else if hitLeft {
		return true, collLeft
	} else {
		return true, collRight
	}
}

// NewBVH constructs a BVH from a list of shapes using the Binned Surface Area Heuristic.
func NewBVH(shapes []Shape, depth int) *BVH {
	const maxDepth = 32
	const minShapesPerLeaf = 4
	const numBins = 16

	if len(shapes) == 0 {
		return nil
	}

	// If only one shape, create a leaf node.
	if len(shapes) == 1 {
		return &BVH{
			Left:   shapes[0],
			Right:  shapes[0],
			bounds: shapes[0].Bounds(),
		}
	}

	// If maximum depth reached or few shapes, create a leaf node.
	if depth >= maxDepth || len(shapes) <= minShapesPerLeaf {
		// Group shapes into a leaf node.
		group := &Group{Shapes: shapes}
		return &BVH{
			Left:   group,
			Right:  group,
			bounds: group.Bounds(),
		}
	}

	// Compute bounding box of all shapes.
	var bbox AABB
	bbox = shapes[0].Bounds()
	for _, shape := range shapes[1:] {
		bbox = bbox.Union(shape.Bounds())
	}

	// Choose the best axis to split along.
	axis := bbox.LongestAxis()

	// Precompute shape information.
	type shapeInfo struct {
		shape    Shape
		bounds   AABB
		centroid float64
	}
	shapeInfos := make([]shapeInfo, len(shapes))
	for i, shape := range shapes {
		bounds := shape.Bounds()
		centroid := bounds.center().Get(axis)
		shapeInfos[i] = shapeInfo{
			shape:    shape,
			bounds:   bounds,
			centroid: centroid,
		}
	}

	// Implement Binned SAH.
	type bin struct {
		bounds AABB
		count  int
	}
	bins := make([]bin, numBins)
	for i := range bins {
		bins[i].bounds = AABB{
			Min: r3.Point{X: math.Inf(1), Y: math.Inf(1), Z: math.Inf(1)},
			Max: r3.Point{X: math.Inf(-1), Y: math.Inf(-1), Z: math.Inf(-1)},
		}
	}

	// Compute binning.
	for _, si := range shapeInfos {
		centroid := si.centroid
		binIndex := int(numBins * ((centroid - bbox.Min.Get(axis)) / (bbox.Max.Get(axis) - bbox.Min.Get(axis))))
		if binIndex == numBins {
			binIndex = numBins - 1
		}
		bins[binIndex].count++
		bins[binIndex].bounds = bins[binIndex].bounds.Union(si.bounds)
	}

	// Compute SAH cost for each possible split.
	leftCounts := make([]int, numBins)
	rightCounts := make([]int, numBins)
	leftBounds := make([]AABB, numBins)
	rightBounds := make([]AABB, numBins)

	// Initialize left counts and bounds.
	count := 0
	bounds := AABB{
		Min: r3.Point{X: math.Inf(1), Y: math.Inf(1), Z: math.Inf(1)},
		Max: r3.Point{X: math.Inf(-1), Y: math.Inf(-1), Z: math.Inf(-1)},
	}
	for i := 0; i < numBins; i++ {
		count += bins[i].count
		bounds = bounds.Union(bins[i].bounds)
		leftCounts[i] = count
		leftBounds[i] = bounds
	}

	// Initialize right counts and bounds.
	count = 0
	bounds = AABB{
		Min: r3.Point{X: math.Inf(1), Y: math.Inf(1), Z: math.Inf(1)},
		Max: r3.Point{X: math.Inf(-1), Y: math.Inf(-1), Z: math.Inf(-1)},
	}
	for i := numBins - 1; i >= 0; i-- {
		count += bins[i].count
		bounds = bounds.Union(bins[i].bounds)
		rightCounts[i] = count
		rightBounds[i] = bounds
	}

	// Find the best split.
	totalSA := bbox.surfaceArea()
	bestCost := math.MaxFloat64
	bestSplit := -1

	for i := 0; i < numBins-1; i++ {
		pLeft := leftBounds[i].surfaceArea() / totalSA
		pRight := rightBounds[i+1].surfaceArea() / totalSA
		cost := 1 + (float64(leftCounts[i])*pLeft + float64(rightCounts[i+1])*pRight)

		if cost < bestCost {
			bestCost = cost
			bestSplit = i
		}
	}

	// Partition shapes into left and right based on the split.
	if bestSplit == -1 {
		// If no good split found, split shapes equally.
		mid := len(shapeInfos) / 2
		sort.Slice(shapeInfos, func(i, j int) bool {
			return shapeInfos[i].centroid < shapeInfos[j].centroid
		})
		leftShapes := make([]Shape, mid)
		rightShapes := make([]Shape, len(shapeInfos)-mid)
		for i := 0; i < mid; i++ {
			leftShapes[i] = shapeInfos[i].shape
		}
		for i := mid; i < len(shapeInfos); i++ {
			rightShapes[i-mid] = shapeInfos[i].shape
		}

		// Recursively build child BVH nodes.
		return &BVH{
			Left:   NewBVH(leftShapes, depth+1),
			Right:  NewBVH(rightShapes, depth+1),
			bounds: bbox,
		}
	}

	leftShapes := []Shape{}
	rightShapes := []Shape{}
	for _, si := range shapeInfos {
		centroid := si.centroid
		binIndex := int(numBins * ((centroid - bbox.Min.Get(axis)) / (bbox.Max.Get(axis) - bbox.Min.Get(axis))))
		if binIndex == numBins {
			binIndex = numBins - 1
		}
		if binIndex <= bestSplit {
			leftShapes = append(leftShapes, si.shape)
		} else {
			rightShapes = append(rightShapes, si.shape)
		}
	}

	// Parallelize BVH construction if the number of shapes is large enough.
	var left, right *BVH
	if len(shapes) > 1 { // Adjust the threshold based on performance testing.
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			left = NewBVH(leftShapes, depth+1)
			wg.Done()
		}()
		go func() {
			right = NewBVH(rightShapes, depth+1)
			wg.Done()
		}()
		wg.Wait()
	} else {
		left = NewBVH(leftShapes, depth+1)
		right = NewBVH(rightShapes, depth+1)
	}

	return &BVH{
		Left:   left,
		Right:  right,
		bounds: bbox,
	}
}

// Implement custom JSON marshalling for BVH
func (b *BVH) MarshalJSON() ([]byte, error) {
	type BVHData struct {
		Type   string          `json:"Type"`
		Left   json.RawMessage `json:"Left"`
		Right  json.RawMessage `json:"Right"`
		Bounds AABB            `json:"Bounds"`
	}
	leftData, err := marshalInterface(b.Left)
	if err != nil {
		return nil, err
	}
	rightData, err := marshalInterface(b.Right)
	if err != nil {
		return nil, err
	}
	data := BVHData{
		Type:   "BVH",
		Left:   leftData,
		Right:  rightData,
		Bounds: b.bounds,
	}
	return json.Marshal(data)
}

// Implement custom JSON unmarshalling for BVH
func (b *BVH) UnmarshalJSON(data []byte) error {
	type BVHData struct {
		Type   string          `json:"Type"`
		Left   json.RawMessage `json:"Left"`
		Right  json.RawMessage `json:"Right"`
		Bounds AABB            `json:"Bounds"`
	}
	var temp BVHData
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if temp.Type != "BVH" {
		return fmt.Errorf("invalid type: expected BVH, got %s", temp.Type)
	}
	leftShape, err := unmarshalInterface(temp.Left)
	if err != nil {
		return err
	}
	rightShape, err := unmarshalInterface(temp.Right)
	if err != nil {
		return err
	}
	b.Left = leftShape.(Shape)
	b.Right = rightShape.(Shape)
	b.bounds = temp.Bounds
	return nil
}

// Implement custom JSON marshalling for Group
func (g *Group) MarshalJSON() ([]byte, error) {
	type GroupData struct {
		Type   string            `json:"Type"`
		Shapes []json.RawMessage `json:"Shapes"`
	}
	shapesData := make([]json.RawMessage, len(g.Shapes))
	for i, shape := range g.Shapes {
		data, err := marshalInterface(shape)
		if err != nil {
			return nil, err
		}
		shapesData[i] = data
	}
	data := GroupData{
		Type:   "Group",
		Shapes: shapesData,
	}
	return json.Marshal(data)
}

// Implement custom JSON unmarshalling for Group
func (g *Group) UnmarshalJSON(data []byte) error {
	type GroupData struct {
		Type   string            `json:"Type"`
		Shapes []json.RawMessage `json:"Shapes"`
	}
	var temp GroupData
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if temp.Type != "Group" {
		return fmt.Errorf("invalid type: expected Group, got %s", temp.Type)
	}
	shapes := make([]Shape, len(temp.Shapes))
	for i, shapeData := range temp.Shapes {
		shape, err := unmarshalInterface(shapeData)
		if err != nil {
			return err
		}
		shapes[i] = shape.(Shape)
	}
	g.Shapes = shapes
	return nil
}

// Group represents a group of shapes, used as leaf nodes in the BVH.
type Group struct {
	Shapes []Shape
}

// Ensure Group implements the Shape interface.
var _ Shape = (*Group)(nil)

// Validate checks if the Group is valid.
func (g *Group) Validate() error {
	if len(g.Shapes) == 0 {
		return fmt.Errorf("Group must contain at least one shape")
	}
	for i, shape := range g.Shapes {
		if shape == nil {
			return fmt.Errorf("Group shape at index %d is nil", i)
		}
		if err := shape.Validate(); err != nil {
			return fmt.Errorf("Group shape at index %d is invalid: %v", i, err)
		}
	}
	return nil
}

// Bounds computes the bounding box of the group.
func (g *Group) Bounds() AABB {
	bbox := g.Shapes[0].Bounds()
	for _, shape := range g.Shapes[1:] {
		bbox = bbox.Union(shape.Bounds())
	}
	return bbox
}

// Collide checks for collision between a ray and any shape in the group.
func (g *Group) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	hitAnything := false
	var closestCollision collision
	closestT := tmax
	for _, shape := range g.Shapes {
		hit, coll := shape.Collide(r, tmin, closestT)
		if hit {
			hitAnything = true
			closestT = coll.t
			closestCollision = coll
		}
	}
	return hitAnything, closestCollision
}

// String returns a string representation of the BVH.
func (b *BVH) String() string {
	return fmt.Sprintf("BVH{Left: %v, Right: %v, Bounds: %v}", b.Left, b.Right, b.bounds)
}

func init() {
	RegisterInterfaceType(BVH{})
	RegisterInterfaceType(Group{})
}
