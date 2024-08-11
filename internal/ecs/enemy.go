package ecs

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyMemory struct {
	LastPosition Position
}

type Enemy struct {
	Position
	Sprite      Sprite
	Velocity    float64
	BoundingBox BoundingBox
	Eyesight    float64
	EnemyMemory EnemyMemory
}

func Distance(enemy *Enemy, player Player) float64 {
	return math.Sqrt(math.Pow(player.X-enemy.X, 2) + math.Pow(player.Y-enemy.Y, 2))
}

func GetMovementDirection(enemy *Enemy, player Player, env Environment, screenWidth, screenHeight int) (float64, float64) {
	var x, y float64
	if player.X <= enemy.X {
		if enemy.X+enemy.Velocity < float64(screenWidth) && !env.Colliding(enemy.BoundingBox.BoundingBoxShiftRight(2)) {
			enemy.BoundingBox.Polygon.Vertices = enemy.BoundingBox.BoundingBoxShiftRight(enemy.Velocity).Vertices
			x = enemy.Velocity
		} else {
			x = 0
		}
	} else {
		if enemy.X-enemy.Velocity > 0 && !env.Colliding(enemy.BoundingBox.BoundingBoxShiftLeft(2)) {
			enemy.BoundingBox.Polygon.Vertices = enemy.BoundingBox.BoundingBoxShiftLeft(enemy.Velocity).Vertices
			x = -1 * enemy.Velocity
		} else {
			x = 0
		}
	}
	if player.Y <= enemy.Y {
		if enemy.Y+enemy.Velocity <= float64(screenHeight) && !env.Colliding(enemy.BoundingBox.BoundingBoxShiftDown(2)) {
			enemy.BoundingBox.Polygon.Vertices = enemy.BoundingBox.BoundingBoxShiftDown(enemy.Velocity).Vertices
			y = enemy.Velocity
		} else {
			y = 0
		}
	} else {
		if enemy.Y-enemy.Velocity >= 0 && !env.Colliding(enemy.BoundingBox.BoundingBoxShiftUp(2)) {
			enemy.BoundingBox.Polygon.Vertices = enemy.BoundingBox.BoundingBoxShiftUp(enemy.Velocity).Vertices
			y = -1 * enemy.Velocity
		} else {
			y = 0
		}
	}
	if x == 0 && y == 0 {
		screenCenterX, screenCenterY := float64(screenWidth/2), float64(screenHeight/2)
		if screenCenterX > enemy.X {
			x = enemy.Velocity
		} else {
			x = -1 * enemy.Velocity
		}
		if screenCenterY > enemy.Y {
			y = enemy.Velocity
		} else {
			y = -1 * enemy.Velocity
		}
	}
	return x, y
}

func (e *Enemy) Update(p Player, env *Environment, screenHeight, screenWidth int) error {
	distanceFromPlayer := Distance(e, p)
	if distanceFromPlayer <= e.Eyesight {
		x, y := GetMovementDirection(e, p, *env, screenWidth, screenHeight)
		e.MoveTo(e.X+x, e.Y+y, &p.BoundingBox, *env, screenHeight, screenWidth)
		e.EnemyMemory = EnemyMemory{
			LastPosition: Position{
				X: e.X,
				Y: e.Y,
			},
		}
	}
	if e.BoundingBox.PolygonCollision(p.BoundingBox.Polygon.Vertices) && e.Sprite.IsActive {
		e.Sprite.IsActive = false
		env.AliveEnemyCount--
	}
	return nil
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	e.Sprite.Draw(screen, e.X, e.Y)
}
