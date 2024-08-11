package ecs

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Position struct {
	X, Y float64
}

type Polygon struct {
	Vertices []ebiten.Vertex
	Indices  []uint16
	Color    Color
}

type BoundingBox struct {
	Polygon Polygon
}

func (b *BoundingBox) Init(X, Y, Height, Width float64) {
	b.Polygon = Polygon{
		Vertices: []ebiten.Vertex{
			{DstX: float32(X), DstY: float32(Y), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			{DstX: float32(X) + float32(Width)*15, DstY: float32(Y), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			{DstX: float32(X) + float32(Width)*15, DstY: float32(Y) + float32(Width)*20, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			{DstX: float32(X), DstY: float32(Y) + float32(Width)*20, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		},
		Indices: []uint16{1, 0, 3, 1, 2, 3}, // square
		Color:   Color{R: 1, G: 1, B: 1, A: 1},
	}
}

func (b *BoundingBox) AppluBoundingBox(polygon Polygon) {
	b.Polygon = polygon
}

func (b *BoundingBox) BoundingBoxShiftRight(shiftby float64) Polygon {
	VerticesCopy := make([]ebiten.Vertex, len(b.Polygon.Vertices))
	copy(VerticesCopy, b.Polygon.Vertices)
	for i := 0; i < len(VerticesCopy); i += 1 {

		VerticesCopy[i].DstX += float32(shiftby)
	}
	return Polygon{Vertices: VerticesCopy}
}

func (b *BoundingBox) BoundingBoxShiftLeft(shiftby float64) Polygon {
	VerticesCopy := make([]ebiten.Vertex, len(b.Polygon.Vertices))
	copy(VerticesCopy, b.Polygon.Vertices)
	for i := 0; i < len(VerticesCopy); i += 1 {

		VerticesCopy[i].DstX -= float32(shiftby)
	}
	return Polygon{Vertices: VerticesCopy}
}

func (b *BoundingBox) BoundingBoxShiftUp(shiftby float64) Polygon {
	VerticesCopy := make([]ebiten.Vertex, len(b.Polygon.Vertices))
	copy(VerticesCopy, b.Polygon.Vertices)
	for i := 0; i < len(VerticesCopy); i += 1 {
		VerticesCopy[i].DstY -= float32(shiftby)
	}
	return Polygon{Vertices: VerticesCopy}
}

func (b *BoundingBox) BoundingBoxShiftDown(shiftby float64) Polygon {
	VerticesCopy := make([]ebiten.Vertex, len(b.Polygon.Vertices))
	copy(VerticesCopy, b.Polygon.Vertices)
	for i := 0; i < len(VerticesCopy); i += 1 {

		VerticesCopy[i].DstY += float32(shiftby)
	}
	return Polygon{Vertices: VerticesCopy}
}

// collider using Separating Axis Theorem (SAT)
func Dot(v1, v2 ebiten.Vertex) float64 {
	return float64(v1.DstX)*float64(v2.DstX) + float64(v1.DstY)*float64(v2.DstY)
}

// ProjectPolygon projects the vertices of a polygon onto an axis and returns the min and max projection values.
func ProjectPolygon(vertices []ebiten.Vertex, axis ebiten.Vertex) (float64, float64) {
	min := Dot(axis, vertices[0])
	max := min
	for _, vertex := range vertices {
		proj := Dot(axis, vertex)
		if proj < min {
			min = proj
		}
		if proj > max {
			max = proj
		}
	}
	return min, max
}

// GetAxes returns the perpendicular (normal) axes of a polygon's edges.
func GetAxes(polygon []ebiten.Vertex) []ebiten.Vertex {
	axes := make([]ebiten.Vertex, len(polygon))
	for i := range polygon {
		p1 := polygon[i]
		p2 := polygon[(i+1)%len(polygon)]
		edge := ebiten.Vertex{DstX: p2.DstX - p1.DstX, DstY: p2.DstY - p1.DstY}
		normal := ebiten.Vertex{DstX: -edge.DstY, DstY: edge.DstX} // Perpendicular vector
		length := math.Sqrt(float64(normal.DstX*normal.DstX) + float64(normal.DstY*normal.DstY))
		axes[i] = ebiten.Vertex{DstX: float32(float64(normal.DstX) / length), DstY: float32(float64(normal.DstY) / length)} // Normalize
	}
	return axes
}

// PolygonCollision checks if two convex polygons collide using the Separating Axis Theorem (SAT).
func (b *BoundingBox) PolygonCollision(poly2Vert []ebiten.Vertex) bool {
	axes := GetAxes(b.Polygon.Vertices)
	axes = append(axes, GetAxes(poly2Vert)...)

	for _, axis := range axes {
		min1, max1 := ProjectPolygon(b.Polygon.Vertices, axis)
		min2, max2 := ProjectPolygon(poly2Vert, axis)
		if max1 < min2 || max2 < min1 {
			return false
		}
	}
	return true
}

func (p *Position) MoveTo(x, y float64, box *BoundingBox, env Environment, screenHeight, screenWidth int) {
	if p.X < x && x <= float64(screenWidth) {
		velocity := x - p.X
		if !env.Colliding(box.BoundingBoxShiftRight(velocity)) {
			box.Polygon.Vertices = box.BoundingBoxShiftRight(velocity).Vertices
			p.X += velocity
		}
	}
	if p.X > x && x > 0 {
		velocity := p.X - x
		if !env.Colliding(box.BoundingBoxShiftLeft(velocity)) {
			box.Polygon.Vertices = box.BoundingBoxShiftLeft(velocity).Vertices
			p.X -= velocity
		}
	}
	if p.Y < y && y <= float64(screenHeight) {
		velocity := y - p.Y
		if !env.Colliding(box.BoundingBoxShiftDown(velocity)) {
			box.Polygon.Vertices = box.BoundingBoxShiftDown(velocity).Vertices
			p.Y += velocity
		}
	}
	if p.Y > y && y > 0 {
		velocity := p.Y - y
		if !env.Colliding(box.BoundingBoxShiftUp(velocity)) {
			box.Polygon.Vertices = box.BoundingBoxShiftUp(velocity).Vertices
			p.Y -= velocity
		}
	}

}

func (p *Polygon) Draw(screen *ebiten.Image) {

	// Draw lines over the edges to make only the borders visible
	ebitenutil.DrawLine(screen, float64(p.Vertices[0].DstX), float64(p.Vertices[0].DstY), float64(p.Vertices[1].DstX), float64(p.Vertices[1].DstY), color.White)
	ebitenutil.DrawLine(screen, float64(p.Vertices[1].DstX), float64(p.Vertices[1].DstY), float64(p.Vertices[2].DstX), float64(p.Vertices[2].DstY), color.White)
	ebitenutil.DrawLine(screen, float64(p.Vertices[2].DstX), float64(p.Vertices[2].DstY), float64(p.Vertices[3].DstX), float64(p.Vertices[3].DstY), color.White)
	ebitenutil.DrawLine(screen, float64(p.Vertices[3].DstX), float64(p.Vertices[3].DstY), float64(p.Vertices[0].DstX), float64(p.Vertices[0].DstY), color.White)
}
