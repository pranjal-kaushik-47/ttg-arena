package entity

import (
	"image/color"
	"math"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Point struct {
	X float64
	Y float64
}

type Polygon struct {
	Vertices []ebiten.Vertex `json:"vertices"`
	Indices  []uint16        `json:"indices"`
	Color    Color           `json:"color"`
	IsActive bool
}

func (p *Polygon) Draw(screen *ebiten.Image) error {
	// op := &ebiten.DrawTrianglesOptions{}
	// op.Address = ebiten.AddressUnsafe
	// img := ebiten.NewImage(1, 1) // A 1x1 image to be used as a placeholder
	// // img.Fill(p.Color.ToColor())
	// screen.DrawTriangles(p.Vertices, p.Indices, img, op)
	// return nil
	for i := 0; i < len(p.Indices)-1; i++ {
		p1 := p.Vertices[p.Indices[i]]
		p2 := p.Vertices[p.Indices[i+1]]
		vector.StrokeLine(screen, p1.DstX, p1.DstY, p2.DstX, p2.DstY, 1, color.White, false)
	}
	return nil
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
func (poly1 *Polygon) PolygonCollision(poly2Vert []ebiten.Vertex) bool {
	axes := GetAxes(poly1.Vertices)
	axes = append(axes, GetAxes(poly2Vert)...)

	for _, axis := range axes {
		min1, max1 := ProjectPolygon(poly1.Vertices, axis)
		min2, max2 := ProjectPolygon(poly2Vert, axis)
		if max1 < min2 || max2 < min1 {
			return false
		}
	}
	return true
}
